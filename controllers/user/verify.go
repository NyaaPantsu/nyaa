package userController

import (
	"github.com/NyaaPantsu/nyaa/templates"
	"github.com/NyaaPantsu/nyaa/utils/email"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/gin-gonic/gin"
)

// UserVerifyEmailHandler : Controller when verifying email, needs a token
func UserVerifyEmailHandler(c *gin.Context) {
	token := c.Param("token")
	messages := msg.GetMessages(c)

	_, errEmail := email.EmailVerification(token, c)
	if errEmail != nil {
		messages.ImportFromError("errors", errEmail)
	}
	templates.Static(c, "site/static/verify_success.jet.html")
}
