package controllers

import(
	"net/http"
	"rest-api/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)


func AddCommentHandler(coll *mongo.Collection) gin.HandlerFunc{
	fun := func(c *gin.Context){
		var comment models.CommentRequest
		if err := c.BindJSON(&comment); err != nil{
			c.JSON(http.StatusBadRequest, models.GeneralResponse{Message: err.Error()})
			return
		}
		_ , err := comment.AddComment(coll)
		if err != nil{
			c.JSON(http.StatusInternalServerError, models.GeneralResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, models.GeneralResponse{Message:"Successfully added comment"})	
	}
	return fun
}

func UpdateCommentHandler(coll *mongo.Collection) gin.HandlerFunc{
	fun := func(c *gin.Context){
		var comment models.Comment
		if err:=c.BindJSON(&comment); err != nil{
			c.JSON(http.StatusBadRequest, models.GeneralResponse{Message: err.Error()})
			return
		}
		result, err := comment.UpdateComment(coll) 
		if err != nil{
			c.JSON(http.StatusBadRequest, models.GeneralResponse{Message: err.Error()})
			return
		}
		if result.ModifiedCount == 0{
			c.JSON(http.StatusInternalServerError, models.GeneralResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, models.GeneralResponse{Message: "Successfully updated comment"})
	}
	return fun
}

func DeleteCommentHandler(coll *mongo.Collection) gin.HandlerFunc{
	fun := func(c *gin.Context){
		var commentId struct{ CommentId string `json:"commentId"`} 
		if err := c.BindJSON(&commentId); err != nil{
			c.JSON(http.StatusBadRequest, models.GeneralResponse{Message: err.Error()})
			return
		}
		result, err := models.DeleteComment(commentId.CommentId, coll)
		if err != nil{
			c.JSON(http.StatusBadRequest, models.GeneralResponse{Message: err.Error()})
			return
		}
		if result.DeletedCount == 0{
			c.JSON(http.StatusInternalServerError, models.GeneralResponse{Message:"Unable to delete!"})
			return
		}
		c.JSON(http.StatusOK, models.GeneralResponse{Message:"Successfully deleted comment"})
	}
	return fun
}


func GetCommentsHandler(coll *mongo.Collection) gin.HandlerFunc{
	fun := func(c *gin.Context){
		var commentIds struct{CommentIds []string `json:"commentIds"`} 
		if err:=c.BindJSON(&commentIds); err != nil{
			c.JSON(http.StatusBadRequest, models.GeneralResponse{Message: err.Error()})
			return
		}
		comments, err := models.GetComments(commentIds.CommentIds, coll)
		if err != nil{
			c.JSON(http.StatusInternalServerError, models.GeneralResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, comments)
	}
	return fun
}