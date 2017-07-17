package moderatorController

import (
	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/NyaaPantsu/nyaa/models/users"
	"github.com/NyaaPantsu/nyaa/templates"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/validator/torrent"
	"github.com/gin-gonic/gin"
)

// TorrentReassignModPanel : Controller for reassigning a torrent, after GET request
func TorrentReassignModPanel(c *gin.Context) {
	templates.Form(c, "admin/reassign.jet.html", torrentValidator.ReassignForm{})
}

// ExecuteAction : Function for applying the changes from ReassignForm
func ExecuteReassign(f *torrentValidator.ReassignForm) (int, error) {
	var toBeChanged []uint
	var err error
	if f.By == "olduser" {
		toBeChanged, err = users.FindOldUploadsByUsername(f.Data)
		if err != nil {
			return 0, err
		}
	} else if f.By == "torrentid" {
		toBeChanged = f.Torrents
	}

	num := 0
	for _, torrentID := range toBeChanged {
		torrent, err2 := torrents.FindRawByID(torrentID)
		if err2 == nil {
			torrent.UploaderID = f.AssignTo
			torrent.Update(true)
			num++
		}
	}
	return num, nil
}

// TorrentPostReassignModPanel : Controller for reassigning a torrent, after POST request
func TorrentPostReassignModPanel(c *gin.Context) {
	var rForm torrentValidator.ReassignForm
	messages := msg.GetMessages(c)

	if rForm.ExtractInfo(c) {
		count, err2 := ExecuteReassign(&rForm)
		if err2 != nil {
			messages.AddErrorT("errors", "something_went_wrong")
		} else {
			messages.AddInfoTf("infos", "nb_torrents_updated", count)
		}
	}
	templates.Form(c, "admin/reassign.jet.html", rForm)
}
