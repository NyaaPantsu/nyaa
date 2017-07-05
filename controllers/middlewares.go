package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func errorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.Writer.Status() != http.StatusOK && c.Writer.Size() <= 0 {
			httpError(c, c.Writer.Status())
		}
	}
}

// Make sure the user is a moderator, otherwise return forbidden
func modMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser := getUser(c)
		if !currentUser.HasAdmin() {
			NotFoundHandler(c)
		}
		c.Next()
	}
}
