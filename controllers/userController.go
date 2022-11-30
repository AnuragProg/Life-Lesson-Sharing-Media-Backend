package controllers

import (
	"context"
	"net/http"
	"rest-api/components"
	"rest-api/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Only for admin
func GetUsersHandler(coll *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){
		users, err := models.GetUsers(coll)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"message":err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}


func UpdateUserHandler(coll *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){

		// Retrieving user id from token verification
		userId, exists := c.Get(components.USERIDKEY)
		if !exists{
			c.JSON(http.StatusInternalServerError, gin.H{"message":"unable to retreive user id"})
			c.Abort()
			return
		}

		// Retrieving User update request from body
		var userData models.UserUpdateRequest
		if err := c.BindJSON(&userData); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"message":err.Error()})
			c.Abort()
			return
		}

		// Extracting User associated with user id that is extracted from JWT token
		id, err := primitive.ObjectIDFromHex(userId.(string))
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"message":"token contains invalid user id"})
			c.Abort()
			return
		}
		var user models.User
		filter := bson.M{"_id": id}
		result := coll.FindOne(context.TODO(), filter)
		err = result.Decode(&user)
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"message":"no such user with userid in token exists"})
			c.Abort()
			return
		}

		// Generating hashed password
		hashedPassword,err:= bcrypt.GenerateFromPassword([]byte(user.Password), models.COST)
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"message":err.Error()})
			c.Abort()
			return
		}

		// Generating new token with updated email and password
		token, err := components.GenerateJWTToken(user.IsAdmin, user.Email, string(hashedPassword))
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"message":err.Error()})
			c.Abort()
			return
		}

		// Updating user and last token
		_ , err = userData.UpdateUser(userId.(string), string(hashedPassword), token, coll)
		if err!= nil{
			c.JSON(http.StatusInternalServerError, gin.H{"message":err.Error()})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{"token":token})
	}
}


// Requires Query (id: userId)
func DeleteUserHandler(pllColl *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){

		// Retreiving userid from token verification
		userId:= c.GetString(components.USERIDKEY)
		if userId == ""{
			c.JSON(http.StatusInternalServerError, gin.H{"message":"not able to find user id from token"})
			c.Abort()
			return
		}

		// Deleting the user with userId
		deleteResult , err := models.DeleteUser(userId, pllColl)
		if err!= nil{
			c.JSON(http.StatusInternalServerError, gin.H{"message":err.Error()})
			c.Abort()
			return
		}
		if deleteResult.DeletedCount == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"message":"Unable to delete user!"})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, gin.H{"message":"Successfully deleted User!"})
	} 
}