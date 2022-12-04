package controllers

import (
	"net/http"
	"rest-api/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddCategoryHandler(coll *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){
		var category models.CategoryRequest
		if err:=c.BindJSON(&category); err!=nil{
			c.JSON(http.StatusBadRequest, models.GeneralResponse{Message: err.Error()})
			return 
		}
		_, err := category.AddCategory(coll)
		if err != nil{
			c.JSON(http.StatusInternalServerError, models.GeneralResponse{Message:err.Error()})
			return
		}
		c.JSON(http.StatusOK, models.GeneralResponse{Message:"Successfully added category"})
	}
}

func UpdateCategoryHandler(coll *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){
		var category models.Category
		if err := c.BindJSON(&category); err !=nil{
			c.JSON(http.StatusBadRequest, models.GeneralResponse{Message:err.Error()})
			return
		}
		result, err:= category.UpdateCategory(coll)
		if err != nil{
			c.JSON(http.StatusInternalServerError, models.GeneralResponse{Message:err.Error()})
			return
		}
		if result.ModifiedCount == 0{
			c.JSON(http.StatusInternalServerError, models.GeneralResponse{Message:"Unable to update category!"})
			return
		}
		c.JSON(http.StatusOK, models.GeneralResponse{Message: "Successfully updated category"})
	}
}

func GetCategoriesHandler(coll *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){
		categories, err := models.GetCategories(coll)
		if err != nil{
			c.JSON(http.StatusInternalServerError, models.GeneralResponse{Message:err.Error()})
			return
		}
		c.JSON(http.StatusOK, categories)
	}
}

/*
Requires Query (id: categoryId)
*/
func GetCategoryHandler(coll *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){
		categoryId := c.Query("id")
		if categoryId == ""{
			c.JSON(http.StatusBadRequest, models.GeneralResponse{Message:"Couldn't find 'id' in Query"})
			return
		}
		category, err := models.GetCategory(categoryId, coll)
		if err != nil{ 
			c.JSON(http.StatusInternalServerError, models.GeneralResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, category)
	}
}

/*
Requires Query (id: categoryId)
*/
func DeleteCategoryHandler(coll *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){
		categoryId := c.Query("id")
		if categoryId == ""{
			c.JSON(http.StatusBadRequest, models.GeneralResponse{Message:"Cannot find 'id' in query"})
			return
		}
		result, err := models.DeleteCategory(categoryId, coll)
		if err != nil{
			c.JSON(http.StatusInternalServerError, models.GeneralResponse{Message:err.Error()})
			return
		}
		if result.DeletedCount == 0{
			c.JSON(http.StatusInternalServerError, models.GeneralResponse{Message:"Unable to delete category"})
			return
		}
		c.JSON(http.StatusOK, models.GeneralResponse{Message: "Successfully deleted category"})
	}
}