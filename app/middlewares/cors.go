package middlewares

import "github.com/gin-gonic/gin"

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Origin, Content-Type, Accept, token")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 对于预检请求，直接返回 200
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}
