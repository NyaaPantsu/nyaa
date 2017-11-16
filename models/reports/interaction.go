package reports

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	elastic "gopkg.in/olivere/elastic.v5"

	"github.com/NyaaPantsu/nyaa/models"
)

// Query : Interface to pass for torrents query
type Query interface {
	String() string
	ToDBQuery() (string, []interface{})
	ToESQuery(client *elastic.Client) (*elastic.SearchService, error)
	Append(string, ...interface{})
	Prepend(string, ...interface{})
}

func Delete(id uint) (*models.TorrentReport, int, error) {
	var torrentReport models.TorrentReport
	db := models.ORM.Unscoped()
	if db.First(&torrentReport, id).RecordNotFound() {
		return &torrentReport, http.StatusNotFound, errors.New("try_to_delete_report_inexistant")
	}
	if _, err := torrentReport.Delete(); err != nil {
		return &torrentReport, http.StatusInternalServerError, err
	}
	return &torrentReport, http.StatusOK, nil
}

func DeleteAll() {
	models.ORM.Delete(&models.TorrentReport{})
}

func findOrderBy(parameters Query, orderBy string, limit int, offset int, countAll bool) (
	torrentReports []models.TorrentReport, count int, err error,
) {
	var conditionArray []string
	var params []interface{}
	if parameters != nil { // if there is where parameters
		condition, wheres := parameters.ToDBQuery()
		if len(condition) > 0 {
			conditionArray = append(conditionArray, condition)
			params = wheres
		}
	}
	conditions := strings.Join(conditionArray, " AND ")
	if countAll {
		err = models.ORM.Model(&torrentReports).Where(conditions, params...).Count(&count).Error
		if err != nil {
			return
		}
	}
	
	var blankReport models.TorrentReport

	// build custom db query for performance reasons
	dbQuery := "SELECT * FROM " + blankReport.TableName()
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
	err = models.ORM.Preload("Torrent").Preload("User").Raw(dbQuery, params...).Find(&torrentReports).Error
	return
}

// GetOrderBy : Get torrents reports based on search parameters with order
func FindOrderBy(parameters Query, orderBy string, limit int, offset int) ([]models.TorrentReport, int, error) {
	return findOrderBy(parameters, orderBy, limit, offset, true)
}

// GetAll : Get all torrents report
func GetAll(limit int, offset int) ([]models.TorrentReport, int, error) {
	return FindOrderBy(nil, "", limit, offset)
}
