package middlewares

import (
	"log"
	"context"
	"errors"
	"net/http"
	"rest-api/models"
	"rest-api/components"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserAuthMiddlwareHandler(coll *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){

		// Retrieving JWT token from request
		tokenString, err := components.GetBearerToken(c)
		if err != nil{
			c.Abort()
			return
		}

		// Retrieving JWT Secret for verifying JWT signature
		secret, err := components.GetJWTSecret()
		if err != nil{
			c.Abort()
			return 
		}

		// Parsing the JWT token retreiving from request
		token, err:= jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok{
				return "", errors.New("unable to parse token")
			}
			return secret, nil
		})
		if err != nil{
			c.JSON(404, gin.H{"message":err.Error()})
			c.Abort()
			return
		}

		// Checking Validity of parsed JWT Token
		if claims, ok := token.Claims.(jwt.MapClaims); token.Valid && ok{
			email := claims["email"]
			pass := claims["password"]

			// Retrieving user with credentials in JWT token
			filter := bson.M{"email": email}
			result := coll.FindOne(context.TODO(), filter)
			var user models.User
			err := result.Decode(&user)
			if err!=nil{
				c.JSON(http.StatusUnauthorized, gin.H{"message":"user info not in database, need to sign up again"})
				c.Abort()
				return
			}

			// Verifying whether JWT Token retrieved was the latest assigned one
			if user.LastToken != tokenString{
				c.JSON(http.StatusUnauthorized, gin.H{"message":"use newly assigned token"})
				c.Abort()
				return
			}

			// Verifying credentials of JWT token with credentials saved in db
			// JWT contains hashed password and db also contains hashed password so simple == will do work
			if user.Email != email || user.Password != pass{
				log.Println("Expected email", user.Email, " and got from token email", email)
				log.Println("Expected password", user.Password, " and got from token email", pass)
				c.JSON(http.StatusUnauthorized, gin.H{"message":"credentials not valid due to incorrect password"})
				c.Abort()
				return
			}
			c.Set(components.USERIDKEY, user.ID)
			c.Next()
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"message":"Invalid jwt token"})
		c.Abort()
	}
}



func AdminAuthMiddlwareHandler(coll *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){

		// Retrieving token from request
		tokenString, err := components.GetBearerToken(c)
		if err != nil{
			c.Abort()
			return
		}

		// Retrieving JWT Secret for verifying JWT signature
		secret, err := components.GetJWTSecret()
		if err != nil{
			c.Abort()
			return 
		}

		// Parsing token string retrieving from request and verifying signature with JWT secret 
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok{
				return "", errors.New("unable to parse token")
			}
			return secret, nil
		})
		if err != nil{
			c.JSON(404, gin.H{"message":err.Error()})
			c.Abort()
			return
		}

		// Checking whether JWT token is valid or not
		if claims, ok := token.Claims.(jwt.MapClaims); token.Valid && ok{
			admin := claims["admin"].(bool)
			email := claims["email"].(string)
			pass := claims["password"].(string)

			// Retrieving user from db to validate credentials of JWT token claims
			var user models.User
			filter := bson.M{"email": email}
			result := coll.FindOne(context.TODO(), filter)
			err := result.Decode(&user)
			if err!=nil{
				c.JSON(http.StatusUnauthorized, gin.H{"message":"user info not in database, need to sign up again"})
				c.Abort()
				return
			}

			// Verifying whether JWT Token retrieved was the latest assigned one
			if user.LastToken != tokenString{
				c.JSON(http.StatusUnauthorized, gin.H{"message":"use newly issued token"})
				c.Abort()
				return
			}

			// Validating credentials of JWT token with credentails in db
			if user.Email != email || user.Password != pass{
				log.Println("Expected email", user.Email, " and got from token email", email)
				log.Println("Expected password", user.Password, " and got from token email", pass)
				c.JSON(http.StatusUnauthorized, gin.H{"message":"credentials not valid due to incorrect password"})
				c.Abort()
				return
			}

			// Checking whether the provided user is admin or not
			if !admin || !user.IsAdmin{
				c.JSON(http.StatusUnauthorized, gin.H{"message":"Not admin user"})
				c.Abort()
				return
			}else{

				// User is admin with correct credentials
				// Let the user move to the next handler
				c.Set(components.USERIDKEY, user.ID)
				c.Next()
				return
			}
		}

		c.JSON(http.StatusUnauthorized, gin.H{"message":"Invalid jwt token"})
		c.Abort()
	}
}
