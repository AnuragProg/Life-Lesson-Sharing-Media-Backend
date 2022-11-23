package main

import (
	"context"
	"log"
	"os"
	"rest-api/controllers"
	"rest-api/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func main(){
	parentRouter := gin.Default()
	
	router := parentRouter.Group("/v1")

	if err := godotenv.Load(); err!=nil{
		log.Fatal("Cannot find .env file!")
	}
	var uri string
	if uri = os.Getenv("MONGO_URI"); uri == ""{
		log.Fatal("Cannot find MONGO_URI in .env file")
	}
	
	log.Println("Connecting to DB...")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	log.Println("Connected to DB")

	// Disconnecting from the db
	defer func(){
		log.Println("Disconnecting from DB...")
		if err = client.Disconnect(context.TODO()); err != nil{
			log.Fatal(err.Error())
		}
		log.Println("Disconnected DB")
	}()

	db := client.Database("personalLifeLessons_db")

	if err != nil{
		log.Fatal("Cannot connect to DB!")
	}

	userCollection := db.Collection("Users")
	pllCollection := db.Collection("Pll")	
	commentCollection := db.Collection("Comments")
	categoryCollection := db.Collection("Categories")

	user := router.Group("/user", middlewares.AuthMiddleware)
	{
		user.POST("/", controllers.AddUserHandler(userCollection))
		user.GET("/", controllers.GetUsersHandler(userCollection))
		user.PATCH("/", controllers.UpdateUserHandler(userCollection))
		user.DELETE("/", controllers.DeleteUserHandler(userCollection, commentCollection))
	}

	pll := router.Group("/pll", middlewares.AuthMiddleware)
	{
		pll.GET("/plls", controllers.GetPllsHandler(pllCollection))
		pll.GET("/pll", controllers.GetPllHandler(pllCollection))
		pll.PATCH("/", controllers.UpdatePllHandler(pllCollection))
		pll.POST("/", controllers.AddPllHandler(pllCollection))
	}

	category := router.Group("/category", middlewares.AuthMiddleware)
	{
		category.GET("/categories", controllers.GetCategoriesHandler(categoryCollection))
		category.GET("/category", controllers.GetCategoryHandler(categoryCollection))
		category.POST("/", controllers.AddCategoryHandler(categoryCollection))
		category.DELETE("/", controllers.DeleteCategoryHandler(categoryCollection))
		category.PATCH("/", controllers.UpdateCategoryHandler(categoryCollection))
	}

	comments := router.Group("/comment", middlewares.AuthMiddleware)
	{
		comments.GET("/", controllers.GetCommentsHandler(commentCollection))
		comments.POST("/", controllers.AddCommentHandler(pllCollection,commentCollection))
		comments.DELETE("/", controllers.DeleteCommentHandler(pllCollection,commentCollection))
		comments.PATCH("/", controllers.UpdateCommentHandler(commentCollection))
	}


	parentRouter.Run("localhost:5000")
}
