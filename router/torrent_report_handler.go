package router

/*import (
	"net/http"
	"strconv"

	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/service/moderation"
	"github.com/NyaaPantsu/nyaa/service/user/permission"
	"github.com/gorilla/mux"
)*/

/*
func SanitizeTorrentReport(torrentReport *model.TorrentReport) {
	// TODO unescape html ?
	return
}

func IsValidTorrentReport() bool {
	// TODO Validate, see if user_id already reported, see if torrent exists
	return true
}

// TODO Only allow moderators for each action in this file
func CreateTorrentReportHandler(c *gin.Context) {
	var torrentReport model.TorrentReport
	var err error

	modelHelper.BindValueForm(&torrentReport, r)
	if IsValidTorrentReport() {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	SanitizeTorrentReport(&torrentReport)
	_, err = moderationService.CreateTorrentReport(torrentReport)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func DeleteTorrentReportHandler(c *gin.Context) {
	// TODO Figure out how to get torrent report id from form
	var id int
	var err error
	_, err = moderationService.DeleteTorrentReport(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
*/

/*

func DeleteTorrentHandler(c *gin.Context) {
	// TODO Figure out how to get torrent report id from form
	var err error
	var id string
	_, err = torrentService.DeleteTorrent(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}*/

/*func GetTorrentReportHandler(c *gin.Context) {
    currentUser := getUser(c)
    if userPermission.HasAdmin(currentUser) {
		vars := mux.Vars(r)
		page, _ := strconv.Atoi(c.Query("page"))
		offset := 100
		userid := c.Query("userid")
		var conditions string
		var values []interface{}
		if (userid != "") {
			conditions = "user_id = ?"
			values = append(values, userid)
		}

	torrentReports, nbReports, err := moderationService.GetTorrentReports(offset, page * offset, conditions, values...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = torrentReportTemplate.ExecuteTemplate(w, "admin_index.html", ViewTorrentReportsVariables{model.TorrentReportsToJSON(torrentReports), newSearchForm(), navigation{nbReports, offset, page, "mod_trlist_page"}, currentUser, r.URL})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
    } else {
            http.Error(w, "admins only", http.StatusForbidden)
    }
}*/
