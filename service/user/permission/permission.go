package userPermission

import (
	"github.com/NyaaPantsu/nyaa/db"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/util/log"
)

// HasAdmin checks that user has an admin permission.
func HasAdmin(user *model.User) bool {
	return user.IsModerator()
}

// CurrentOrAdmin check that user has admin permission or user is the current user.
func CurrentOrAdmin(user *model.User, userID uint) bool {
	log.Debugf("user.ID == userID %d %d %s", user.ID, userID, user.ID == userID)
	return (HasAdmin(user) || user.ID == userID)
}

// CurrentUserIdentical check that userID is same as current user's ID.
// TODO: Inline this (won't go do this for us?)
func CurrentUserIdentical(user *model.User, userID uint) bool {
	return user.ID == userID
}

func NeedsCaptcha(user *model.User) bool {
	// Trusted members & Moderators don't
	return !(user.IsTrusted() || user.IsModerator())
}

func GetRole(user *model.User) string {
	switch user.Status {
	case model.UserStatusBanned:
		return "Banned"
	case model.UserStatusMember:
		return "Member"
	case model.UserStatusTrusted:
		return "Trusted Member"
	case model.UserStatusModerator:
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
