package userPermission

import (
	"github.com/ewhal/nyaa/model"
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
func CurrentUserIdentical(user *model.User, userId uint) (bool) {
	if user.Id != userId {
		return false
	}

	return true
}
