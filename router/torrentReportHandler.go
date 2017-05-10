package router

import (
	"net/http"

	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/service/moderation"
	"github.com/ewhal/nyaa/service/torrent"
	"github.com/ewhal/nyaa/util/modelHelper"
)

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

func GetTorrentReportHandler(w http.ResponseWriter, r *http.Request) {
	torrentReports, err := moderationService.GetTorrentReports()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = torrentReportTemplate.ExecuteTemplate(w, "torrent_report.html", model.TorrentReportsToJSON(torrentReports))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func DeleteTorrentHandler(w http.ResponseWriter, r *http.Request) {
	// TODO Figure out how to get torrent report id from form
	var err error
	var id string
	_, err = torrentService.DeleteTorrent(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
