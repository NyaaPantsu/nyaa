package moderationService

import (
	"errors"
	"net/http"

	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
)

// Return torrentReport in case we did modified it (ie: CreatedAt field)
func CreateTorrentReport(torrentReport model.TorrentReport) (model.TorrentReport, error) {
	if db.ORM.Create(&torrentReport).Error != nil {
		return torrentReport, errors.New("TorrentReport was not created")
	}
	return torrentReport, nil
}

func DeleteTorrentReport(id int) (int, error) {
	var torrentReport model.TorrentReport
	if db.ORM.First(&torrentReport, id).RecordNotFound() {
		return http.StatusNotFound, errors.New("Trying to delete a torrent report that does not exists.")
	}
	if db.ORM.Delete(&torrentReport).Error != nil {
		return http.StatusInternalServerError, errors.New("User is not deleted.")
	}
	return http.StatusOK, nil
}

// TODO Add WhereParams to filter the torrent reports (ie: searching description)
// TODO Use limit, offset
func GetTorrentReports() ([]model.TorrentReport, error) {
	var torrentReports []model.TorrentReport
	if db.ORM.Preload("User").Preload("Torrent").Find(&torrentReports).Error != nil {
		return nil, errors.New("Problem finding all torrent reports.")
	}
	return torrentReports, nil
}
