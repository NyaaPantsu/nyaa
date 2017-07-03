package users

import (
	"errors"
	"net/http"

	"github.com/NyaaPantsu/nyaa/utils/log"

	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/utils/validator"
	"github.com/NyaaPantsu/nyaa/utils/validator/user"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// UpdateFromRequest updates a user.
func UpdateFromRequest(c *gin.Context, form *userValidator.UserForm, formSet *userValidator.UserSettingsForm, currentUser *models.User, id uint) (*models.User, int, error) {
	var user = &models.User{}
	if models.ORM.First(user, id).RecordNotFound() {
		return user, http.StatusNotFound, errors.New("user_not_found")
	}

	if !currentUser.IsModerator() { // We don't want users to be able to modify some fields
		form.Status = user.Status
		form.Username = user.Username
	}
	if form.Password != "" {
		errBcrypt := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.CurrentPassword))
		if errBcrypt != nil && !currentUser.IsModerator() {
			return user, http.StatusInternalServerError, errors.New("incorrect_password")
		}
		newPassword, errBcrypt := bcrypt.GenerateFromPassword([]byte(form.Password), 10)
		if errBcrypt != nil {
			return user, http.StatusInternalServerError, errors.New("error_password_generating")
		}
		form.Password = string(newPassword)
	} else { // Then no change of password
		form.Password = user.Password
	}
	log.Debugf("form %+v\n", form)
	validator.Bind(user, form)

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

	status, err := user.Update()
	return user, status, err
}
