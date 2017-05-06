package userPermission

import (
	"errors"
	"net/http"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/service/user"
	"github.com/ewhal/nyaa/util/log"
)

// HasAdmin checks that user has an admin permission.
func HasAdmin(user *model.User) bool {
	name := "admin"
	for _, role := range user.Roles {
		log.Debugf("HasAdmin role.Name : %s", role.Name)
		if role.Name == name {
			return true
		}
	}
	return false
}

// CurrentOrAdmin check that user has admin permission or user is the current user.
func CurrentOrAdmin(user *model.User, userId uint) bool {
	log.Debugf("user.Id == userId %d %d %s", user.Id, userId, user.Id == userId)
	return (HasAdmin(user) || user.Id == userId)
}

// CurrentUserIdentical check that userId is same as current user's Id.
func CurrentUserIdentical(r *http.Request, userId uint) (bool, error) {
	currentUser, err := userService.CurrentUser(r)
	if err != nil {
		return false, errors.New("Auth failed.")
	}
	if currentUser.Id != userId {
		return false, errors.New("User is not identical.")
	}

	return true, nil
}
