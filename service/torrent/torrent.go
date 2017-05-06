package torrentService
import (
	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/config"
	"github.com/jinzhu/gorm"
	"errors"
	"strings"
)

type WhereParams struct {
	Conditions string // Ex : name LIKE ? AND category_id LIKE ?
	Params     []interface{}
}

/* Function to interact with Models
 *
 * Get the torrents with where clause
 *
 */

// don't need raw SQL once we get MySQL
func GetFeeds() []model.Feed {
	var result []model.Feed
	rows, err := db.ORM.DB().
		Query(
			"SELECT `torrent_id` AS `id`, `torrent_name` AS `name`, `torrent_hash` AS `hash`, `timestamp` FROM `torrents` " +
				"ORDER BY `timestamp` desc LIMIT 50")
	if err == nil {
		for rows.Next() {
			item := model.Feed{}
			rows.Scan(&item.Id, &item.Name, &item.Hash, &item.Timestamp)
			magnet := "magnet:?xt=urn:btih:" + strings.TrimSpace(item.Hash) + "&dn=" + item.Name + config.Trackers
			item.Magnet = magnet
			// memory hog
			result = append(result, item)
		}
		rows.Close()
	}
	return result
}

func GetTorrentById(id string) (model.Torrents, error) {
	var torrent model.Torrents

	if db.ORM.Where("torrent_id = ?", id).Preload("Sub_Categories").Find(&torrent).RecordNotFound() {
		return torrent, errors.New("Article is not found.")
	}

	return torrent, nil
}

func GetTorrentsOrderBy(parameters *WhereParams, orderBy string, limit int, offset int) ([]model.Torrents, int) {
	var torrents []model.Torrents
	var dbQuery *gorm.DB
	var count int
	conditions := "torrent_hash is not null" //filter out broken entries
	var params []interface{}
	if parameters != nil { // if there is where parameters
		conditions += " AND " + parameters.Conditions
		params = parameters.Params
	}
	db.ORM.Model(&torrents).Where(conditions, params...).Count(&count)
	dbQuery = db.ORM.Model(&torrents).Where(conditions, params...)

	if orderBy == "" {
		orderBy = "torrent_id DESC"
	} // Default OrderBy
	if limit != 0 || offset != 0 { // if limits provided
		dbQuery = dbQuery.Limit(limit).Offset(offset)
	}
	dbQuery.Order(orderBy).Preload("Categories").Preload("Sub_Categories").Find(&torrents)
	return torrents, count

}

/* Functions to simplify the get parameters of the main function
 *
 * Get Torrents with where parameters and limits, order by default
 */
func GetTorrents(parameters WhereParams, limit int, offset int) ([]model.Torrents, int) {
	return GetTorrentsOrderBy(&parameters, "", limit, offset)
}

/* Get Torrents with where parameters but no limit and order by default (get all the torrents corresponding in the db)
 */
func GetTorrentsDB(parameters WhereParams) ([]model.Torrents, int) {
	return GetTorrentsOrderBy(&parameters, "", 0, 0)
}

/* Function to get all torrents
 */

func GetAllTorrentsOrderBy(orderBy string, limit int, offset int) ([]model.Torrents, int) {

	return GetTorrentsOrderBy(nil, orderBy, limit, offset)
}

func GetAllTorrents(limit int, offset int) ([]model.Torrents, int) {
	return GetTorrentsOrderBy(nil, "", limit, offset)
}

func GetAllTorrentsDB() ([]model.Torrents, int) {
	return GetTorrentsOrderBy(nil, "", 0, 0)
}

/* Function to get all categories with/without torrents (bool)
 */
func GetAllCategories(populatedWithTorrents bool) []model.Categories {
	var categories []model.Categories
	if populatedWithTorrents {
		db.ORM.Preload("Torrents").Preload("Sub_Categories").Find(&categories)
	} else {
		db.ORM.Preload("Sub_Categories").Find(&categories)
	}
	return categories
}

func CreateWhereParams(conditions string, params ...string) WhereParams {
	whereParams := WhereParams{}
	whereParams.Conditions = conditions
	for i, _ := range params {
		whereParams.Params = append(whereParams.Params, params[i])
	}

	return whereParams
}
