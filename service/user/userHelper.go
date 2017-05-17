package userService

import (
	"errors"
	"github.com/NyaaPantsu/nyaa/db"
	"github.com/NyaaPantsu/nyaa/model"
	"net/http"

	"github.com/NyaaPantsu/nyaa/util/log"
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
		var newUser model.User
		newUser.Username = username
		log.Debugf("user %+v\n", newUser)
		if db.ORM.Create(&newUser).Error != nil {
			return newUser, http.StatusBadRequest, errors.New("user not created")
		}
		log.Debugf("retrieved User %v\n", newUser)
		return newUser, http.StatusOK, nil
	}
	return user, http.StatusBadRequest, nil
}
