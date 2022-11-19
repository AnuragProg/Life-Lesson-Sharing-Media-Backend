package controllers

import (
	"rest-api/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"go.mongodb.org/mongo-driver/mongo"
)


func AddPllHandler(coll *mongo.Collection) gin.HandlerFunc{
	fun := func(c *gin.Context){
		var pllRequest models.PersonalLifeLessonRequest
		if err := c.BindJSON(&pllRequest); err != nil{
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error(), Response: ""})
			return
		}
		_ , err := pllRequest.AddPll(coll)

		if err != nil{
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error(), Response: ""})
			return
		}
		c.JSON(http.StatusOK, models.GeneralResponse{Message:"Successfully added personal life lesson"})	
	}
	return fun
}

func GetPllsHandler(coll *mongo.Collection) gin.HandlerFunc{
	fun := func(c *gin.Context){
		plls, err := models.GetPlls(coll)
		if err != nil{
			c.JSON(http.StatusInternalServerError, models.GeneralResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, plls)	
	}
	return fun
}

func GetPllHandler(coll *mongo.Collection) gin.HandlerFunc{
	fun := func(c *gin.Context){
		var pllId struct{PllId string `json:"pllId"`}
		if err := c.BindJSON(&pllId); err !=nil{
			c.JSON(http.StatusBadRequest, models.GeneralResponse{Message: "Invalid request body!"})
			return
		}
		pll, err := models.GetPll(pllId.PllId, coll)
		
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
	return fun
}

func DeletePllHandler(coll *mongo.Collection) gin.HandlerFunc{
	fun := func(c *gin.Context){
		var pllId struct{PllId string `json:"pllId"`}
		if err := c.BindJSON(&pllId); err !=nil{
			c.JSON(http.StatusBadRequest, models.GeneralResponse{Message: "Invalid request body!"})
			return
		}
		result, err := models.DeletePll(pllId.PllId, coll)
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
	return fun
}

func UpdatePllHandler(coll *mongo.Collection) gin.HandlerFunc{
	fun := func(c *gin.Context){
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
	return fun
}