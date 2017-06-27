package router

import (
	"net/http"
	"strconv"
	"time"

	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/service/captcha"
	"github.com/NyaaPantsu/nyaa/service/notifier"
	"github.com/NyaaPantsu/nyaa/service/user"
	"github.com/NyaaPantsu/nyaa/service/user/form"
	"github.com/NyaaPantsu/nyaa/service/user/permission"
	"github.com/NyaaPantsu/nyaa/util/crypto"
	msg "github.com/NyaaPantsu/nyaa/util/messages"
	"github.com/NyaaPantsu/nyaa/util/modelHelper"
	"github.com/NyaaPantsu/nyaa/util/publicSettings"
	"github.com/NyaaPantsu/nyaa/util/search"
	"github.com/gin-gonic/gin"
)

// UserRegisterFormHandler : Getting View User Registration
func UserRegisterFormHandler(c *gin.Context) {
	_, errorUser := userService.CurrentUser(c)
	// User is already connected, redirect to home
	if errorUser == nil {
		SearchHandler(c)
		return
	}
	registrationForm := form.RegistrationForm{}
	c.Bind(&registrationForm)
	registrationForm.CaptchaID = captcha.GetID()
	formTemplate(c, "user/register.jet.html", registrationForm)
}

// UserLoginFormHandler : Getting View User Login
func UserLoginFormHandler(c *gin.Context) {
	_, errorUser := userService.CurrentUser(c)
	// User is already connected, redirect to home
	if errorUser == nil {
		SearchHandler(c)
		return
	}

	loginForm := form.LoginForm{}
	formTemplate(c, "user/register.jet.html", loginForm)
}

// UserProfileHandler :  Getting User Profile
func UserProfileHandler(c *gin.Context) {
	id := c.Query("id")
	Ts, _ := publicSettings.GetTfuncAndLanguageFromRequest(c)
	messages := msg.GetMessages(c)

	userProfile, _, errorUser := userService.RetrieveUserForAdmin(id)
	if errorUser == nil {
		currentUser := getUser(c)
		follow := c.Request.URL.Query()["followed"]
		unfollow := c.Request.URL.Query()["unfollowed"]
		deleteVar := c.Request.URL.Query()["delete"]

		if (deleteVar != nil) && (userPermission.CurrentOrAdmin(currentUser, userProfile.ID)) {
			_ = userService.DeleteUser(c, currentUser, id)
			staticTemplate(c, "user/delete_success.jet.html")
		} else {
			if follow != nil {
				messages.AddInfof("infos", Ts("user_followed_msg"), userProfile.Username)
			}
			if unfollow != nil {
				messages.AddInfof("infos", Ts("user_unfollowed_msg"), userProfile.Username)
			}
			userProfile.ParseSettings()
			query := c.Request.URL.Query()
			query.Set("userID", id)
			query.Set("max", "16")
			c.Request.URL.RawQuery = query.Encode()
			var torrents []model.Torrent
			var err error
			if userPermission.CurrentOrAdmin(currentUser, userProfile.ID) {
				_, torrents, _, err = search.SearchByQuery(c, 1)
			} else {
				_, torrents, _, err = search.SearchByQueryNoHidden(c, 1)
			}
			if err != nil {
				messages.AddErrorT("errors", "retrieve_torrents_error")
			}
			userProfile.Torrents = torrents
			userProfileTemplate(c, &userProfile)
		}
	} else {
		NotFoundHandler(c)
	}
}

// UserDetailsHandler : Getting User Profile Details View
func UserDetailsHandler(c *gin.Context) {
	id := c.Query("id")
	currentUser := getUser(c)

	userProfile, _, errorUser := userService.RetrieveUserForAdmin(id)
	if errorUser == nil && userPermission.CurrentOrAdmin(currentUser, userProfile.ID) {
		if userPermission.CurrentOrAdmin(currentUser, userProfile.ID) {
			b := form.UserForm{}
			c.Bind(&b)
			availableLanguages := publicSettings.GetAvailableLanguages()
			userProfile.ParseSettings()
			userProfileEditTemplate(c, &userProfile, b, availableLanguages)
		}
	} else {
		NotFoundHandler(c)
	}
}

