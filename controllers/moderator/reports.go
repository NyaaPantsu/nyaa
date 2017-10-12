package moderatorController

import (
	"fmt"
	"html"
	"net/http"
	"strconv"

	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/reports"
	"github.com/NyaaPantsu/nyaa/templates"
	"github.com/NyaaPantsu/nyaa/utils/log"
	"github.com/gin-gonic/gin"
)

// TorrentReportListPanel : Controller for listing torrent reports, can accept pages
func TorrentReportListPanel(c *gin.Context) {
	page := c.Param("page")
	pagenum := 1
	offset := 100
	var err error

	if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if !log.CheckError(err) {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	torrentReports, nbReports, _ := reports.GetAll(offset, (pagenum-1)*offset)

	reportJSON := models.TorrentReportsToJSON(torrentReports)
	nav := templates.Navigation{nbReports, offset, pagenum, "mod/reports/p"}
	templates.ModelList(c, "admin/torrent_report.jet.html", reportJSON, nav, templates.NewSearchForm(c))
}

// TorrentReportDeleteModPanel : Controller for deleting a torrent report
func TorrentReportDeleteModPanel(c *gin.Context) {
	id := c.PostForm("id")

	fmt.Println(id)
	idNum, _ := strconv.ParseUint(id, 10, 64)
	_, _, _ = reports.Delete(uint(idNum))
	/* If we need to log report delete activity
	if err == nil {
		activity.Log(&models.User{}, torrent.Identifier(), "delete", "torrent_report_deleted_by", strconv.Itoa(int(report.ID)), router.GetUser(c).Username)
	}
	*/
	c.Redirect(http.StatusSeeOther, "/mod/reports?deleted")
}
