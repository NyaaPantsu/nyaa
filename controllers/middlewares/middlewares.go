package middlewares

import (
	"net/http"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/controllers/router"
	"github.com/NyaaPantsu/nyaa/templates"
	"github.com/NyaaPantsu/nyaa/utils/log"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/oauth2"
	"github.com/gin-gonic/gin"
	"github.com/ory/fosite"
)

// ErrorMiddleware for managing errors on status
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.Writer.Status() >= 300 && config.Get().Environment == "DEVELOPMENT" {
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

func ScopesRequired(scopes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		mySessionData := oauth2.NewSession("", "")
		ctx, err := oauth2.Oauth2.IntrospectToken(c, fosite.AccessTokenFromRequest(c.Request), fosite.AccessToken, mySessionData, scopes...)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
		// All required scopes are found
		c.Set("fosite", ctx)
		c.Next()
	}
}

// CSP set Content Security Policy http header
func CSP() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Security-Policy", "default-src 'self'; img-src * data:; media-src *; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline'")
		c.Next()
	}
}
