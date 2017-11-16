package torrentController

import (
	"net/http"
	"strconv"

	"github.com/NyaaPantsu/nyaa/controllers/router"
	"github.com/NyaaPantsu/nyaa/models/notifications"
	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/NyaaPantsu/nyaa/templates"
	"github.com/NyaaPantsu/nyaa/utils/captcha"
	"github.com/NyaaPantsu/nyaa/utils/filelist"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/gin-gonic/gin"
)

// ViewHandler : Controller for displaying a torrent
func ViewHandler(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)
	messages := msg.GetMessages(c)
	user := router.GetUser(c)

	// Retrieve the torrent
	torrent, err := torrents.FindByID(uint(id))

	// Toggle the notifications related to this torrent as read
	if user.ID > 0 {
		notifications.ToggleReadNotification(torrent.Identifier(), user)
	}

	// If torrent not found, display 404
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	// We load tags for user and torrents
	user.LoadTags(torrent)
	torrent.LoadTags()

	// Convert torrent to the JSON Model used to display a torrent
	// Since many datas need to be parsed from a simple torrent model to the actual display
	b := torrent.ToJSON()
	// Get the folder root for the filelist view
	folder := filelist.FileListToFolder(torrent.FileList, "root")
	captchaID := ""
	//Generate a captcha
	if user.NeedsCaptcha() {
		captchaID = captcha.GetID()
	}
	
	// Display success message on upload
	if c.Request.URL.Query()["success"] != nil {
		if torrent.IsBlocked() {
			messages.AddInfoT("infos", "torrent_uploaded_locked")
		} else {
			messages.AddInfoT("infos", "torrent_uploaded")
		}
	}
	// Display success message on edit
	if c.Request.URL.Query()["success_edit"] != nil {
		messages.AddInfoT("infos", "torrent_updated")
	}
	// Display wrong captcha error message
	if c.Request.URL.Query()["badcaptcha"] != nil {
		messages.AddErrorT("errors", "bad_captcha")
	}
	// Display reported successful message
	if c.Request.URL.Query()["reported"] != nil {
		messages.AddInfoTf("infos", "report_msg", id)
	}

	if c.Request.URL.Query()["followed"] != nil {
		messages.AddInfoTf("infos", "user_followed_msg", b.UploaderName)
	}
	if c.Request.URL.Query()["unfollowed"] != nil {
		messages.AddInfoTf("infos", "user_unfollowed_msg", b.UploaderName)
	}
	
	// Display finally the view
	templates.Torrent(c, b, folder, captchaID)
}

// ViewHeadHandler : Controller for checking a torrent
func ViewHeadHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		return
	}

	_, err = torrents.FindRawByID(uint(id))

	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.AbortWithStatus(http.StatusOK)
}
