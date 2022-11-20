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
	router := gin.Default()
	
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

	user := router.Group("/user", middlewares.AuthMiddleware)
	{

		// *Working*
		user.POST("/", controllers.AddUserHandler(userCollection))

		// *Working*
		user.GET("/", controllers.GetUsersHandler(userCollection))

		// *Working* | Note: Giving error "message": "the provided hex string is not a valid ObjectID" => should instead give proper error of missing id
		user.PATCH("/", controllers.UpdateUserHandler(userCollection))

		// *Not working*
		user.DELETE("/", controllers.DeleteUserHandler(userCollection, commentCollection))
	}

	pll := router.Group("/pll", middlewares.AuthMiddleware)
	{
		// *Working*
		pll.GET("/plls", controllers.GetPllsHandler(pllCollection))

		// *Working*
		pll.GET("/pll", controllers.GetPllHandler(pllCollection))

		// *Working*
		pll.PATCH("/", controllers.UpdatePllHandler(pllCollection))

		// *Working*
		pll.POST("/", controllers.AddPllHandler(pllCollection))

		//
		pll.DELETE("/", controllers.DeletePllAndCommentsFromTableHandler(pllCollection, commentCollection))

		// *Working*
		pll.POST("/comment", controllers.AddCommentToTableAndListHandler(pllCollection,commentCollection))

		// *Working*
		pll.DELETE("/comment", controllers.DeleteCommentFromTableAndListHandler(pllCollection, commentCollection))
	}

	category := router.Group("/category", middlewares.AuthMiddleware)
	{
		categoryCollection := db.Collection("Categories")

		// *Working*
		category.GET("/", controllers.GetCategoriesHandler(categoryCollection))

		// *Working*
		category.GET("/category", controllers.GetCategoryHandler(categoryCollection))

		// *Working*
		category.POST("/", controllers.AddCategoryHandler(categoryCollection))

		// *Working*
		category.DELETE("/", controllers.DeleteCategoryHandler(categoryCollection))

		// *Working*
		category.PATCH("/", controllers.UpdateCategoryHandler(categoryCollection))
	}

	router.Run("localhost:5000")
}
