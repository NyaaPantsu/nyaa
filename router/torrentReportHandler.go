package router

/*import (
	"net/http"
	"strconv"

	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/service/moderation"
	"github.com/ewhal/nyaa/service/user/permission"
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
func CreateTorrentReportHandler(w http.ResponseWriter, r *http.Request) {
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

func DeleteTorrentReportHandler(w http.ResponseWriter, r *http.Request) {
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

func DeleteTorrentHandler(w http.ResponseWriter, r *http.Request) {
	// TODO Figure out how to get torrent report id from form
	var err error
	var id string
	_, err = torrentService.DeleteTorrent(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}*/

/*func GetTorrentReportHandler(w http.ResponseWriter, r *http.Request) {
    currentUser := GetUser(r)
    if userPermission.HasAdmin(currentUser) {
		vars := mux.Vars(r)
		page, _ := strconv.Atoi(vars["page"])
		offset := 100
		userid := r.URL.Query().Get("userid")
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
	err = torrentReportTemplate.ExecuteTemplate(w, "admin_index.html", ViewTorrentReportsVariables{model.TorrentReportsToJSON(torrentReports), NewSearchForm(), Navigation{nbReports, offset, page, "mod_trlist_page"}, currentUser, r.URL})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
    } else {
            http.Error(w, "admins only", http.StatusForbidden)
    }
}*/
