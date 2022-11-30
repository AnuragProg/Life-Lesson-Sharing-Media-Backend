package controllers

import (
	"log"
	"context"
	"net/http"
	"rest-api/models"
	"rest-api/components"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func AddCommentHandler(pllColl, commentColl, userColl *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){

		// Retreiving body from request
		var comment models.CommentRequest
		err := c.BindJSON(&comment);
		if err != nil{
			c.JSON(404, gin.H{"message":err.Error()})
			return
		}

		// Extracting userid from token verification
		userId, exists := c.Get(components.USERIDKEY)
		if !exists{
			c.JSON(http.StatusBadRequest, gin.H{"message":"Unable to find userId"})
			c.Abort()
			return
		}

		// Retreiving user from db to convert commentRequest to commentRequestIntermediate
		var user models.User
		id, err := primitive.ObjectIDFromHex(userId.(string))
		log.Println("userid:", userId, "and objectid is", id)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"message":"invalid userid found"})
			c.Abort()
			return
		}
		filter := bson.M{"_id":id}
		result := userColl.FindOne(context.TODO(), filter)
		if err= result.Decode(&user); err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"message":"cannot find user"})
			c.Abort()
			return
		}

		// Converting CommentRequest to CommentRequestIntermediate
		_, err = comment.ToCommentRequestIntermediate(user.ID, user.Username).AddComment(pllColl, commentColl)
		if err != nil{
			c.JSON(404, gin.H{"message":err.Error()})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{"message":"Successfully added comment"})
	}
}


func UpdateCommentHandler(coll *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){

		// Retreive body from request
		var comment models.CommentUpdateRequest
		err := c.BindJSON(&comment)
		if err != nil{
			c.JSON(404, gin.H{"message":err.Error()})
			return
		}

		// Extracting userid from token verification
		userId, exists := c.Get(components.USERIDKEY)
		if !exists{
			c.JSON(http.StatusBadRequest, gin.H{"message":"Unable to find userId"})
			c.Abort()
			return
		}

		// Verify authority of user over the comment to be updated
		authorized, err := components.CheckAuthority(userId.(string), comment.ID, components.COMMENT, coll)
		if err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{"message":err.Error()})
			c.Abort()
			return
		}
		if !authorized{
			c.JSON(http.StatusUnauthorized, gin.H{"message":"not authorized to update comment"})
			c.Abort()
			return
		}

		// Updating the comment
		result, err := comment.UpdateComment(coll)
		if err != nil{
			c.JSON(404, gin.H{"message":err.Error()})
			return
		}
		if result.ModifiedCount == 0{
			c.JSON(404, gin.H{"message": "Unable to update comment"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message":"Successfully updated comment"})
	}
}


func DeleteCommentHandler(pllColl, commentColl *mongo.Collection)gin.HandlerFunc{
	return func(c *gin.Context){

		// Retrieving commentid from request query
		commentId := c.Query("id")
		if commentId == ""{
			c.JSON(http.StatusBadRequest, gin.H{"message":"unable to find 'id' in query"})
			c.Abort()
			return
		}

		// Extracting userid from token verification
		userId, exists := c.Get(components.USERIDKEY)
		if !exists{
			c.JSON(http.StatusBadRequest, gin.H{"message":"Unable to find userId"})
			c.Abort()
			return
		}

		// Is given User authorized to delete this comment
		authorized, err := components.CheckAuthority(userId.(string), commentId, components.COMMENT, commentColl)
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"message":err.Error()})
			c.Abort()
			return
		}
		if !authorized{
			c.JSON(http.StatusUnauthorized, gin.H{"message":"not authorized to delete this comment"})
			c.Abort()
			return
		}

		// Delete the comment
		err = models.DeleteComment(commentId, pllColl, commentColl)
		if err != nil{
			c.JSON(404, gin.H{"message":err.Error()})
			return
		}
		c.JSON(http.StatusOK,gin.H{"message":"Successfully deleted comment"})
	}
}

func GetCommentsHandler(coll *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){
		var commentIds []string
		err := c.BindJSON(&commentIds)
		if err != nil{
			c.JSON(404, gin.H{"message":err.Error()})
			return
		}
		comments := models.GetComments(commentIds, coll)
		c.JSON(http.StatusOK, comments)
	}
}


// TODO Only for Admin (not to be implemented now)
func DeleteCommentsHandler(coll *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){
		var commentDeletionRequest struct{
			UserId string `json:"userId"`
			CommentIds []string `json:"commentIds"`
		}

		//TODO Verify authority of given user

		err := c.BindJSON(&commentDeletionRequest)
		
		if err != nil{
			c.JSON(404, gin.H{"message":err.Error()})
			return
		}
		err = models.DeleteComments(commentDeletionRequest.CommentIds, coll)
		if err != nil{
			c.JSON(404, gin.H{"message":err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message":"Successfully deleted comment"})
	}
}
