package controllers

import (
	"net/http"

	"github.com/NyaaPantsu/nyaa/util/cookies"
	"github.com/gin-gonic/gin"
)

func errorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.Writer.Status() == http.StatusNotFound && c.Writer.Size() == 0 {
			NotFoundHandler(c)
		}
	}
}

// Make sure the user is a moderator, otherwise return forbidden
func modMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser := getUser(c)
		if !userPermission.HasAdmin(currentUser) {
			NotFoundHandler(c)
		}
		c.Next()
	}
}
