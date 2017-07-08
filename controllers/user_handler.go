package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/notifications"
	"github.com/NyaaPantsu/nyaa/models/users"
	"github.com/NyaaPantsu/nyaa/utils/captcha"
	"github.com/NyaaPantsu/nyaa/utils/cookies"
	"github.com/NyaaPantsu/nyaa/utils/crypto"
	"github.com/NyaaPantsu/nyaa/utils/email"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/publicSettings"
	"github.com/NyaaPantsu/nyaa/utils/search"
	"github.com/NyaaPantsu/nyaa/utils/validator"
	"github.com/NyaaPantsu/nyaa/utils/validator/user"
	"github.com/gin-gonic/gin"
)

// UserRegisterFormHandler : Getting View User Registration
func UserRegisterFormHandler(c *gin.Context) {
	_, _, errorUser := cookies.CurrentUser(c)
	// User is already connected, redirect to home
	if errorUser == nil {
		SearchHandler(c)
		return
	}
	registrationForm := userValidator.RegistrationForm{}
	c.Bind(&registrationForm)
	registrationForm.CaptchaID = captcha.GetID()
	formTemplate(c, "site/user/register.jet.html", registrationForm)
}

// UserLoginFormHandler : Getting View User Login
func UserLoginFormHandler(c *gin.Context) {
	_, _, errorUser := cookies.CurrentUser(c)
	// User is already connected, redirect to home
	if errorUser == nil {
		SearchHandler(c)
		return
	}

	loginForm := userValidator.LoginForm{
		RedirectTo: c.DefaultQuery("redirectTo", ""),
	}
	formTemplate(c, "site/user/login.jet.html", loginForm)
}

// UserProfileHandler :  Getting User Profile
func UserProfileHandler(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	fmt.Printf("User ID: %s", id)
	Ts, _ := publicSettings.GetTfuncAndLanguageFromRequest(c)
	messages := msg.GetMessages(c)

	userProfile, _, errorUser := users.FindForAdmin(uint(id))
	if errorUser == nil {
		currentUser := getUser(c)
		follow := c.Request.URL.Query()["followed"]
		unfollow := c.Request.URL.Query()["unfollowed"]
		deleteVar := c.Request.URL.Query()["delete"]

		if (deleteVar != nil) && (currentUser.CurrentOrAdmin(userProfile.ID)) {
			_, err := userProfile.Delete(currentUser)
			if err == nil && currentUser.CurrentUserIdentical(userProfile.ID) {
				cookies.Clear(c)
			}
			staticTemplate(c, "site/static/delete_success.jet.html")
		} else {
			if follow != nil {
				messages.AddInfof("infos", Ts("user_followed_msg"), userProfile.Username)
			}
			if unfollow != nil {
				messages.AddInfof("infos", Ts("user_unfollowed_msg"), userProfile.Username)
			}
			userProfile.ParseSettings()
			query := c.Request.URL.Query()
			query.Set("userID", strconv.Itoa(int(id)))
			query.Set("max", "16")
			c.Request.URL.RawQuery = query.Encode()
			var torrents []models.Torrent
			var err error
			if currentUser.CurrentOrAdmin(userProfile.ID) {
				_, torrents, _, err = search.ByQuery(c, 1)
			} else {
				_, torrents, _, err = search.ByQueryNoHidden(c, 1)
			}
			if err != nil {
				messages.AddErrorT("errors", "retrieve_torrents_error")
			}
			userProfile.Torrents = torrents
			userProfileTemplate(c, userProfile)
		}
	} else {
		NotFoundHandler(c)
	}
}

// UserDetailsHandler : Getting User Profile Details View
func UserDetailsHandler(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	currentUser := getUser(c)

	userProfile, _, errorUser := users.FindForAdmin(uint(id))
	if errorUser == nil && currentUser.CurrentOrAdmin(userProfile.ID) {
		b := userValidator.UserForm{}
		c.Bind(&b)
		availableLanguages := publicSettings.GetAvailableLanguages()
		userProfile.ParseSettings()
		userProfileEditTemplate(c, userProfile, b, availableLanguages)
	} else {
		NotFoundHandler(c)
	}
}

