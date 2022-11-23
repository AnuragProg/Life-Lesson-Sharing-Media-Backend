package controllers

import(
	"net/http"
	"rest-api/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddCommentHandler(pllColl, commentColl *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){
		var comment models.CommentRequest

		// TODO Verify userid authority

		err := c.BindJSON(&comment);
		if err != nil{
			c.JSON(404, models.GeneralResponse{Message:err.Error()})
			return
		}

		_, err = comment.AddComment(pllColl, commentColl)
		
		if err != nil{
			c.JSON(404, models.GeneralResponse{Message:err.Error()})
			return
		}

		
		c.JSON(http.StatusOK, models.GeneralResponse{Message:"Successfully added comment"})
	}
}

func UpdateCommentHandler(coll *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){
		var comment models.Comment
		err := c.BindJSON(&comment)
		if err != nil{
			c.JSON(404, models.GeneralResponse{Message:err.Error()})
			return
		}

		// TODO Verify authority of given userId

		result, err := comment.UpdateComment(coll)
		if err != nil{
			c.JSON(404, models.GeneralResponse{Message:err.Error()})
			return
		}
		if result.ModifiedCount == 0{
			c.JSON(404, models.GeneralResponse{Message: "Unable to update comment"})
			return
		}
		c.JSON(http.StatusOK, models.GeneralResponse{Message:"Successfully updated comment"})
	}
}

func DeleteCommentHandler(pllColl, commentColl *mongo.Collection)gin.HandlerFunc{
	return func(c *gin.Context){
		var commentDeletionRequest struct{
			UserId string `json:"userId"`
			CommentId string `json:"commentId"`
		}

		// TODO Verify authority of given userID on comment

		err := c.BindJSON(&commentDeletionRequest)
		if err != nil{
			c.JSON(404, models.GeneralResponse{Message:err.Error()})
			return
		}
		err = models.DeleteComment(commentDeletionRequest.CommentId, pllColl, commentColl)
		if err != nil{
			c.JSON(404, models.GeneralResponse{Message:err.Error()})
			return
		}
		c.JSON(http.StatusOK, models.GeneralResponse{Message:"Successfully deleted comment"})
	}
}

func DeleteCommentsHandler(coll *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){
		var commentDeletionRequest struct{
			UserId string `json:"userId"`
			CommentIds []string `json:"commentIds"`
		}

		//TODO Verify authority of given user

		err := c.BindJSON(&commentDeletionRequest)
		
		if err != nil{
			c.JSON(404, models.GeneralResponse{Message:err.Error()})
			return
		}
		err = models.DeleteComments(commentDeletionRequest.CommentIds, coll)
		if err != nil{
			c.JSON(404, models.GeneralResponse{Message:err.Error()})
			return
		}
		c.JSON(http.StatusOK, models.GeneralResponse{Message:"Successfully deleted comment"})
	}
}

func GetCommentsHandler(coll *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){
		var commentIds []string
		err := c.BindJSON(&commentIds)
		if err != nil{
			c.JSON(404, models.GeneralResponse{Message:err.Error()})
			return
		}
		comments := models.GetComments(commentIds, coll)
		c.JSON(http.StatusOK, comments)
	}
}