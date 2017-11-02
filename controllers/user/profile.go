package userController

import (
	"strconv"
	"time"
	"fmt"

	"net/http"

	"github.com/NyaaPantsu/nyaa/controllers/router"
	"github.com/NyaaPantsu/nyaa/models/notifications"
	"github.com/NyaaPantsu/nyaa/models/users"
	"github.com/NyaaPantsu/nyaa/templates"
	"github.com/NyaaPantsu/nyaa/utils/cookies"
	"github.com/NyaaPantsu/nyaa/utils/crypto"
	"github.com/NyaaPantsu/nyaa/utils/email"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/publicSettings"
	"github.com/NyaaPantsu/nyaa/utils/validator"
	"github.com/NyaaPantsu/nyaa/utils/validator/user"
	"github.com/gin-gonic/gin"
)

// UserProfileDelete :  Deleting User Profile
func UserProfileDelete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	userProfile, _, errorUser := users.FindForAdmin(uint(id))
	if errorUser == nil{
		currentUser := router.GetUser(c)
		if (currentUser.CurrentOrAdmin(userProfile.ID)) {
			_, err := userProfile.Delete(currentUser)
			if err == nil && currentUser.CurrentUserIdentical(userProfile.ID) {
				cookies.Clear(c)
			}
		}
		templates.Static(c, "site/static/delete_success.jet.html")
	}
}

// UserProfileHandler :  Getting User Profile
func UserProfileHandler(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	Ts, _ := publicSettings.GetTfuncAndLanguageFromRequest(c)
	messages := msg.GetMessages(c)

	if c.Param("id") != "0" && id == 0 && ContainsNonNumbersChars(c.Param("id")) {
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/username/%s", c.Param("id")))
		return
	}
	
	userProfile, _, errorUser := users.FindForAdmin(uint(id))
	if errorUser == nil {
		currentUser := router.GetUser(c)
		follow := c.Request.URL.Query()["followed"]
		unfollow := c.Request.URL.Query()["unfollowed"]
		deleteVar := c.Request.URL.Query()["delete"]

		if !((deleteVar != nil) && (currentUser.CurrentOrAdmin(userProfile.ID))) {
			if follow != nil {
				messages.AddInfof("infos", Ts("user_followed_msg"), userProfile.Username)
			}
			if unfollow != nil {
				messages.AddInfof("infos", Ts("user_unfollowed_msg"), userProfile.Username)
			}
			userProfile.ParseSettings()

			templates.UserProfile(c, userProfile)
		}
	} else {
		variables := templates.Commonvariables(c)
		templates.Render(c, "errors/user_not_found.jet.html", variables)
	}
}

func ContainsNonNumbersChars(source string) bool {
	for char := range source {
		if char < 30 || char > 39 {
			return true
		}
	}
	return false
}

func UserGetFromName(c *gin.Context) {
	username := c.Param("username")
	
	Ts, _ := publicSettings.GetTfuncAndLanguageFromRequest(c)
	messages := msg.GetMessages(c)

	userProfile, _, _, err := users.FindByUsername(username)
	if err == nil {
		currentUser := router.GetUser(c)
		follow := c.Request.URL.Query()["followed"]
		unfollow := c.Request.URL.Query()["unfollowed"]
		deleteVar := c.Request.URL.Query()["delete"]

		if (deleteVar != nil) && (currentUser.CurrentOrAdmin(userProfile.ID)) {
			_, err := userProfile.Delete(currentUser)
			if err == nil && currentUser.CurrentUserIdentical(userProfile.ID) {
				cookies.Clear(c)
			}
			templates.Static(c, "site/static/delete_success.jet.html")
		} else {
			if follow != nil {
				messages.AddInfof("infos", Ts("user_followed_msg"), userProfile.Username)
			}
			if unfollow != nil {
				messages.AddInfof("infos", Ts("user_unfollowed_msg"), userProfile.Username)
			}
			userProfile.ParseSettings()

			templates.UserProfile(c, userProfile)
		}
	} else {
		variables := templates.Commonvariables(c)
		searchForm := templates.NewSearchForm(c)
		searchForm.User = username
		variables.Set("Search", searchForm)
		templates.Render(c, "errors/user_not_found.jet.html", variables)
	}
}

func RedirectToUserSearch(c *gin.Context) {
	username := c.Query("username")
	
	if username == "" {
		variables := templates.Commonvariables(c)
		templates.Render(c, "errors/user_not_found.jet.html", variables)
	} else {
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/username/%s", username))
	}
}

// UserDetailsHandler : Getting User Profile Details View
func UserDetailsHandler(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	currentUser := router.GetUser(c)

	userProfile, _, errorUser := users.FindForAdmin(uint(id))
	if errorUser == nil && currentUser.CurrentOrAdmin(userProfile.ID) {
		b := userValidator.UserForm{}
		c.Bind(&b)
		availableLanguages := publicSettings.GetAvailableLanguages()
		userProfile.ParseSettings()
		templates.UserProfileEdit(c, userProfile, b, availableLanguages)
	} else {
		variables := templates.Commonvariables(c)
		templates.Render(c, "errors/user_not_found.jet.html", variables)
	}
}

// UserProfileFormHandler : Getting View User Profile Update
func UserProfileFormHandler(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	currentUser := router.GetUser(c)
	userProfile, _, errorUser := users.FindForAdmin(uint(id))
	if errorUser != nil || !currentUser.CurrentOrAdmin(userProfile.ID) || userProfile.ID == 0 {
		c.Status(http.StatusNotFound)
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
				if currentUser.HasAdmin() {
					userProfile.Email = userForm.Email
				} else {
					email.SendVerificationToUser(currentUser, userForm.Email)
					messages.AddInfoTf("infos", "email_changed", userForm.Email)
					userForm.Email = userProfile.Email // reset, it will be set when user clicks verification
				}
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
				userProfile = user
			}
		}
	}
	availableLanguages := publicSettings.GetAvailableLanguages()
	templates.UserProfileEdit(c, userProfile, userForm, availableLanguages)
}

// UserNotificationsHandler : Controller to show user notifications
func UserNotificationsHandler(c *gin.Context) {
	currentUser := router.GetUser(c)
	if currentUser.ID > 0 {
		if c.Request.URL.Query()["clear"] != nil {
			notifications.DeleteNotifications(currentUser, false)
			
		} else if c.Request.URL.Query()["clear_all"] != nil {
			notifications.DeleteNotifications(currentUser, true)
		} else if c.Request.URL.Query()["read_all"] != nil {
			notifications.MarkAllNotificationsAsRead(currentUser)
		}
		templates.UserProfileNotifications(c, currentUser)
	} else {
		c.Status(http.StatusNotFound)
	}
}

// UserAPIKeyResetHandler : Controller to reset user api key
func UserAPIKeyResetHandler(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	currentUser := router.GetUser(c)

	messages := msg.GetMessages(c)
	userProfile, _, errorUser := users.FindForAdmin(uint(id))
	if errorUser != nil || !currentUser.CurrentOrAdmin(userProfile.ID) || userProfile.ID == 0 {
		c.Status(http.StatusNotFound)
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
	UserDetailsHandler(c)
}
