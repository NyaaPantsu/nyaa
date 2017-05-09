package userPermission

import (
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/util/log"
)

// HasAdmin checks that user has an admin permission.
func HasAdmin(user *model.User) bool {
	return user.Status == 2
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

func GetRole(user *model.User) string {
	switch user.Status {
	case -1 :
		return "Banned"
	case 0 :
		return "Member"
	case 1 :
		return "Trusted Member"
	case 2 :
		return "Moderator"
	}
	return "Member"
}
