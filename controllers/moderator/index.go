package moderatorController

import (
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/comments"
	"github.com/NyaaPantsu/nyaa/models/reports"
	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/NyaaPantsu/nyaa/models/users"
	"github.com/NyaaPantsu/nyaa/templates"
	"github.com/gin-gonic/gin"
)

// IndexModPanel : Controller for showing index page of Mod Panel
func IndexModPanel(c *gin.Context) {
	offset := 10
	torrents, _, _ := torrents.FindAllOrderBy("torrent_id DESC", offset, 0)
	users, _ := users.FindUsersForAdmin(offset, 0)
	comments, _ := comments.FindAll(offset, 0, "", "")
	torrentReports, _, _ := reports.GetAll(offset, 0)

	templates.PanelAdmin(c, torrents, models.TorrentReportsToJSON(torrentReports), users, comments)
}
