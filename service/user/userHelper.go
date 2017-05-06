package userService

import (
	"github.com/ewha/nyaa/db"
	"github.com/ewha/nyaa/model"
	//   "github.com/gin-gonic/gin"
	"errors"
	"net/http"

	"github.com/ewha/nyaa/util/log"
	"github.com/ewha/nyaa/util/retrieveHelper"
	//   "github.com/dorajistyle/goyangi/util/crypto"
)

// FindUserByUserName creates a user.
func FindUserByUserName(appID int64, userName string) (model.User, int, error) {
	var user model.User
	var err error
	// token := c.Request.Header.Get("X-Auth-Token")
	if db.ORM.Where("app_id=? and name=?", appID, userName).First(&user).RecordNotFound() {
		return user, http.StatusUnauthorized, err
	}
	return user, http.StatusOK, nil
}

// FindOrCreateUser creates a user.
func FindOrCreateUser(appID int64, userName string) (model.User, int, error) {
	var user model.User
	var err error

	// if len(token) > 0 {
	// 	log.Debug("header token exist.")
	// } else {
	// 	token, err = Token(c)
	// 	log.Debug("header token not exist.")
	// 	if err != nil {
	// 		return user, http.StatusUnauthorized, err
	// 	}
	// }
	log.Debugf("userName : %s\n", userName)
	// log.Debugf("Error : %s\n", err.Error())
	if db.ORM.Where("app_id=? and name=?", appID, userName).First(&user).RecordNotFound() {
		var user model.User
		// return user, http.StatusBadRequest, err
		user.Name = userName
		// user.Token = token
		user.AppID = appID
		log.Debugf("user %+v\n", user)
		if db.ORM.Create(&user).Error != nil {
			return user, http.StatusBadRequest, errors.New("User is not created.")
		}
		log.Debugf("retrived User %v\n", user)
		return user, http.StatusOK, nil
	}
	return user, http.StatusBadRequest, nil
}
