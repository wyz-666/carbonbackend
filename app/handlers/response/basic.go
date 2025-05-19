package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MakeSuccessAdmin(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"statusCode": code, "message": message, "data": data})
}

func MakeSuccess(c *gin.Context, code int, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"statusCode": code, "data": data})
}

func MakeFail(c *gin.Context, code int, message interface{}) {
	c.JSON(http.StatusOK, gin.H{"statusCode": code, "message": message})
}
