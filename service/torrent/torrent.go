package torrentService

import (
	"errors"
	"strconv"
	"strings"

	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/util"
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
	result := make([]model.Feed, 0, 50)
	rows, err := db.ORM.DB().
		Query(
			"SELECT `torrent_id` AS `id`, `torrent_name` AS `name`, `torrent_hash` AS `hash`, `timestamp` FROM `torrents` " +
				"ORDER BY `timestamp` desc LIMIT 50")
	if err == nil {
		for rows.Next() {
			item := model.Feed{}
			rows.Scan(&item.Id, &item.Name, &item.Hash, &item.Timestamp)
			magnet := util.InfoHashToMagnet(strings.TrimSpace(item.Hash), item.Name, config.Trackers...)
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
	// Postgres DB integer size is 32-bit.
	id_int, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return torrent, err
	}

	tmp := db.ORM.Where("torrent_id = ?", id).Preload("Comments")
	if id_int <= config.LastOldTorrentId {
		// only preload old comments if they could actually exist
		tmp = tmp.Preload("OldComments")
	}
	if tmp.Find(&torrent).RecordNotFound() {
		return torrent, errors.New("Article is not found.")
	}
	// .Preload("Comments.User") doesn't work
	for i := range torrent.Comments {
		torrent.Comments[i].User = new(model.User)
		db.ORM.Where("user_id = ?", torrent.Comments[i].UserId).Find(torrent.Comments[i].User)
	}

	return torrent, nil
}

func GetTorrentsOrderByNoCount(parameters *WhereParams, orderBy string, limit int, offset int) (torrents []model.Torrents, err error) {
	torrents, _, err = getTorrentsOrderBy(parameters, orderBy, limit, offset, false)
	return
}

func GetTorrentsOrderBy(parameters *WhereParams, orderBy string, limit int, offset int) (torrents []model.Torrents, count int, err error) {
	torrents, count, err = getTorrentsOrderBy(parameters, orderBy, limit, offset, true)
	return
}

func getTorrentsOrderBy(parameters *WhereParams, orderBy string, limit int, offset int, countAll bool) (torrents []model.Torrents, count int, err error) {
	var conditionArray []string
	if strings.HasPrefix(orderBy, "filesize") {
		// torrents w/ NULL filesize fuck up the sorting on postgres
		conditionArray = append(conditionArray, "filesize IS NOT NULL")
	}
	var params []interface{}
	if parameters != nil { // if there is where parameters
		if len(parameters.Conditions) > 0 {
			conditionArray = append(conditionArray, parameters.Conditions)
		}
		params = parameters.Params
	}
	conditions := strings.Join(conditionArray, " AND ")
	if countAll {
		err = db.ORM.Model(&torrents).Where(conditions, params...).Count(&count).Error
		if err != nil {
			return
		}
	}
	// TODO: Vulnerable to injections. Use query builder.

	// build custom db query for performance reasons
	dbQuery := "SELECT * FROM torrents"
	if conditions != "" {
		dbQuery = dbQuery + " WHERE " + conditions
	}
	if strings.Contains(conditions, "torrent_name") {
		dbQuery = "WITH t AS (SELECT * FROM torrents WHERE " + conditions + ") SELECT * FROM t"
	}

	if orderBy == "" { // default OrderBy
		orderBy = "torrent_id DESC"
	}
	dbQuery = dbQuery + " ORDER BY " + orderBy
	if limit != 0 || offset != 0 { // if limits provided
		dbQuery = dbQuery + " LIMIT " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(offset)
	}
	err = db.ORM.Raw(dbQuery, params...).Find(&torrents).Error
	return
}

/* Functions to simplify the get parameters of the main function
 *
 * Get Torrents with where parameters and limits, order by default
 */
func GetTorrents(parameters WhereParams, limit int, offset int) ([]model.Torrents, int, error) {
	return GetTorrentsOrderBy(&parameters, "", limit, offset)
}

/* Get Torrents with where parameters but no limit and order by default (get all the torrents corresponding in the db)
 */
func GetTorrentsDB(parameters WhereParams) ([]model.Torrents, int, error) {
	return GetTorrentsOrderBy(&parameters, "", 0, 0)
}

/* Function to get all torrents
 */

func GetAllTorrentsOrderBy(orderBy string, limit int, offset int) ([]model.Torrents, int, error) {
	return GetTorrentsOrderBy(nil, orderBy, limit, offset)
}

func GetAllTorrents(limit int, offset int) ([]model.Torrents, int, error) {
	return GetTorrentsOrderBy(nil, "", limit, offset)
}

func GetAllTorrentsDB() ([]model.Torrents, int, error) {
	return GetTorrentsOrderBy(nil, "", 0, 0)
}

func CreateWhereParams(conditions string, params ...string) WhereParams {
	whereParams := WhereParams{
		Conditions: conditions,
		Params:     make([]interface{}, len(params)),
	}
	for i := range params {
		whereParams.Params[i] = params[i]
	}
	return whereParams
}
