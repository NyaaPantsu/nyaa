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

	if c.Request.URL.Query()["success"] != nil {
		messages.AddInfoT("infos", "torrent_uploaded")
	}
	if c.Request.URL.Query()["success_edit"] != nil {
		messages.AddInfoT("infos", "torrent_updated")
	}
	if c.Request.URL.Query()["badcaptcha"] != nil {
		messages.AddErrorT("errors", "bad_captcha")
	}
	if c.Request.URL.Query()["reported"] != nil {
		messages.AddInfoTf("infos", "report_msg", id)
	}

	torrent, err := torrents.FindByID(uint(id))

	if c.Request.URL.Query()["notif"] != nil {
		notifications.ToggleReadNotification(torrent.Identifier(), user.ID)
	}

	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	b := torrent.ToJSON()
	folder := filelist.FileListToFolder(torrent.FileList, "root")
	captchaID := ""
	if user.NeedsCaptcha() {
		captchaID = captcha.GetID()
	}
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
