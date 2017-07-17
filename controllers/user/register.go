package userController

import (
	"net/http"

	"github.com/NyaaPantsu/nyaa/models/users"
	"github.com/NyaaPantsu/nyaa/templates"
	"github.com/NyaaPantsu/nyaa/utils/captcha"
	"github.com/NyaaPantsu/nyaa/utils/cookies"
	"github.com/NyaaPantsu/nyaa/utils/email"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/validator"
	"github.com/NyaaPantsu/nyaa/utils/validator/user"
	"github.com/gin-gonic/gin"
)

// UserRegisterFormHandler : Getting View User Registration
func UserRegisterFormHandler(c *gin.Context) {
	_, _, errorUser := cookies.CurrentUser(c)
	// User is already connected, redirect to home
	if errorUser == nil {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
	registrationForm := userValidator.RegistrationForm{}
	c.Bind(&registrationForm)
	registrationForm.CaptchaID = captcha.GetID()
	templates.Form(c, "site/user/register.jet.html", registrationForm)
}

// UserRegisterPostHandler : Post Registration controller, we do some check on the form here, the rest on user service
func UserRegisterPostHandler(c *gin.Context) {
	b := userValidator.RegistrationForm{}
	messages := msg.GetMessages(c)

	if !captcha.Authenticate(captcha.Extract(c)) {
		messages.AddErrorT("errors", "bad_captcha")
	}
	if !messages.HasErrors() {
		if len(c.PostForm("email")) > 0 {
			if !userValidator.EmailValidation(c.PostForm("email")) {
				messages.AddErrorT("email", "email_not_valid")
			}
		}
		if !userValidator.ValidateUsername(c.PostForm("username")) {
			messages.AddErrorT("username", "username_illegal")
		}

		if !messages.HasErrors() {
			c.Bind(&b)
			validator.ValidateForm(&b, messages)
			if !messages.HasErrors() {
				user, _ := users.CreateUser(c)
				if !messages.HasErrors() {
					_, err := cookies.SetLogin(c, user)
					if err != nil {
						messages.Error(err)
					}
					if b.Email != "" {
						email.SendVerificationToUser(user, b.Email)
					}
					if !messages.HasErrors() {
						templates.Static(c, "site/static/signup_success.jet.html")
					}
				}
			}
		}
	}
	if messages.HasErrors() {
		UserRegisterFormHandler(c)
	}
}
