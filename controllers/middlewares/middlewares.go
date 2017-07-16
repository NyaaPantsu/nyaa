package middlewares

import (
	"net/http"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/controllers/router"
	"github.com/NyaaPantsu/nyaa/templates"
	"github.com/NyaaPantsu/nyaa/utils/log"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/gin-gonic/gin"
)

// Middleware for managing errors on status
func ErrorMiddleware() gin.HandlerFunc {
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
			templates.HttpError(c, c.Writer.Status())
		}
	}
}

// ModMiddleware Make sure the user is a moderator, otherwise return forbidden
func ModMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser := router.GetUser(c)
		if !currentUser.HasAdmin() {
			NotFoundHandler(c)
		}
		c.Next()
	}
}

func pprofHandler(handler http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser := router.GetUser(c)
		if currentUser.HasAdmin() {
			handler.ServeHTTP(c.Writer, c.Request)
		} else {
			templates.HttpError(c, http.StatusNotFound)
		}
	}
}
