package controllers

import (
	"context"
	"errors"
	"net/http"
	"rest-api/components"
	"rest-api/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Everytime user opens the app
// During the splash screen, this login should take place
func LoginUserWithTokenHandler(coll *mongo.Collection)gin.HandlerFunc{
	return func(c *gin.Context){

		// Retrieving token from request
		tokenString, err := components.GetBearerToken(c)
		if err != nil{
			c.JSON(404, gin.H{"message":err.Error()})
			c.Abort()
			return 
		}

		// Retrieving Secret for verifying JWT token retrieved from request
		secret, err := components.GetJWTSecret()
		if err != nil{
			panic("Unable to find jwt secret")
		}
		secretKey := []byte(secret)

		// Parsing JWT token from string and Verifying the signature
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {

			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return "", errors.New("unable to parse token")
			}
			return secretKey, nil
		})
		if err != nil{
			c.JSON(http.StatusUnauthorized, gin.H{"message":err.Error()})
			c.Abort()
			return
		}

		// Checking if Parsed JWT token is valid
		// token.Valid && // to omit whether the token has expired or not
		// because it may have been expired
		if claims, ok := token.Claims.(jwt.MapClaims); ok{
			email, password, admin:= claims["email"], claims["password"], claims["admin"].(bool)

			// Finding user with credentials in JWT token
			var user models.User
			filter := bson.M{"email": email}			
			result := coll.FindOne(context.TODO(), filter)
			if err := result.Decode(&user); err != nil{
				c.JSON(404, gin.H{"message":"Cannot find user"})
				c.Abort()
				return
			}

			// Check whether this was the token last assigned
			if tokenString != user.LastToken{
				c.JSON(http.StatusUnauthorized, gin.H{"message":"not the latest token assigned"})
				c.Abort()
				return
			}

			if user.Email == email && user.Password == password{

				// Generating and sending new JWT token
				tokenString,err := components.GenerateJWTToken(admin, user.Email, user.Password)
				if err != nil{
					c.JSON(404, gin.H{"message":err.Error()})
					c.Abort()
					return
				}

				// Updating user's last token
				err = models.UpdateLastToken(user.Email, tokenString, coll)
				if err!= nil{
					c.JSON(http.StatusInternalServerError, gin.H{"message":err.Error()})	
					c.Abort()
					return
				}

				// Returning new token to user
				c.JSON(http.StatusOK, gin.H{"token":tokenString, "userId": user.ID})
				return
			}
		}
		c.JSON(http.StatusUnauthorized, gin.H{"message":"Invalid token"})
	}
}


/* Initial sign up for new users */
func SignUpUserHandler(coll *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){

		// Retrieving userrequest body from request body
		var user models.UserRequest
		if err := c.BindJSON(&user); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"message":err.Error()})
			c.Abort()
			return
		}
		
		// Generating hash of user password and replacing with user requested password 
		hashedPassword,err := bcrypt.GenerateFromPassword([]byte(user.Password), models.COST)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"message":err.Error()})
			c.Abort()
			return
		}
		user.Password = string(hashedPassword)


		// Admin accounts to check for admins
		admins := []string{"codingsubs73@gmail.com"}

		// Checking if email provided is admin
		isAdmin := false
		for _, admin:= range admins{
			if admin == user.Email{
				isAdmin = true
				break
			} 
		}

		// Generating token with user credentials
		token, err := components.GenerateJWTToken(isAdmin, user.Email, user.Password)
		if err != nil{
			c.JSON(404, gin.H{"message":err.Error()})
			c.Abort()
			return
		}

		// Adding user to db and adding latest token
		// Done after all conversions to make user user is ready to be added to db
		result, err := user.ToUserIntermediate(isAdmin, token).AddUser(coll)
		if err != nil{
			c.JSON(404, gin.H{"message":err.Error()})
			c.Abort()
			return
		}

		// Returning the token
		c.JSON(http.StatusOK, gin.H{"token":token, "userId":result.InsertedID.(primitive.ObjectID).Hex()})
	}
}

/*
note:
	1. JWT token password is hashed password
	2. Hashing only takes place in Signing Up and Logging in using password
*/
func LoginUserWithPasswordHandler(coll *mongo.Collection)gin.HandlerFunc{
	return func(c *gin.Context){
		var credentials struct{
			Email string `json:"email"`
			Password string `json:"password"`
		}

		// Retreiving body from request
		if err:=c.BindJSON(&credentials); err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{"message":"Invalid request body"})
			c.Abort()
			return
		}

		// Finding user with provided credentials
		filter := bson.M{"email":credentials.Email}
		result := coll.FindOne(context.TODO(), filter)
		var user models.User
		err := result.Decode(&user)
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"message":"No such user found with given email"})
			c.Abort()
			return
		}

		// Compare hash password with plain password
		if user.Email != credentials.Email || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)) != nil{
			c.JSON(http.StatusBadRequest, gin.H{"message":"Wrong email or password"})
			c.Abort()
			return
		}
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"message":err.Error()})
			c.Abort()
			return
		}

		// Generating new token
		token, err := components.GenerateJWTToken(user.IsAdmin, user.Email, user.Password)
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"message":err.Error()})
			c.Abort()
			return
		}

		// Updating previous token
		models.UpdateLastToken(user.Email, token, coll)

		c.JSON(http.StatusOK, gin.H{"token":token, "userId":user.ID})
	}
}