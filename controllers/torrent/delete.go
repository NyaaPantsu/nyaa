package torrentController

import (
	"net/http"
	"strconv"

	"github.com/NyaaPantsu/nyaa/controllers/router"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/activities"
	"github.com/NyaaPantsu/nyaa/models/reports"
	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/NyaaPantsu/nyaa/utils/search/structs"
	"github.com/gin-gonic/gin"
)

// TorrentDeleteUserPanel : Controller for deleting a user torrent by a user
func TorrentDeleteUserPanel(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 32)
	currentUser := router.GetUser(c)
	torrent, _ := torrents.FindByID(uint(id))
	if currentUser.CurrentOrAdmin(torrent.UploaderID) {
		_, _, err := torrent.Delete(false)
		if err == nil {
			if torrent.Uploader == nil {
				torrent.Uploader = &models.User{}
			}
			_, username := torrents.HideUser(torrent.UploaderID, torrent.Uploader.Username, torrent.Hidden)
			if currentUser.HasAdmin() { // We hide username on log activity if user is not admin and torrent is hidden
				activities.Log(&models.User{}, torrent.Identifier(), "delete", "torrent_deleted_by", strconv.Itoa(int(torrent.ID)), username, currentUser.Username)
			} else {
				activities.Log(&models.User{}, torrent.Identifier(), "delete", "torrent_deleted_by", strconv.Itoa(int(torrent.ID)), username, username)
			}
			//delete reports of torrent
			whereParams := structs.CreateWhereParams("torrent_id = ?", id)
			torrentReports, _, _ := reports.FindOrderBy(&whereParams, "", 0, 0)
			for _, report := range torrentReports {
				report.Delete(false)
			}
		}
		c.Redirect(http.StatusSeeOther, "/?deleted")
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}
