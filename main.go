package main

import (
	"context"
	"fmt"
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

	user := router.Group("/user")
	{
		user.POST("/signUp", middlewares.SignUpUserHandler(userCollection))
		user.POST("/signIn", middlewares.LoginUserWithTokenHandler(userCollection))
		user.POST("/signInWithPassword", middlewares.LoginUserWithPasswordHandler(userCollection))
		user.GET("/", middlewares.AdminAuthMiddlwareHandler(userCollection),controllers.GetUsersHandler(userCollection))
		user.PATCH("/edit", middlewares.UserAuthMiddlwareHandler(userCollection),controllers.UpdateUserHandler(userCollection))
		user.DELETE("/", middlewares.UserAuthMiddlwareHandler(userCollection),controllers.DeleteUserHandler(userCollection, commentCollection))
	}

	pll := router.Group("/pll", middlewares.UserAuthMiddlwareHandler(userCollection))
	{
		pll.GET("/plls", controllers.GetPllsHandler(pllCollection))
		pll.GET("/pll", controllers.GetPllHandler(pllCollection))
		pll.PATCH("/", controllers.UpdatePllHandler(pllCollection))
		pll.POST("/", controllers.AddPllHandler(pllCollection))
		pll.POST("/like", controllers.LikePllsHandler(pllCollection))
		pll.POST("/dislike", controllers.DislikePllsHandler(pllCollection))
	}

	category := router.Group("/category")
	{
		category.GET("/categories",middlewares.UserAuthMiddlwareHandler(userCollection), controllers.GetCategoriesHandler(categoryCollection))
		category.GET("/category", middlewares.UserAuthMiddlwareHandler(userCollection), controllers.GetCategoryHandler(categoryCollection))
		category.POST("/", middlewares.AdminAuthMiddlwareHandler(userCollection),controllers.AddCategoryHandler(categoryCollection))
		category.DELETE("/", middlewares.AdminAuthMiddlwareHandler(userCollection),controllers.DeleteCategoryHandler(categoryCollection))
		category.PATCH("/", middlewares.AdminAuthMiddlwareHandler(userCollection),controllers.UpdateCategoryHandler(categoryCollection))
	}

	comments := router.Group("/comment", middlewares.UserAuthMiddlwareHandler(userCollection))
	{
		comments.GET("/", controllers.GetCommentsHandler(commentCollection))
		comments.POST("/", controllers.AddCommentHandler(pllCollection,commentCollection))
		comments.DELETE("/", controllers.DeleteCommentHandler(pllCollection,commentCollection))
		comments.PATCH("/", controllers.UpdateCommentHandler(commentCollection))
	}

	port := 5000
	addr := fmt.Sprintf("localhost:%v", port)

	parentRouter.Run(addr)
}
