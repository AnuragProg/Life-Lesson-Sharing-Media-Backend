package main

import (
	"context"
	"log"
	"os"
	"rest-api/middlewares"
	"rest-api/controllers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func main(){
	router := gin.New()
	
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

	user := router.Group("/user", middlewares.AuthMiddleware)
	{
		userCollection := db.Collection("Users")
		user.POST("/", controllers.AddUserHandler(userCollection))
		user.GET("/", controllers.GetUsersHandler(userCollection))
		user.PUT("/", controllers.UpdateUserHandler(userCollection))
		user.DELETE("/", controllers.DeleteUserHandler(userCollection))
	}

	pll := router.Group("/pll", middlewares.AuthMiddleware)
	{
		pllCollection := db.Collection("Pll")	
		pll.POST("/", controllers.AddPllHandler(pllCollection))
		pll.GET("/plls", controllers.GetPllsHandler(pllCollection))
		pll.GET("/pll", controllers.GetPllHandler(pllCollection))
		pll.PUT("/", controllers.UpdatePllHandler(pllCollection))
		pll.DELETE("/", controllers.DeletePllHandler(pllCollection))
	}

	comment := router.Group("/comment", middlewares.AuthMiddleware)
	{
		commentCollection := db.Collection("Comments")
		comment.GET("/comments", controllers.GetCommentsHandler(commentCollection))
		comment.POST("/", controllers.AddCommentHandler(commentCollection))
		comment.DELETE("/", controllers.DeleteCommentHandler(commentCollection))
		comment.PUT("/", controllers.UpdateCommentHandler(commentCollection))
	}


	router.Run("localhost:5000")
}
