package controllers

import (
	"net/http"
	"rest-api/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)


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