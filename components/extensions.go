package components

import(
	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context){
	if pc := recover(); pc!=nil{
		c.JSON(404, gin.H{"message":pc})
	}
}