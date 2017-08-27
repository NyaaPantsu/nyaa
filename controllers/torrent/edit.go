package torrentController

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/NyaaPantsu/nyaa/controllers/router"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/NyaaPantsu/nyaa/templates"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/upload"
	"github.com/NyaaPantsu/nyaa/utils/validator/torrent"
	"github.com/gin-gonic/gin"
)

// TorrentEditUserPanel : Controller for editing a user torrent by a user, after GET request
func TorrentEditUserPanel(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 32)
	torrent, _ := torrents.FindByID(uint(id))

	currentUser := router.GetUser(c)
	if currentUser.CurrentOrAdmin(torrent.UploaderID) && torrent.ID > 0 {
		uploadForm := torrentValidator.TorrentRequest{}
		uploadForm.Name = torrent.Name
		uploadForm.Category = strconv.Itoa(torrent.Category) + "_" + strconv.Itoa(torrent.SubCategory)
		uploadForm.Remake = torrent.Status == models.TorrentStatusRemake
		uploadForm.WebsiteLink = string(torrent.WebsiteLink)
		uploadForm.Description = string(torrent.Description)
		uploadForm.Hidden = torrent.Hidden
		uploadForm.Languages = torrent.Languages
		uploadForm.Tags.Bind(torrent)
		templates.Form(c, "site/torrents/edit.jet.html", uploadForm)
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

// TorrentPostEditUserPanel : Controller for editing a user torrent by a user, after post request
func TorrentPostEditUserPanel(c *gin.Context) {
	var uploadForm torrentValidator.UpdateRequest
	id, _ := strconv.ParseInt(c.Query("id"), 10, 32)
	uploadForm.ID = uint(id)
	messages := msg.GetMessages(c)
	torrent, _ := torrents.FindByID(uint(id))
	torrent.LoadTags()
	currentUser := router.GetUser(c)
	if torrent.ID > 0 && currentUser.CurrentOrAdmin(torrent.UploaderID) {
		errUp := upload.ExtractEditInfo(c, &uploadForm.Update)
		if errUp != nil {
			messages.AddErrorT("errors", "fail_torrent_update")
		}
		if !messages.HasErrors() {
			upload.UpdateTorrent(&uploadForm, torrent, currentUser).Update(currentUser.HasAdmin())
			c.Redirect(http.StatusSeeOther, fmt.Sprintf("/view/%d?success_edit", id))
			return
		}
		templates.Form(c, "site/torrents/edit.jet.html", uploadForm.Update)
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}
