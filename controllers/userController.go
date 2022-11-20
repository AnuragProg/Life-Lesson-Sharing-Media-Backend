package controllers

import(
	"net/http"
	"rest-api/models"
	"github.com/gin-gonic/gin"	
	"go.mongodb.org/mongo-driver/mongo"
)

func AddUserHandler(coll *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){
		var userData models.UserRequest
		if err := c.BindJSON(&userData); err != nil{
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error(), Response: ""})
			return
		}
		_ , err := userData.AddUser(coll)

		if err != nil{
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error(), Response: ""})
			return
		}
		c.JSON(http.StatusOK, models.GeneralResponse{Message:"Successfully added user"})	
	}
}


func GetUsersHandler(coll *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){
		users, err := models.GetUsers(coll)
		if err != nil{
			c.JSON(http.StatusInternalServerError, models.GeneralResponse{Message:err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}


func UpdateUserHandler(coll *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){
		var userData models.User
		if err := c.BindJSON(&userData); err != nil{
			c.JSON(http.StatusBadRequest, models.GeneralResponse{Message:err.Error()})
			return
		}
		_ , err := userData.UpdateUser(coll)
		if err!= nil{
			c.JSON(http.StatusInternalServerError, models.GeneralResponse{Message:err.Error()})
			return
		}
		c.JSON(http.StatusOK, userData)
	}
}

/*
Requires Query (id: userId)
*/
func DeleteUserHandler(pllColl *mongo.Collection, commentColl *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){
		userId := c.Query("id")
		if userId == ""{
			c.JSON(http.StatusBadRequest, models.GeneralResponse{Message:"Couldn't find 'id' in query"})
			return
		}
		deleteResult , err := models.DeleteUser(userId, pllColl, commentColl)
		if err!= nil{
			c.JSON(http.StatusInternalServerError, models.GeneralResponse{Message:err.Error()})
			return
		}

		if deleteResult.DeletedCount == 0 {
			c.JSON(http.StatusInternalServerError, models.GeneralResponse{Message:"Unable to delete user!"})
			return
		}
		c.JSON(http.StatusOK, models.GeneralResponse{Message:"Successfully deleted User!"})
	} 
}