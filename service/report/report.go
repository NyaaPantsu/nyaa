package reportService

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/NyaaPantsu/nyaa/db"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/service"
)

// CreateTorrentReport : Return torrentReport in case we did modified it (ie: CreatedAt field)
func CreateTorrentReport(torrentReport model.TorrentReport) error {
	if db.ORM.Create(&torrentReport).Error != nil {
		return errors.New("TorrentReport was not created")
	}
	return nil
}

// DeleteTorrentReport : Delete a torrent report by id
func DeleteTorrentReport(id uint) (error, int) {
	var torrentReport model.TorrentReport
	if db.ORM.First(&torrentReport, id).RecordNotFound() {
		return errors.New("Trying to delete a torrent report that does not exists"), http.StatusNotFound
	}
	if err := db.ORM.Delete(&torrentReport).Error; err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}

// DeleteDefinitelyTorrentReport : Delete definitely a torrent report by id
func DeleteDefinitelyTorrentReport(id uint) (error, int) {
	var torrentReport model.TorrentReport
	if db.ORM.Unscoped().First(&torrentReport, id).RecordNotFound() {
		return errors.New("Trying to delete a torrent report that does not exists"), http.StatusNotFound
	}
	if err := db.ORM.Unscoped().Delete(&torrentReport).Error; err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}

func getTorrentReportsOrderBy(parameters *serviceBase.WhereParams, orderBy string, limit int, offset int, countAll bool) (
	torrentReports []model.TorrentReport, count int, err error,
) {
	var conditionArray []string
	var params []interface{}
	if parameters != nil { // if there is where parameters
		if len(parameters.Conditions) > 0 {
			conditionArray = append(conditionArray, parameters.Conditions)
		}
		params = parameters.Params
	}
	conditions := strings.Join(conditionArray, " AND ")
	if countAll {
		err = db.ORM.Model(&torrentReports).Where(conditions, params...).Count(&count).Error
		if err != nil {
			return
		}
	}

	// build custom db query for performance reasons
	dbQuery := "SELECT * FROM torrent_reports"
	if conditions != "" {
		dbQuery = dbQuery + " WHERE " + conditions
	}

	if orderBy == "" { // default OrderBy
		orderBy = "torrent_report_id DESC"
	}
	dbQuery = dbQuery + " ORDER BY " + orderBy
	if limit != 0 || offset != 0 { // if limits provided
		dbQuery = dbQuery + " LIMIT " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(offset)
	}
	err = db.ORM.Preload("Torrent").Preload("User").Raw(dbQuery, params...).Find(&torrentReports).Error
	return
}

// GetTorrentReportsOrderBy : Get torrents based on search parameters with order
func GetTorrentReportsOrderBy(parameters *serviceBase.WhereParams, orderBy string, limit int, offset int) ([]model.TorrentReport, int, error) {
	return getTorrentReportsOrderBy(parameters, orderBy, limit, offset, true)
}

// GetAllTorrentReports : Get all torrents
func GetAllTorrentReports(limit int, offset int) ([]model.TorrentReport, int, error) {
	return GetTorrentReportsOrderBy(nil, "", limit, offset)
}
