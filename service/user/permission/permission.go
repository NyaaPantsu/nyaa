package userPermission

import (
	"github.com/ewhal/nyaa/db"
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

// CurrentUserIdentical check that userID is same as current user's ID.
// TODO: Inline this
func CurrentUserIdentical(user *model.User, userID uint) bool {
	return user.ID == userID
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

func IsFollower(user *model.User, currentUser *model.User) bool {
	var likingUserCount int
	db.ORM.Model(&model.UserFollows{}).Where("user_id = ? and following = ?", user.ID, currentUser.ID).Count(&likingUserCount)
	if likingUserCount != 0 {
		return true
	}
	return false
}
