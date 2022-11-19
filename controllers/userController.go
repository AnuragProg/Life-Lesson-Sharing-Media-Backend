package controllers

import(
	"rest-api/models"
	"net/http"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/gin-gonic/gin"	
)

func AddUserHandler(coll *mongo.Collection) gin.HandlerFunc{
	fun := func(c *gin.Context){
		var userData models.UserRequest
		if err := c.BindJSON(&userData); err != nil{
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error(), Response: ""})
			return
		}
		_ , err := userData.AddUser(coll)

		if err != nil{
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error(), Response: ""})
		}
		c.JSON(http.StatusOK, models.GeneralResponse{Message:"Successfully added user"})	
	}
	return fun
}


func GetUsersHandler(coll *mongo.Collection) gin.HandlerFunc{
	fun := func(c *gin.Context){
		users, err := models.GetUsers(coll)
		if err != nil{
			c.JSON(http.StatusInternalServerError, models.GeneralResponse{Message:err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	}
	return fun
}


func UpdateUserHandler(coll *mongo.Collection) gin.HandlerFunc{
	fun := func(c *gin.Context){
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
	return fun
}
func DeleteUserHandler(coll *mongo.Collection) gin.HandlerFunc{
	fun := func(c *gin.Context){
		var userId struct{ UserId string `json:"userId"`}
		if err := c.BindJSON(&userId); err != nil{
			c.JSON(http.StatusBadRequest, models.GeneralResponse{Message:err.Error()})
			return
		}
		deleteResult , err := models.DeleteUser(userId.UserId, coll)
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
	return fun
}