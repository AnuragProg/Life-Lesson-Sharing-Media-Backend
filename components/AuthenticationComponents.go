package components

import(
	"context"
	"os"
	"strings"
	"errors"
	"time"
	"rest-api/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


const(
	USERIDKEY string = "UserId"
)

// Parses authorization token from request and construct necessary error message
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

// Retreives JWT Secret from environment variable
func GetJWTSecret() ([]byte, error){
	token := os.Getenv("JWT_SECRET")
	if token == ""{
		return []byte(""), errors.New("no jwt secret found")
	}
	return []byte(token), nil
}


// Generate JWT token with given data
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


type DataType uint16

const (
	PERSONALLIFELESSON DataType = iota
	COMMENT 
)
/*
 Checks whether "data" with "dataId"
 can be manipulated by "user" with given "userId"
*/
func CheckAuthority(userId, dataId string, datatype DataType, coll *mongo.Collection)(bool, error){

	// Generating object id and finding the document associated with the dataId
	objectId, err := primitive.ObjectIDFromHex(dataId)
	if err != nil{
		return false,errors.New("invalid data id") 
	}
	filter := bson.M{"_id": objectId}
	result:= coll.FindOne(context.TODO(), filter)

	// Decoding the document according to the requested data type
	switch datatype{
	case PERSONALLIFELESSON:
		var pll models.PersonalLifeLesson
		if err := result.Decode(&pll); err != nil{
			return false, err
		}
		return pll.UserId == userId, nil
	case COMMENT:
		var comment models.Comment
		if err := result.Decode(&comment); err != nil{
			return false, err
		}
		return comment.UserId == userId, nil
	default:
		return false, errors.New("invalid data sent for checking authority")
	}
}
