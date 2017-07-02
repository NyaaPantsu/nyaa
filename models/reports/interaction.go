package reports

import (
	"errors"
	"net/http"
	"nyaa-master/db"
	"strconv"
	"strings"

	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/utils/search/structs"
)

// Create : Return torrentReport in case we did modified it (ie: CreatedAt field)
func Create(torrentReport models.TorrentReport) error {
	if models.ORM.Create(&torrentReport).Error != nil {
		return errors.New("torrent_report_not_created")
	}
	return nil
}

// Delete : Delete a torrent report by id
func Delete(id uint) (*models.TorrentReport, int, error) {
	return delete(id, false)
}

// DeleteDefinitely : Delete definitely a torrent report by id
func DeleteDefinitely(id uint) (*models.TorrentReport, int, error) {
	return delete(id, true)
}

func delete(id uint, definitely bool) (*models.TorrentReport, int, error) {
	var torrentReport models.TorrentReport
	db := models.ORM
	if definitely {
		db = models.ORM.Unscoped()
	}
	if db.First(&torrentReport, id).RecordNotFound() {
		return &torrentReport, http.StatusNotFound, errors.New("try_to_delete_report_inexistant")
	}
	if _, err := torrentReport.Delete(false); err != nil {
		return &torrentReport, http.StatusInternalServerError, err
	}
	return &torrentReport, http.StatusOK, nil
}

func findOrderBy(parameters *structs.WhereParams, orderBy string, limit int, offset int, countAll bool) (
	torrentReports []models.TorrentReport, count int, err error,
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

// GetOrderBy : Get torrents reports based on search parameters with order
func FindOrderBy(parameters *structs.WhereParams, orderBy string, limit int, offset int) ([]models.TorrentReport, int, error) {
	return findOrderBy(parameters, orderBy, limit, offset, true)
}

// GetAll : Get all torrents report
func GetAll(limit int, offset int) ([]models.TorrentReport, int, error) {
	return FindOrderBy(nil, "", limit, offset)
}
