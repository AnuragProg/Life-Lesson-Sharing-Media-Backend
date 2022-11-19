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
		user.POST("/add", controllers.AddUserHandler(userCollection))
		user.GET("/users", controllers.GetUsersHandler(userCollection))
		user.PUT("/update", controllers.UpdateUserHandler(userCollection))
		user.DELETE("/delete", controllers.DeleteUserHandler(userCollection))
	}

	pll := router.Group("/pll", middlewares.AuthMiddleware)
	{
		pllCollection := db.Collection("Pll")	
		pll.POST("/add", controllers.AddPllHandler(pllCollection))
		pll.GET("/plls", controllers.GetPllsHandler(pllCollection))
		pll.GET("/pll", controllers.GetPllHandler(pllCollection))
		pll.PUT("/update", controllers.UpdatePllHandler(pllCollection))
		pll.DELETE("/delete", controllers.DeletePllHandler(pllCollection))
	}


	router.Run("localhost:5000")
}
