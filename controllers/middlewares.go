package controllers

import (
	"net/http"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/utils/log"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/gin-gonic/gin"
)

func errorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if config.Get().Environment == "DEVELOPMENT" {
			messages := msg.GetMessages(c)
			if messages.HasErrors() {
				log.Errorf("Request has errors: %v", messages.GetAllErrors())
			}
		}
		if c.Writer.Status() != http.StatusOK && c.Writer.Size() <= 0 {
			if c.ContentType() == "application/json" {
				messages := msg.GetMessages(c)
				messages.AddErrorT("errors", "404_not_found")
				c.JSON(c.Writer.Status(), messages.GetAllErrors())
				return
			}
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

func pprofHandler(handler http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser := getUser(c)
		if currentUser.HasAdmin() {
			handler.ServeHTTP(c.Writer, c.Request)
		} else {
			httpError(c, http.StatusNotFound)
		}
	}
}
