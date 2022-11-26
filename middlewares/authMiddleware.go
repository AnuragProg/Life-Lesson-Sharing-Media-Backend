package middlewares

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"
	"rest-api/models"
	"golang.org/x/crypto/bcrypt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
note:
	1. JWT token password is hashed password
	2. Hashing only takes place in Signing Up and Logging in using password
*/

const COST = 5
const(
	USERIDKEY string = "UserId"
)


func GetBearerToken(c *gin.Context) (string, error){
	tokenString := c.Request.Header.Get("Authorization")	
	if tokenString == "" || tokenString == "Bearer "{
		return "", errors.New("no token string found")
	}
	tokenSlice := strings.Split(tokenString, " ")
	if len(tokenSlice) != 2{
		return "", errors.New("invalid token")
	}
	token := tokenSlice[1]
	return token, nil
}

func GetJWTSecret() ([]byte, error){
	token := os.Getenv("JWT_SECRET")
	if token == ""{
		return []byte(""), errors.New("no jwt secret found")
	}
	return []byte(token), nil
}

func LoginUserWithPasswordHandler(coll *mongo.Collection)gin.HandlerFunc{
	return func(c *gin.Context){
		var credentials struct{
			Email string `json:"email"`
			Password string `json:"password"`
		}

		// Finding user with provided credentials
		filter := bson.M{"email":credentials.Email}
		result := coll.FindOne(context.TODO(), filter)
		var user models.User
		err := result.Decode(&user)
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"message":err.Error()})
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
		token, err := GenerateJWTToken(user.IsAdmin, user.Email, user.Password)
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"message":err.Error()})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{"token":token})
	}
}


func LoginUserWithTokenHandler(coll *mongo.Collection)gin.HandlerFunc{
	return func(c *gin.Context){

		// Retrieving token from request
		tokenString, err := GetBearerToken(c)
		if err != nil{
			c.JSON(404, gin.H{"message":err.Error()})
			c.Abort()
			return 
		}

		// Retrieving Secret for verifying JWT token retrieved from request
		secret, err := GetJWTSecret()
		if err != nil{
			panic("Unable to find jwt secret")
		}
		secretKey := []byte(secret)

		// Parsing JWT token from string
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

		// Verifying Parsed JWT token
		if claims, ok := token.Claims.(jwt.MapClaims); token.Valid && ok{
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


			// Verifying user with credentials contained in JWT token
			if user.Email == email && user.Password == password{

				// Generating and sending new JWT token
				tokenString,err := GenerateJWTToken(admin, user.Email, user.Password)
				if err != nil{
					c.JSON(404, gin.H{"message":err.Error()})
					c.Abort()
					return
				}
				c.JSON(http.StatusOK, gin.H{"token":tokenString})
				return
			}
		}
		c.JSON(http.StatusUnauthorized, gin.H{"message":"Invalid token"})
	}
}


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
		hashedPassword,err := bcrypt.GenerateFromPassword([]byte(user.Password), COST)
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
		token, err := GenerateJWTToken(isAdmin, user.Email, user.Password)
		if err != nil{
			c.JSON(404, gin.H{"message":err.Error()})
			c.Abort()
			return
		}

		// Adding user to db
		// Done after all conversions to make user user is ready to be added to db
		_, err = user.ToUserIntermediate(isAdmin).AddUser(coll)
		if err != nil{
			c.JSON(404, gin.H{"message":err.Error()})
			c.Abort()
			return
		}

		// Returning the token
		c.JSON(http.StatusOK, gin.H{"token":token})
	}
}

func UserAuthMiddlwareHandler(coll *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){

		// Retrieving JWT token from request
		tokenString, err := GetBearerToken(c)
		if err != nil{
			c.Abort()
			return
		}

		// Retrieving JWT Secret for verifying JWT signature
		secret, err := GetJWTSecret()
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
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message:"user info not in database, need to sign up again"})
				c.Abort()
				return
			}

			// Verifying credentials of JWT token with credentials saved in db
			// JWT contains hashed password and db also contains hashed password so simple == will do work
			if user.Email != email || user.Password != pass{
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message:"credentials not valid due to incorrect password"})
				c.Abort()
				return
			}
			c.Set(USERIDKEY, user.ID)
			c.Next()
			return
		}
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message:"Invalid jwt token"})
		c.Abort()
	}
}



func AdminAuthMiddlwareHandler(coll *mongo.Collection) gin.HandlerFunc{
	return func(c *gin.Context){

		// Retrieving token from request
		tokenString, err := GetBearerToken(c)
		if err != nil{
			c.Abort()
			return
		}


		// Retrieving JWT Secret for verifying JWT signature
		secret, err := GetJWTSecret()
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
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message:"user info not in database, need to sign up again"})
				c.Abort()
				return
			}

			// Validating credentials of JWT token with credentails in db
			if user.Email != email || user.Password != pass{
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message:"credentials not valid due to incorrect password"})
				c.Abort()
				return
			}

			// Checking whether the provided user is admin or not
			if !admin || !user.IsAdmin{
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message:"Not admin user"})
				c.Abort()
				return
			}else{

				// User is admin with correct credentials
				// Let the user move to the next handler
				c.Set(USERIDKEY, user.ID)
				c.Next()
				return
			}
		}

		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message:"Invalid jwt token"})
		c.Abort()
	}
}


func GenerateJWTToken(admin bool, email, password string)(string, error){

	secret, err := GetJWTSecret()
	if err != nil{
		return "", err
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["admin"] = admin
	claims["email"] = email
	claims["password"] = password
	claims["exp"] = time.Now().Add(time.Hour*24).Unix()
	return token.SignedString(secret)
}