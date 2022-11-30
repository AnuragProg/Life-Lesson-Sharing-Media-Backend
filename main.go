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

	user := router.Group("/user")
	{
		// Without any middleware
		user.POST("/signUp", controllers.SignUpUserHandler(userCollection))
		user.POST("/signInWithPassword", controllers.LoginUserWithPasswordHandler(userCollection))

		// Requires User middleware
		user.POST("/signIn", middlewares.UserAuthMiddlwareHandler(userCollection),controllers.LoginUserWithTokenHandler(userCollection))
		user.PATCH("/", middlewares.UserAuthMiddlwareHandler(userCollection),controllers.UpdateUserHandler(userCollection))
		user.DELETE("/", middlewares.UserAuthMiddlwareHandler(userCollection),controllers.DeleteUserHandler(userCollection))

		// Only for admin
		user.GET("/", middlewares.AdminAuthMiddlwareHandler(userCollection),controllers.GetUsersHandler(userCollection))
	}

	pll := router.Group("/pll", middlewares.UserAuthMiddlwareHandler(userCollection))
	{
		pll.GET("/plls", controllers.GetPllsHandler(pllCollection))
		pll.GET("/pll", controllers.GetPllHandler(pllCollection))
		pll.PATCH("/", controllers.UpdatePllHandler(pllCollection, userCollection, categoryCollection))
		pll.POST("/", controllers.AddPllHandler(pllCollection,userCollection, categoryCollection))
		pll.POST("/like", controllers.LikePllsHandler(pllCollection))
		pll.POST("/dislike", controllers.DislikePllsHandler(pllCollection))
		pll.DELETE("/", controllers.DeletePllHandler(pllCollection))
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
		comments.POST("/", controllers.AddCommentHandler(pllCollection,commentCollection, userCollection))
		comments.DELETE("/", controllers.DeleteCommentHandler(pllCollection,commentCollection))
		comments.PATCH("/", controllers.UpdateCommentHandler(commentCollection))
	}

	parentRouter.Run(os.Getenv("BASE_URL"))
}
