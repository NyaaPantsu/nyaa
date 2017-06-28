package users

import (
	"errors"
	"net/http"

	"github.com/NyaaPantsu/nyaa/util/log"

	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/util/validator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/context"
	"golang.org/x/crypto/bcrypt"
)

// RetrieveUser retrieves a user.
func RetrieveFromRequest(c *gin.Context, id uint) (*models.User, bool, uint, int, error) {
	var user models.User
	var currentUserID uint
	var isAuthor bool

	if models.ORM.First(&user, id).RecordNotFound() {
		return nil, isAuthor, currentUserID, http.StatusNotFound, errors.New("user_not_found")
	}
	currentUser, err := CurrentUser(c)
	if err == nil {
		currentUserID = currentUser.ID
		isAuthor = currentUser.ID == user.ID
	}

	return &user, isAuthor, currentUserID, http.StatusOK, nil
}

// CurrentUser retrieves a current user.
func CurrentUser(c *gin.Context) (user *models.User, status int, err error) {
	encoded := c.Request.Header.Get("X-Auth-Token")
	if len(encoded) == 0 {
		// check cookie instead
		cookie, err = c.Cookie(CookieName)
		if err != nil {
			status = http.StatusInternalServerError
			return
		}
		encoded = cookie
	}
	userID, err = DecodeCookie(encoded)
	if err != nil {
		status = http.StatusInternalServerError
		return
	}

	userFromContext := getUserFromContext(c)

	if userFromContext.ID > 0 && userID == userFromContext.ID {
		user = &userFromContext
	} else {
		if ORM.Preload("Notifications").Where("user_id = ?", userID).First(user).RecordNotFound() { // We only load unread notifications
			return nil, user, errors.New("user_not_found")
		}
		setUserToContext(c, *user)
	}

	if user.IsBanned() {
		// recheck as user might've been banned in the meantime
		return nil, user, errors.New("account_banned")
	}
	if err != nil {
		status = http.StatusInternalServerError
		return
	}
	status = http.StatusOK
	return
}

// UpdateFromRequest updates a user.
func UpdateFromRequest(c *gin.Context, form *formStruct.UserForm, formSet *formStruct.UserSettingsForm, currentUser *models.User, id string) (user *models.User, status int) {
	messages := msg.GetMessages(c)
	if models.ORM.First(&user, id).RecordNotFound() {
		messages.AddErrorT("errors", "user_not_found")
		status = http.StatusNotFound
		return
	}

	if form.Password != "" {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.CurrentPassword))
		if err != nil && !currentUser.IsModerator() {
			messages.AddErrorT("errors", "incorrect_password")
			status = http.StatusInternalServerError
			return
		}
		newPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), 10)
		if err != nil {
			messages.AddErrorT("errors", "error_password_generating")
			status = http.StatusInternalServerError
			return
		}
		form.Password = string(newPassword)
	} else { // Then no change of password
		form.Password = user.Password
	}
	if !currentUser.IsModerator() { // We don't want users to be able to modify some fields
		form.Status = user.Status
		form.Username = user.Username
	}
	if form.Email != user.Email {
		// send verification to new email and keep old
		email.SendVerificationToUser(user, form.Email)
		form.Email = user.Email
	}
	log.Debugf("form %+v\n", form)
	validator.AssignValue(&user, form)

	// We set settings here
	user.ParseSettings()
	user.Settings.Set("new_torrent", formSet.NewTorrent)
	user.Settings.Set("new_torrent_email", formSet.NewTorrentEmail)
	user.Settings.Set("new_comment", formSet.NewComment)
	user.Settings.Set("new_comment_email", formSet.NewCommentEmail)
	user.Settings.Set("new_responses", formSet.NewResponses)
	user.Settings.Set("new_responses_email", formSet.NewResponsesEmail)
	user.Settings.Set("new_follower", formSet.NewFollower)
	user.Settings.Set("new_follower_email", formSet.NewFollowerEmail)
	user.Settings.Set("followed", formSet.Followed)
	user.Settings.Set("followed_email", formSet.FollowedEmail)
	user.SaveSettings()

	status, err = Update(user)
	if err != nil {
		messages.Error(err)
	}
	return
}

func getUserFromContext(c *gin.Context) models.User {
	if rv := context.Get(c.Request, UserContextKey); rv != nil {
		return rv.(models.User)
	}
	return models.User{}
}

func setUserToContext(c *gin.Context, val models.User) {
	context.Set(c.Request, UserContextKey, val)
}
