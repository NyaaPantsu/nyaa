package reportService

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/service"
)

// Return torrentReport in case we did modified it (ie: CreatedAt field)
func CreateTorrentReport(torrentReport model.TorrentReport) error {
	if db.ORM.Create(&torrentReport).Error != nil {
		return errors.New("TorrentReport was not created")
	}
	return nil
}

func DeleteTorrentReport(id int) (error, int) {
	var torrentReport model.TorrentReport
	if db.ORM.First(&torrentReport, id).RecordNotFound() {
		return errors.New("Trying to delete a torrent report that does not exists."), http.StatusNotFound
	}
	if err := db.ORM.Delete(&torrentReport).Error; err != nil {
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
	// TODO: Vulnerable to injections. Use query builder. (is it?)

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
	err = db.ORM.Preload("Torrent").Preload("User").Raw(dbQuery, params...).Find(&torrentReports).Error //fixed !!!!
	return
}

func GetTorrentReportsOrderBy(parameters *serviceBase.WhereParams, orderBy string, limit int, offset int) ([]model.TorrentReport, int, error) {
	return getTorrentReportsOrderBy(parameters, orderBy, limit, offset, true)
}

func GetAllTorrentReports(limit int, offset int) ([]model.TorrentReport, int, error) {
	return GetTorrentReportsOrderBy(nil, "", limit, offset)
}