// UserProfileFormHandler : Getting View User Profile Update
func UserProfileFormHandler(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	currentUser := getUser(c)
	userProfile, _, errorUser := users.FindForAdmin(uint(id))
	if errorUser != nil || !currentUser.CurrentOrAdmin(userProfile.ID) || userProfile.ID == 0 {
		NotFoundHandler(c)
		return
	}
	userProfile.ParseSettings()
	messages := msg.GetMessages(c)
	userForm := userValidator.UserForm{}
	userSettingsForm := userValidator.UserSettingsForm{}

	if len(c.PostForm("email")) > 0 {
		if !userValidator.EmailValidation(c.PostForm("email")) {
			messages.AddErrorT("email", "email_not_valid")
		}
	}
	if len(c.PostForm("username")) > 0 {
		if !userValidator.ValidateUsername(c.PostForm("username")) {
			messages.AddErrorT("username", "username_illegal")
		}
	}

	if !messages.HasErrors() {
		c.Bind(&userForm)
		c.Bind(&userSettingsForm)
		if !currentUser.HasAdmin() {
			userForm.Username = userProfile.Username
			userForm.Status = userProfile.Status
		} else {
			if userProfile.Status != userForm.Status && userForm.Status == 2 {
				messages.AddErrorT("errors", "elevating_user_error")
			}
		}
		validator.ValidateForm(&userForm, messages)
		if !messages.HasErrors() {
			if userForm.Email != userProfile.Email {
				email.SendVerificationToUser(currentUser, userForm.Email)
				messages.AddInfoTf("infos", "email_changed", userForm.Email)
				userForm.Email = userProfile.Email // reset, it will be set when user clicks verification
			}
			user, _, err := users.UpdateFromRequest(c, &userForm, &userSettingsForm, currentUser, uint(id))
			if err != nil {
				messages.Error(err)
			}
			if userForm.Email != user.Email {
				// send verification to new email and keep old
				email.SendVerificationToUser(user, userForm.Email)
			}
			if !messages.HasErrors() {
				messages.AddInfoT("infos", "profile_updated")
			}
		}
	}
	availableLanguages := publicSettings.GetAvailableLanguages()
	userProfileEditTemplate(c, userProfile, userForm, availableLanguages)
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
				_, err := cookies.SetLogin(c, user)
				if err != nil {
					messages.Error(err)
				}
				if b.Email != "" {
					email.SendVerificationToUser(user, b.Email)
				}
				if !messages.HasErrors() {
					staticTemplate(c, "site/static/signup_success.jet.html")
				}
			}
		}
	}
	if messages.HasErrors() {
		UserRegisterFormHandler(c)
	}
}

// UserVerifyEmailHandler : Controller when verifying email, needs a token
func UserVerifyEmailHandler(c *gin.Context) {
	token := c.Param("token")
	messages := msg.GetMessages(c)

	_, errEmail := email.EmailVerification(token, c)
	if errEmail != nil {
		messages.ImportFromError("errors", errEmail)
	}
	staticTemplate(c, "site/static/verify_success.jet.html")
}

// UserLoginPostHandler : Post Login controller
func UserLoginPostHandler(c *gin.Context) {
	b := userValidator.LoginForm{}
	c.Bind(&b)
	messages := msg.GetMessages(c)

	validator.ValidateForm(&b, messages)
	if !messages.HasErrors() {
		_, _, errorUser := cookies.CreateUserAuthentication(c, &b)
		if errorUser == nil {
			url := c.DefaultPostForm("redirectTo", "/")
			c.Redirect(http.StatusSeeOther, url)
			return
		}
		messages.ErrorT(errorUser)
	}
	UserLoginFormHandler(c)
}

// UserLogoutHandler : Controller to logout users
func UserLogoutHandler(c *gin.Context) {
	logout := c.PostForm("logout")
	if logout != "" {
		cookies.Clear(c)
		url := c.DefaultPostForm("redirectTo", "/")
		c.Redirect(http.StatusSeeOther, url)
	} else {
		NotFoundHandler(c)
	}
}

// UserFollowHandler : Controller to follow/unfollow users, need user id to follow
func UserFollowHandler(c *gin.Context) {
	var followAction string
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	currentUser := getUser(c)
	user, _, errorUser := users.FindForAdmin(uint(id))
	if errorUser == nil && user.ID > 0 {
		if !currentUser.IsFollower(user) {
			followAction = "followed"
			currentUser.SetFollow(user)
		} else {
			followAction = "unfollowed"
			currentUser.RemoveFollow(user)
		}
	}
	url := "/user/" + strconv.Itoa(int(user.ID)) + "/" + user.Username + "?" + followAction
	c.Redirect(http.StatusSeeOther, url)
}

// UserNotificationsHandler : Controller to show user notifications
func UserNotificationsHandler(c *gin.Context) {
	currentUser := getUser(c)
	if currentUser.ID > 0 {
		messages := msg.GetMessages(c)
		if c.Request.URL.Query()["clear"] != nil {
			notifications.DeleteAllNotifications(currentUser.ID)
			messages.AddInfoT("infos", "notifications_cleared")
			currentUser.Notifications = []models.Notification{}
		}
		userProfileNotificationsTemplate(c, currentUser)
	} else {
		NotFoundHandler(c)
	}
}

// UserAPIKeyResetHandler : Controller to reset user api key
func UserAPIKeyResetHandler(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	currentUser := getUser(c)

	messages := msg.GetMessages(c)
	userProfile, _, errorUser := users.FindForAdmin(uint(id))
	if errorUser != nil || !currentUser.CurrentOrAdmin(userProfile.ID) || userProfile.ID == 0 {
		NotFoundHandler(c)
		return
	}
	userProfile.APIToken, _ = crypto.GenerateRandomToken32()
	userProfile.APITokenExpiry = time.Unix(0, 0)
	_, errorUser = userProfile.UpdateRaw()
	if errorUser != nil {
		messages.Error(errorUser)
	} else {
		messages.AddInfoT("infos", "profile_updated")
	}
	UserProfileHandler(c)
}
