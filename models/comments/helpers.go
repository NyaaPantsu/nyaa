package comments

import (
	"fmt"
	"strconv"

	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/notifications"
	"github.com/NyaaPantsu/nyaa/utils/publicSettings"
)

// When a new comment is added this is called
func NewCommentEvent(comment *models.Comment, torrent *models.Torrent) {
	comment.Torrent = torrent
	if comment.UserID == torrent.UploaderID {
		return
	}
	url := "/view/" + strconv.FormatUint(uint64(torrent.ID), 10)
	if torrent.UploaderID > 0 {
		torrent.Uploader.ParseSettings()
		if torrent.Uploader.Settings.Get("new_comment") {
			T, _, _ := publicSettings.TfuncAndLanguageWithFallback(torrent.Uploader.Language, torrent.Uploader.Language) // We need to send the notification to every user in their language
			notifications.NotifyUser(torrent.Uploader, comment.Identifier(), fmt.Sprintf(T("new_comment_on_torrent"), torrent.Name), url, torrent.Uploader.Settings.Get("new_comment_email"))
		}
	}
}
