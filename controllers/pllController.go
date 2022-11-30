package controllers

import (
	"context"
	"net/http"
	"rest-api/components"
	"rest-api/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddPllHandler(pllColl, userColl, categoryColl *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){

		//Retrieving userID after token verification
		userId, exists := c.Get(components.USERIDKEY)
		if !exists{
			c.JSON(http.StatusInternalServerError, gin.H{"message":"internal server error"})
			c.Abort()
			return
		}

		// Retrieving Pll from request body 
		var pllRequest models.PersonalLifeLessonRequest
		if err := c.BindJSON(&pllRequest); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			c.Abort()
			return
		}

		// Verifiying whether given category is present
		if _, err := models.GetCategory(pllRequest.CategoryId, categoryColl); err !=nil{
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			c.Abort()
			return
		}

		// Retrieving User from db for converting request to its intermediate
		var user models.User
		id, err := primitive.ObjectIDFromHex(userId.(string))
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"message":"not a valid userid"})
			c.Abort()
			return
		}
		filter := bson.M{"_id":id}
		result := userColl.FindOne(context.TODO(), filter)
		if err:=result.Decode(&user); err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"message":"Not able to find user with the provided token"})
			c.Abort()
			return
		}

		// Converting request to its intermediate and adding the intermediate to the db
		_, err = pllRequest.ToPersonalLifeLessonRequestIntermediate(user.ID, user.Username).AddPll(pllColl)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK,gin.H{"message":"Successfully added personal life lesson"})	
	}
}

func GetPllsHandler(coll *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){
		plls, err := models.GetPlls(coll)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, plls)	
	}
}

func GetPllHandler(coll *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){
		pllId := c.Query("id")
		if pllId == ""{
			c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot find 'id' in query"})
			return
		}
		pll, err := models.GetPll(pllId, coll)
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		if pll == nil{
			c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot find data with give id!"})
			return
		}
		c.JSON(http.StatusOK, *pll)	
	}
}


func UpdatePllHandler(pllColl, userColl, categoryColl *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){

		// Get User id from verified token
		userId, exists := c.Get(components.USERIDKEY)
		if !exists{
			c.JSON(http.StatusInternalServerError, gin.H{"message":"internal sever error"})
			c.Abort()
			return
		}

		// Get Pll from request body
		var pll models.PersonalLifeLessonUpdateRequest
		if err:= c.BindJSON(&pll); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"message":"Invalid request body"})
			c.Abort()
			return
		}

		// Verifiying whether given category is present
		if _, err := models.GetCategory(pll.CategoryId, categoryColl); err !=nil{
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			c.Abort()
			return
		}

		// Sending userid to validate Ownership of given user over given pll
		authorized, err := components.CheckAuthority(userId.(string), pll.ID, components.PERSONALLIFELESSON, pllColl)
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"message":err.Error()})
			c.Abort()
			return
		}
		if !authorized{
			c.JSON(http.StatusUnauthorized, gin.H{"message":"not authorized to make changes in the post"})
			c.Abort()
			return
		}

		// Retreiving user for providing user info for posting personal life lesson
		var user models.User
		id, err := primitive.ObjectIDFromHex(userId.(string))
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"message":"not able to convert userid to object id"})
			c.Abort()
			return
		}
		filter := bson.M{"_id": id}
		result := userColl.FindOne(context.TODO(), filter)
		if err = result.Decode(&user); err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"message":"no such user exists"})
			c.Abort()
			return
		}

		// Updating the pll
		_, err = pll.UpdatePll(user.ID, user.Username, pllColl)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"message":err.Error()})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{"message":"Successfully updated"})
	}
}

func DeletePllHandler(pllColl *mongo.Collection)gin.HandlerFunc{
	return func(c *gin.Context){

		// Retrieving pll id from request
		pllId := c.Query("id")
		if pllId == ""{
			c.JSON(http.StatusBadRequest, gin.H{"message":"unable to find personal life lesson id in query"})
			c.Abort()
			return
		}

		// Retrieving user id from token verification step
		userId, exists := c.Get(components.USERIDKEY)
		if !exists{
			c.JSON(http.StatusInternalServerError, gin.H{"message":"Unable to extract userid from jwt token"})
			c.Abort()
			return
		}
		// checking the authority of the userid over pll
		authorized, err := components.CheckAuthority(userId.(string), pllId, components.PERSONALLIFELESSON, pllColl)
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"message":err.Error()})
			c.Abort()
			return
		}
		if !authorized{
			c.JSON(http.StatusUnauthorized, gin.H{"message":"not authorized to deleted this personal life lesson post"})
			c.Abort()
			return
		}

		// deleting the pll
		err = models.DeletePll(pllId, pllColl)
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"message":err.Error()})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{"message":"Successfully deleted personal life lesson post"})
	}
}


func LikePllsHandler(pllColl *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){

		// Retrieving UserId after token verification
		userId, exists := c.Get(components.USERIDKEY)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"message":"Internal server error"})
			c.Abort()
			return
		}

		// Retrieving Pll ids from request body
		// Sending userId for Authentic likes
		var pllIds []string
		if err := c.BindJSON(&pllIds); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			c.Abort()
			return 
		}
		models.LikePlls(pllIds, userId.(string), pllColl)
		c.JSON(http.StatusOK, gin.H{"message":"Successfully liked provided personal life lessons"})
	}
}


func DislikePllsHandler(pllColl *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){

		// Retrieving UserId after token verification
		userId, exists := c.Get(components.USERIDKEY)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"message":"Internal server error"})
			c.Abort()
			return
		}

		// Retrieving Plls to dislike from request body
		// Sending userId for Authentic dislikes
		var pllIds []string
		if err := c.BindJSON(&pllIds); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			c.Abort()
			return 
		}
		models.DislikePlls(pllIds, userId.(string), pllColl)
		c.JSON(http.StatusOK, gin.H{"message":"Successfully disliked provided personal life lessons"})
	}
}