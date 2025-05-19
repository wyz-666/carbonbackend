package middlewares

import (
	"wxcloudrun-golang/app/handlers/response"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

// CheckPermission
func CheckPermission() func(c *gin.Context) {
	return func(c *gin.Context) {
		u, ok := c.Get("userID")
		userID := u.(string)
		if !ok {
			response.MakeFail(c, http.StatusForbidden, "userID not found!")
			c.Abort()
			return
		}

		t, ok := c.Get("userType")
		userType := t.(int)
		if !ok {
			response.MakeFail(c, http.StatusForbidden, "userType not found!")
			c.Abort()
			return
		}
		result := false
		if userID[4] == 'B' && userType == 3 {
			result = true
		}
		if userID[4] == 'A' && userType == 2 {
			result = true
		}
		if userID[4] == '0' && userType == 1 {
			result = true
		}
		if result {
			glog.Infoln(userID, " pass!")
			c.Next()
		} else {
			glog.Infoln(userID, " deny!")
			response.MakeFail(c, http.StatusForbidden, "Not Permission")
			c.Abort()
		}
	}
}
