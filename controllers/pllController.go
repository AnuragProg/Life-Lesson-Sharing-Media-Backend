package controllers

import (
	"context"
	"net/http"
	"rest-api/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeleteCommentFromTableAndListHandler(
	pllColl *mongo.Collection, commentColl *mongo.Collection,
)gin.HandlerFunc{
	return func(c *gin.Context){
		commentId := c.Query("id")
		if commentId == ""{
			c.JSON(http.StatusBadRequest, models.GeneralResponse{Message:"Couldn't find comment 'id' in Query"})
			return 
		}
		id, err := primitive.ObjectIDFromHex(commentId)
		if err != nil{
			c.JSON(http.StatusBadRequest, models.GeneralResponse{Message: err.Error()})
			return 
		}
		var comment models.Comment
		filter := bson.M{"_id":id}
		request := commentColl.FindOne(context.TODO(), filter)
		err = request.Decode(&comment)
		if err != nil{
			c.JSON(http.StatusBadRequest, models.GeneralResponse{Message: err.Error()})
			return 
		}
		err = models.DeleteCommentFromTableAndList(comment.PllId, comment.ID, pllColl, commentColl)
		if err != nil{
			c.JSON(http.StatusBadRequest, models.GeneralResponse{Message: err.Error()})
			return 
		}
		c.JSON(http.StatusBadRequest, models.GeneralResponse{Message: "Successfully deleted comment"})
	}
}

func AddCommentToTableAndListHandler(
	pllColl *mongo.Collection, commentColl *mongo.Collection,
)gin.HandlerFunc{
	return func(c *gin.Context){
		var comment models.CommentRequest
		if err:=c.BindJSON(&comment); err != nil {
			c.JSON(http.StatusBadRequest, models.GeneralResponse{Message: "Invalid request body!"} )
			return
		}
		err := comment.AddCommentToTableAndList(comment.PllId, pllColl, commentColl)
		if err != nil{
			c.JSON(http.StatusInternalServerError, models.GeneralResponse{Message: err.Error()} )
			return
		}
		c.JSON(http.StatusOK, models.GeneralResponse{Message:"Successfully commented"})
	}
}

func AddPllHandler(coll *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){
		var pllRequest models.PersonalLifeLessonRequest
		if err := c.BindJSON(&pllRequest); err != nil{
			c.JSON(http.StatusBadRequest, models.GeneralResponse{Message: err.Error()})
			return
		}
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

/*
Requires Query (id: PllId)
*/
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

/*
Requires Query (id: PllId)
*/
func DeletePllAndCommentsFromTableHandler(pllColl *mongo.Collection, commentColl *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){
		pllId := c.Query("id")
		if pllId == ""{
			c.JSON(http.StatusBadRequest, models.GeneralResponse{Message: "Cannot find 'id' in query"})
			return
		}

		result, err := models.DeletePllAndCommentsFromTable(pllId, pllColl, commentColl)
		if err != nil{
			c.JSON(http.StatusInternalServerError, models.GeneralResponse{Message: err.Error()})
			return
		}
		if result.DeletedCount == 0{
			c.JSON(http.StatusInternalServerError, models.GeneralResponse{Message:"Unable to delete!"})
			return
		}

		c.JSON(http.StatusOK, models.GeneralResponse{Message:"Successfully deleted"})
	}
}

/*
Requires Body PersonalLifeLesson
*/
func UpdatePllHandler(coll *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){
		var pll models.PersonalLifeLesson
		if err:= c.BindJSON(&pll); err != nil{
			c.JSON(http.StatusBadRequest, models.GeneralResponse{Message:err.Error()})
			return
		}
		result, err := pll.UpdatePll(coll)
		if err != nil{
			c.JSON(http.StatusInternalServerError, models.GeneralResponse{Message:err.Error()})
			return
		}

		if result.ModifiedCount == 0{
			c.JSON(http.StatusInternalServerError, models.GeneralResponse{Message: "Unable to modify!"})
			return
		}
		c.JSON(http.StatusOK, models.GeneralResponse{Message:"Successfully updated"})
	}
}

/*
Requires Query (id: userId)
Returns Comment
*/
func GetPllCommentsHandler(coll *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){
		pllId := c.Query("id")
		if pllId == ""{
			c.JSON(http.StatusBadRequest, models.GeneralResponse{Message: "Didn't find 'id' in query!"})	
			return 
		}
		comments, err := models.GetPllComments(pllId, coll)
		if err != nil{
			c.JSON(http.StatusInternalServerError, models.GeneralResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, comments)
	}
}