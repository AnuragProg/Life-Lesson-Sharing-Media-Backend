package controllers

import (
	"net/http"
	"rest-api/middlewares"
	"rest-api/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddPllHandler(coll *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){

		//Retrieving userID after token verification
		userId, exists := c.Get(middlewares.USERIDKEY)
		if !exists{
			c.JSON(http.StatusInternalServerError, gin.H{"message":"internal server error"})
			c.Abort()
			return
		}

		// Retrieving Pll from request body 
		var pllRequest models.PersonalLifeLessonRequest
		if err := c.BindJSON(&pllRequest); err != nil{
			c.JSON(http.StatusBadRequest, models.GeneralResponse{Message: err.Error()})
			return
		}

		// setting userId to ensure true ownership of pll
		pllRequest.UserId = userId.(string)
		_ , err := pllRequest.AddPll(coll)

		if err != nil{
			c.JSON(http.StatusInternalServerError, models.GeneralResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, models.GeneralResponse{Message:"Successfully added personal life lesson"})	
	}
}

func GetPllsHandler(coll *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){
		plls, err := models.GetPlls(coll)
		if err != nil{
			c.JSON(http.StatusInternalServerError, models.GeneralResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, plls)	
	}
}

func GetPllHandler(coll *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){
		pllId := c.Query("id")
		if pllId == ""{
			c.JSON(http.StatusBadRequest, models.GeneralResponse{Message: "Cannot find 'id' in query"})
			return
		}
		pll, err := models.GetPll(pllId, coll)
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.GeneralResponse{Message: err.Error()})
			return
		}
		if pll == nil{
			c.JSON(http.StatusBadRequest, models.GeneralResponse{Message: "Cannot find data with give id!"})
			return
		}
		c.JSON(http.StatusOK, *pll)	
	}
}


func UpdatePllHandler(coll *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){

		// Get User id from verified token
		userId, exists := c.Get(middlewares.USERIDKEY)
		if !exists{
			c.JSON(http.StatusInternalServerError, gin.H{"message":"internal sever error"})
			c.Abort()
			return
		}

		// Get Pll from request body
		var pll models.PersonalLifeLesson
		if err:= c.BindJSON(&pll); err != nil{
			c.JSON(http.StatusBadRequest, models.GeneralResponse{Message:err.Error()})
			return
		}

		// Sending userid to validate Ownership of given user over given pll
		result, err := pll.UpdatePll(userId.(string), coll)
		if err != nil{
			c.JSON(http.StatusInternalServerError, models.GeneralResponse{Message:err.Error()})
			return
		}

		// Final check to whether modification took place or not
		if result.ModifiedCount == 0{
			c.JSON(http.StatusInternalServerError, models.GeneralResponse{Message: "Unable to modify!"})
			return
		}
		c.JSON(http.StatusOK, models.GeneralResponse{Message:"Successfully updated"})
	}
}

func LikePllsHandler(pllColl *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){

		// Retrieving UserId after token verification
		userId, exists := c.Get(middlewares.USERIDKEY)
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
		userId, exists := c.Get(middlewares.USERIDKEY)
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