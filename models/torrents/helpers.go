package torrents

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/activities"
	"github.com/NyaaPantsu/nyaa/models/notifications"
	"github.com/NyaaPantsu/nyaa/utils/publicSettings"
)

// ExistOrDelete : Check if a torrent exist with the same hash and if it can be replaced, it is replaced
func ExistOrDelete(hash string, user *models.User) error {
	torrentIndb := models.Torrent{}
	models.ORM.Unscoped().Model(&models.Torrent{}).Where("torrent_hash = ?", hash).First(&torrentIndb)
	if torrentIndb.ID > 0 {
		if user.CurrentUserIdentical(torrentIndb.UploaderID) && torrentIndb.IsDeleted() && !torrentIndb.IsBlocked() { // if torrent is not locked and is deleted and the user is the actual owner
			torrent, _, err := DefinitelyDelete(torrentIndb.ID)
			if err != nil {
				return err
			}
			activities.Log(&models.User{}, torrent.Identifier(), "delete", "torrent_deleted_by", strconv.Itoa(int(torrent.ID)), torrent.Uploader.Username, user.Username)
		} else {
			return errors.New("Torrent already in database")
		}
	}
	return nil
}

// NewTorrentEvent : Should be called when you create a new torrent
func NewTorrentEvent(user *models.User, torrent *models.Torrent) error {
	url := "/view/" + strconv.FormatUint(uint64(torrent.ID), 10)
	if user.ID > 0 && config.Conf.Users.DefaultUserSettings["new_torrent"] { // If we are a member and notifications for new torrents are enabled
		user.GetFollowers()          // We populate the liked field for users
		if len(user.Followers) > 0 { // If we are followed by at least someone
			for _, follower := range user.Followers {
				follower.ParseSettings() // We need to call it before checking settings
				if follower.Settings.Get("new_torrent") {
					T, _, _ := publicSettings.TfuncAndLanguageWithFallback(follower.Language, follower.Language) // We need to send the notification to every user in their language
					notifications.NotifyUser(&follower, torrent.Identifier(), fmt.Sprintf(T("new_torrent_uploaded"), torrent.Name, user.Username), url, follower.Settings.Get("new_torrent_email"))
				}
			}
		}
	}
	return nil
}
