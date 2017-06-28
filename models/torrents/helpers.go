package torrents

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/util/publicSettings"
)

// ExistOrDelete : Check if a torrent exist with the same hash and if it can be replaced, it is replaced
func ExistOrDelete(hash string, user *models.User) error {
	torrentIndb := models.Torrent{}
	models.ORM.Unscoped().Model(&models.Torrent{}).Where("torrent_hash = ?", hash).First(&torrentIndb)
	if torrentInmodels.ID > 0 {
		if user.CurrentUserIdentical(torrentInmodels.UploaderID) && torrentInmodels.IsDeleted() && !torrentInmodels.IsBlocked() { // if torrent is not locked and is deleted and the user is the actual owner
			torrent, _, err := DefinitelyDeleteTorrent(strconv.Itoa(int(torrentInmodels.ID)))
			if err != nil {
				return err
			}
			activity.Log(&models.User{}, torrent.Identifier(), "delete", "torrent_deleted_by", strconv.Itoa(int(torrent.ID)), torrent.Uploader.Username, user.Username)
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
		userService.GetFollowers(user) // We populate the liked field for users
		if len(user.Followers) > 0 {   // If we are followed by at least someone
			for _, follower := range user.Followers {
				follower.ParseSettings() // We need to call it before checking settings
				if follower.Settings.Get("new_torrent") {
					T, _, _ := publicSettings.TfuncAndLanguageWithFallback(follower.Language, follower.Language) // We need to send the notification to every user in their language
					notifierService.NotifyUser(&follower, torrent.Identifier(), fmt.Sprintf(T("new_torrent_uploaded"), torrent.Name, user.Username), url, follower.Settings.Get("new_torrent_email"))
				}
			}
		}
	}
	return nil
}
