package moderatorController

import (
	"html"
	"net/http"
	"strconv"

	"github.com/NyaaPantsu/nyaa/controllers/router"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/activities"
	"github.com/NyaaPantsu/nyaa/models/reports"
	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/NyaaPantsu/nyaa/templates"
	"github.com/NyaaPantsu/nyaa/utils/log"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/search"
	"github.com/NyaaPantsu/nyaa/utils/search/structs"
	"github.com/NyaaPantsu/nyaa/utils/upload"
	"github.com/NyaaPantsu/nyaa/utils/validator/torrent"
	"github.com/gin-gonic/gin"
)

// TorrentsListPanel : Controller for listing torrents, can accept common search arguments
func TorrentsListPanel(c *gin.Context) {
	page := c.Param("page")
	messages := msg.GetMessages(c)
	deleted := c.Request.URL.Query()["deleted"]
	unblocked := c.Request.URL.Query()["unblocked"]
	blocked := c.Request.URL.Query()["blocked"]
	if deleted != nil {
		messages.AddInfoTf("infos", "torrent_deleted", "")
	}
	if blocked != nil {
		messages.AddInfoT("infos", "torrent_blocked")
	}
	if unblocked != nil {
		messages.AddInfoT("infos", "torrent_unblocked")
	}

	var err error
	pagenum := 1
	if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if !log.CheckError(err) {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	searchParam, torrents, count, err := search.ByQueryWithUser(c, pagenum)
	if err != nil {
		messages.Error(err)
	}
	category := ""
	if len(searchParam.Category) > 0 {
		category = searchParam.Category[0].String()
	}
	searchForm := templates.SearchForm{
		TorrentParam:     searchParam,
		Category:         category,
		ShowItemsPerPage: true,
	}

	nav := templates.Navigation{count, int(searchParam.Max), pagenum, "mod/torrents/p"}

	templates.ModelList(c, "admin/torrentlist.jet.html", torrents, nav, searchForm)
}

// TorrentEditModPanel : Controller for editing a torrent after GET request
func TorrentEditModPanel(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 32)
	torrent, _ := torrents.FindUnscopeByID(uint(id))

	torrentJSON := torrent.ToJSON()
	uploadForm := upload.NewTorrentRequest()
	uploadForm.Name = torrentJSON.Name
	uploadForm.Category = torrentJSON.Category + "_" + torrentJSON.SubCategory
	uploadForm.Status = torrentJSON.Status
	uploadForm.Hidden = torrent.Hidden
	uploadForm.WebsiteLink = string(torrentJSON.WebsiteLink)
	uploadForm.Description = string(torrentJSON.Description)
	uploadForm.Languages = torrent.Languages

	templates.Form(c, "admin/paneltorrentedit.jet.html", uploadForm)
}

// TorrentPostEditModPanel : Controller for editing a torrent after POST request
func TorrentPostEditModPanel(c *gin.Context) {
	var uploadForm torrentValidator.UpdateRequest
	id, _ := strconv.ParseInt(c.Query("id"), 10, 32)
	messages := msg.GetMessages(c)
	torrent, _ := torrents.FindUnscopeByID(uint(id))
	currentUser := router.GetUser(c)
	if torrent.ID > 0 {
		errUp := upload.ExtractEditInfo(c, &uploadForm.Update)
		uploadForm.ID = uint(id)
		if errUp != nil {
			messages.AddErrorT("errors", "fail_torrent_update")
		}
		if !messages.HasErrors() {
			_, err := upload.UpdateUnscopeTorrent(&uploadForm, &torrent, currentUser).UpdateUnscope()
			messages.AddInfoT("infos", "torrent_updated")
			if err == nil { // We only log edit torrent for admins
				if torrent.Uploader == nil {
					torrent.Uploader = &models.User{}
				}
				_, username := torrents.HideUser(torrent.UploaderID, torrent.Uploader.Username, torrent.Hidden)
				activities.Log(&models.User{}, torrent.Identifier(), "edit", "torrent_edited_by", strconv.Itoa(int(torrent.ID)), username, currentUser.Username)
			}
		}
	}
	templates.Form(c, "admin/paneltorrentedit.jet.html", uploadForm.Update)
}

// TorrentDeleteModPanel : Controller for deleting a torrent
func TorrentDeleteModPanel(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 32)
	definitely := c.Request.URL.Query()["definitely"]

	var returnRoute = "/mod/torrents"
	torrent, errFind := torrents.FindByID(uint(id))
	if errFind == nil {
		var err error
		if definitely != nil {
			_, _, err = torrent.DefinitelyDelete()

			//delete reports of torrent
			whereParams := structs.CreateWhereParams("torrent_id = ?", id)
			reports, _, _ := reports.FindOrderBy(&whereParams, "", 0, 0)
			for _, report := range reports {
				report.Delete(true)
			}
			returnRoute = "/mod/torrents/deleted"
		} else {
			_, _, err = torrent.Delete(false)

			//delete reports of torrent
			whereParams := structs.CreateWhereParams("torrent_id = ?", id)
			reports, _, _ := reports.FindOrderBy(&whereParams, "", 0, 0)
			for _, report := range reports {
				report.Delete(false)
			}
		}
		if err == nil {
			if torrent.Uploader == nil {
				torrent.Uploader = &models.User{}
			}
			_, username := torrents.HideUser(torrent.UploaderID, torrent.Uploader.Username, torrent.Hidden)
			activities.Log(&models.User{}, torrent.Identifier(), "delete", "torrent_deleted_by", strconv.Itoa(int(torrent.ID)), username, router.GetUser(c).Username)
		}
	}

	c.Redirect(http.StatusSeeOther, returnRoute+"?deleted")
}

// TorrentBlockModPanel : Controller to lock torrents, redirecting to previous page
func TorrentBlockModPanel(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 32)
	torrent, _, err := torrents.ToggleBlock(uint(id))
	var returnRoute, action string
	if torrent.IsDeleted() {
		returnRoute = "/mod/torrents/deleted"
	} else {
		returnRoute = "/mod/torrents"
	}
	if torrent.IsBlocked() {
		action = "blocked"
	} else {
		action = "unblocked"
	}
	if err == nil {
		if torrent.Uploader == nil {
			torrent.Uploader = &models.User{}
		}
		_, username := torrents.HideUser(torrent.UploaderID, torrent.Uploader.Username, torrent.Hidden)
		activities.Log(&models.User{}, torrent.Identifier(), action, "torrent_"+action+"_by", strconv.Itoa(int(torrent.ID)), username, router.GetUser(c).Username)
	}

	c.Redirect(http.StatusSeeOther, returnRoute+"?"+action)
}

// TorrentsPostListPanel : Controller for listing torrents, after POST request when mass update
func TorrentsPostListPanel(c *gin.Context) {
	torrentManyAction(c)
	TorrentsListPanel(c)
}
