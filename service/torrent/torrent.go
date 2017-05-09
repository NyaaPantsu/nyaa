package torrentService

import (
	"errors"
	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/util"
	"strconv"
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
			rows.Scan(&item.ID, &item.Name, &item.Hash, &item.Timestamp)
			magnet := util.InfoHashToMagnet(strings.TrimSpace(item.Hash), item.Name, config.Trackers...)
			item.Magnet = magnet
			// memory hog
			result = append(result, item)
		}
		rows.Close()
	}
	return result
}

func GetTorrentById(id string) (model.Torrent, error) {
	var torrent model.Torrent
	intID, err := strconv.Atoi(id)
	if err != nil {
		return torrent, err
	}

	tmp := db.ORM.Where("torrent_id = ?", id).Preload("Comments")
	if intID <= config.LastOldTorrentID {
		// only preload old comments if they could actually exist
		tmp = tmp.Preload("OldComments")
	}
	if tmp.Find(&torrent).RecordNotFound() {
		return torrent, errors.New("article not found")
	}
	// TODO: .Preload("Comments.User") doesn't work
	for i := range torrent.Comments {
		torrent.Comments[i].User = new(model.User)
		db.ORM.Where("user_id = ?", torrent.Comments[i].UserID).Find(torrent.Comments[i].User)
	}

	return torrent, nil
}

func GetTorrentsOrderBy(parameters *WhereParams, orderBy string, limit int, offset int) ([]model.Torrent, int) {
	var torrents []model.Torrent
	var count int
	var conditionArray []string
	if strings.HasPrefix(orderBy, "filesize") {
		// torrents w/ NULL filesize fuck up the sorting on Postgres
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
	db.ORM.Model(&torrents).Where(conditions, params...).Count(&count)

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
	db.ORM.Raw(dbQuery, params...).Find(&torrents)
	return torrents, count
}

// GetTorrents obtain a list of torrents matching 'parameters' from the
// database. The list will be of length 'limit' and in default order.
// GetTorrents returns the first records found. Later records may be retrieved
// by providing a positive 'offset'
func GetTorrents(parameters WhereParams, limit int, offset int) ([]model.Torrent, int) {
	return GetTorrentsOrderBy(&parameters, "", limit, offset)
}

// Get Torrents with where parameters but no limit and order by default (get all the torrents corresponding in the db)
func GetTorrentsDB(parameters WhereParams) ([]model.Torrent, int) {
	return GetTorrentsOrderBy(&parameters, "", 0, 0)
}

func GetAllTorrentsOrderBy(orderBy string, limit int, offset int) ([]model.Torrent, int) {

	return GetTorrentsOrderBy(nil, orderBy, limit, offset)
}

func GetAllTorrents(limit int, offset int) ([]model.Torrent, int) {
	return GetTorrentsOrderBy(nil, "", limit, offset)
}

func GetAllTorrentsDB() ([]model.Torrent, int) {
	return GetTorrentsOrderBy(nil, "", 0, 0)
}

func CreateWhereParams(conditions string, params ...string) WhereParams {
	whereParams := WhereParams{}
	whereParams.Conditions = conditions
	for i := range params {
		whereParams.Params = append(whereParams.Params, params[i])
	}

	return whereParams
}