// UserProfileFormHandler : Getting View User Profile Update
func UserProfileFormHandler(c *gin.Context) {
	id := c.Query("id")
	currentUser := getUser(c)
	userProfile, _, errorUser := userService.RetrieveUserForAdmin(id)
	if errorUser != nil || !userPermission.CurrentOrAdmin(currentUser, userProfile.ID) || userProfile.ID == 0 {
		NotFoundHandler(c)
		return
	}
	userProfile.ParseSettings()
	messages := msg.GetMessages(c)
	userForm := form.UserForm{}
	userSettingsForm := form.UserSettingsForm{}

	if len(c.PostForm("email")) > 0 {
		form.EmailValidation(c.PostForm("email"), messages)
	}
	if len(c.PostForm("username")) > 0 {
		form.ValidateUsername(c.PostForm("username"), messages)
	}

	if !messages.HasErrors() {
		c.Bind(&userForm)
		c.Bind(&userSettingsForm)
		if !userPermission.HasAdmin(currentUser) {
			userForm.Username = userProfile.Username
			userForm.Status = userProfile.Status
		} else {
			if userProfile.Status != userForm.Status && userForm.Status == 2 {
				messages.AddErrorT("errors", "elevating_user_error")
			}
		}
		modelHelper.ValidateForm(&userForm, messages)
		if !messages.HasErrors() {
			if userForm.Email != userProfile.Email {
				userService.SendVerificationToUser(*currentUser, userForm.Email)
				messages.AddInfoTf("infos", "email_changed", userForm.Email)
				userForm.Email = userProfile.Email // reset, it will be set when user clicks verification
			}
			userProfile, _ = userService.UpdateUser(c, &userForm, &userSettingsForm, currentUser, id)
			if !messages.HasErrors() {
				messages.AddInfoT("infos", "profile_updated")
			}
		}
	}
	availableLanguages := publicSettings.GetAvailableLanguages()
	userProfileEditTemplate(c, &userProfile, userForm, availableLanguages)
}

// UserRegisterPostHandler : Post Registration controller, we do some check on the form here, the rest on user service
func UserRegisterPostHandler(c *gin.Context) {
	b := form.RegistrationForm{}
	messages := msg.GetMessages(c)

	if !captcha.Authenticate(captcha.Extract(c)) {
		messages.AddErrorT("errors", "bad_captcha")
	}
	if !messages.HasErrors() {
		if len(c.PostForm("email")) > 0 {
			form.EmailValidation(c.PostForm("email"), messages)
		}
		form.ValidateUsername(c.PostForm("username"), messages)
		if !messages.HasErrors() {
			c.Bind(&b)
			modelHelper.ValidateForm(&b, messages)
			if !messages.HasErrors() {
				_ = userService.CreateUser(c)
				if !messages.HasErrors() {
					staticTemplate(c, "user/signup_success.jet.html")
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
	token := c.Query("token")
	messages := msg.GetMessages(c)

	_, errEmail := userService.EmailVerification(token, c)
	if errEmail != nil {
		messages.ImportFromError("errors", errEmail)
	}
	staticTemplate(c, "user/verify_success.jet.html")
}

// UserLoginPostHandler : Post Login controller
func UserLoginPostHandler(c *gin.Context) {
	b := form.LoginForm{}
	c.Bind(&b)
	messages := msg.GetMessages(c)

	modelHelper.ValidateForm(&b, messages)
	if !messages.HasErrors() {
		_, errorUser := userService.CreateUserAuthentication(c)
		if errorUser == nil {
			c.Redirect(http.StatusSeeOther, "/")
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
		userService.ClearCookie(c)
		url := c.DefaultPostForm("redirectTo", "/")
		c.Redirect(http.StatusSeeOther, url)
	} else {
		NotFoundHandler(c)
	}
}

// UserFollowHandler : Controller to follow/unfollow users, need user id to follow
func UserFollowHandler(c *gin.Context) {
	var followAction string
	id := c.Query("id")
	currentUser := getUser(c)
	user, _, errorUser := userService.RetrieveUserForAdmin(id)
	if errorUser == nil && user.ID > 0 {
		if !userPermission.IsFollower(&user, currentUser) {
			followAction = "followed"
			userService.SetFollow(&user, currentUser)
		} else {
			followAction = "unfollowed"
			userService.RemoveFollow(&user, currentUser)
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
			notifierService.DeleteAllNotifications(currentUser.ID)
			messages.AddInfoT("infos", "notifications_cleared")
			currentUser.Notifications = []model.Notification{}
		}
		userProfileTemplate(c, currentUser)
	} else {
		NotFoundHandler(c)
	}
}

// UserAPIKeyResetHandler : Controller to reset user api key
func UserAPIKeyResetHandler(c *gin.Context) {
	id := c.Query("id")
	currentUser := getUser(c)

	messages := msg.GetMessages(c)
	userProfile, _, errorUser := userService.RetrieveUserForAdmin(id)
	if errorUser != nil || !userPermission.CurrentOrAdmin(currentUser, userProfile.ID) || userProfile.ID == 0 {
		NotFoundHandler(c)
		return
	}
	userProfile.APIToken, _ = crypto.GenerateRandomToken32()
	userProfile.APITokenExpiry = time.Unix(0, 0)
	_, errorUser = userService.UpdateRawUser(&userProfile)
	if errorUser != nil {
		messages.Error(errorUser)
	} else {
		messages.AddInfoT("infos", "profile_updated")
	}
	UserProfileHandler(c)
}
