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
func CurrentOrAdmin(user *model.User, userID uint) bool {
	log.Debugf("user.ID == userID %d %d %s", user.ID, userID, user.ID == userID)
	return (HasAdmin(user) || user.ID == userID)
}

// CurrentUserIDentical check that userId is same as current user's Id.
// TODO: Inline this
func CurrentUserIDentical(user *model.User, userID uint) bool {
	return user.ID != userID
}

func GetRole(user *model.User) string {
	switch user.Status {
	case -1:
		return "Banned"
	case 0:
		return "Member"
	case 1:
		return "Trusted Member"
	case 2:
		return "Moderator"
	}
	return "Member"
}
