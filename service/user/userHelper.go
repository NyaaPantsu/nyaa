package userService

import (
	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
	"errors"
	"net/http"

	"github.com/ewhal/nyaa/util/log"
)

// FindUserByUserName creates a user.
func FindUserByUserName(userName string) (model.User, int, error) {
	var user model.User
	var err error
	if db.ORM.Where("name=?", userName).First(&user).RecordNotFound() {
		return user, http.StatusUnauthorized, err
	}
	return user, http.StatusOK, nil
}

// FindOrCreateUser creates a user.
func FindOrCreateUser(username string) (model.User, int, error) {
	var user model.User
	if db.ORM.Where("username=?", username).First(&user).RecordNotFound() {
		var user model.User
		user.Username = username
		log.Debugf("user %+v\n", user)
		if db.ORM.Create(&user).Error != nil {
			return user, http.StatusBadRequest, errors.New("User is not created.")
		}
		log.Debugf("retrived User %v\n", user)
		return user, http.StatusOK, nil
	}
	return user, http.StatusBadRequest, nil
}
