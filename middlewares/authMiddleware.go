package middlewares

import(
	"github.com/gin-gonic/gin"
	"net/http"
	"rest-api/models"
)

// Verify JWT Token
func verifyJwt(token string) bool {
	return true
}

// JWT TOKEN
func AuthMiddleware(c *gin.Context){
	token := c.Request.Header.Get("Authorization")

	if token == "" || token =="Bearer "{
		c.IndentedJSON(http.StatusUnauthorized, models.ErrorResponse{ Message: "No token found" })
		return
	}

	if verifyJwt(token){
		c.Next()
		return
	}
	c.IndentedJSON(http.StatusUnauthorized, models.ErrorResponse{ Message: "Invalid token"})
}