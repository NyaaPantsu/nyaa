package moderatorController

import (
	"html"
	"net/http"
	"strconv"

	"github.com/NyaaPantsu/nyaa/templates"
	"github.com/NyaaPantsu/nyaa/utils/log"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/search"
	"github.com/gin-gonic/gin"
)

// DeletedTorrentsModPanel : Controller for viewing deleted torrents, accept common search arguments
func DeletedTorrentsModPanel(c *gin.Context) {
	page := c.Param("page")
	messages := msg.GetMessages(c) // new utils for errors and infos
	deleted := c.Request.URL.Query()["deleted"]
	unblocked := c.Request.URL.Query()["unblocked"]
	blocked := c.Request.URL.Query()["blocked"]
	if deleted != nil {
		messages.AddInfoT("infos", "torrent_deleted_definitely")
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

	searchParam, torrents, count, err := search.ByQueryDeleted(c, pagenum)
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

	nav := templates.Navigation{count, int(searchParam.Max), pagenum, "mod/torrents/deleted/p"}
	search := searchForm
	templates.ModelList(c, "admin/torrentlist.jet.html", torrents, nav, search)
}

// DeletedTorrentsPostPanel : Controller for viewing deleted torrents after a mass update, accept common search arguments
func DeletedTorrentsPostPanel(c *gin.Context) {
	torrentManyAction(c)
	DeletedTorrentsModPanel(c)
}
