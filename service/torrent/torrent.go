package torrentService

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/db"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/service"
	"github.com/NyaaPantsu/nyaa/util"
)

/* Function to interact with Models
 *
 * Get the torrents with where clause
 *
 */

// don't need raw SQL once we get MySQL
func GetFeeds() (result []model.Feed, err error) {
	result = make([]model.Feed, 0, 50)
	rows, err := db.ORM.DB().
		Query(
			"SELECT `torrent_id` AS `id`, `torrent_name` AS `name`, `torrent_hash` AS `hash`, `timestamp` FROM `" + config.TableName +
				"` ORDER BY `timestamp` desc LIMIT 50")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := model.Feed{}
		err = rows.Scan(&item.ID, &item.Name, &item.Hash, &item.Timestamp)
		if err != nil {
			return
		}
		magnet := util.InfoHashToMagnet(strings.TrimSpace(item.Hash), item.Name, config.Trackers...)
		item.Magnet = magnet
		// TODO: memory hog
		result = append(result, item)
	}
	err = rows.Err()
	return
}

func GetTorrentById(id string) (torrent model.Torrent, err error) {
	// Postgres DB integer size is 32-bit
	id_int, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return
	}

	tmp := db.ORM.Table(config.TableName).Where("torrent_id = ?", id).Preload("Comments").Preload("FileList")
	err = tmp.Error
	if err != nil {
		return
	}
	if id_int <= config.LastOldTorrentID {
		// only preload old comments if they could actually exist
		tmp = tmp.Preload("OldComments")
	}
	if tmp.Find(&torrent).RecordNotFound() {
		err = errors.New("Article is not found.")
		return
	}
	// GORM relly likes not doing its job correctly
	// (or maybe I'm just retarded)
	torrent.Uploader = new(model.User)
	db.ORM.Where("user_id = ?", torrent.UploaderID).Find(torrent.Uploader)
	torrent.OldUploader = ""
	if torrent.ID <= config.LastOldTorrentID && torrent.UploaderID == 0 {
		var tmp model.UserUploadsOld
		if !db.ORM.Where("torrent_id = ?", torrent.ID).Find(&tmp).RecordNotFound() {
			torrent.OldUploader = tmp.Username
		}
	}
	for i := range torrent.Comments {
		torrent.Comments[i].User = new(model.User)
		err = db.ORM.Where("user_id = ?", torrent.Comments[i].UserID).Find(torrent.Comments[i].User).Error
		if err != nil {
			return
		}
	}

	return
}

// won't fetch user or comments
func GetRawTorrentById(id uint) (torrent model.Torrent, err error) {
	err = nil
	if db.ORM.Table(config.TableName).Table(config.TableName).Where("torrent_id = ?", id).Find(&torrent).RecordNotFound() {
		err = errors.New("Article is not found.")
	}
	return
}

func GetTorrentsOrderByNoCount(parameters *serviceBase.WhereParams, orderBy string, limit int, offset int) (torrents []model.Torrent, err error) {
	torrents, _, err = getTorrentsOrderBy(parameters, orderBy, limit, offset, false)
	return
}

func GetTorrentsOrderBy(parameters *serviceBase.WhereParams, orderBy string, limit int, offset int) (torrents []model.Torrent, count int, err error) {
	torrents, count, err = getTorrentsOrderBy(parameters, orderBy, limit, offset, true)
	return
}

func getTorrentsOrderBy(parameters *serviceBase.WhereParams, orderBy string, limit int, offset int, countAll bool) (
	torrents []model.Torrent, count int, err error,
) {
	var conditionArray []string
	conditionArray = append(conditionArray, "deleted_at IS NULL")
	var params []interface{}
	if parameters != nil { // if there is where parameters
		if len(parameters.Conditions) > 0 {
			conditionArray = append(conditionArray, parameters.Conditions)
		}
		params = parameters.Params
	}
	conditions := strings.Join(conditionArray, " AND ")
	if countAll {
		// FIXME: `deleted_at IS NULL` is duplicate in here because GORM handles this for us
		err = db.ORM.Model(&torrents).Table(config.TableName).Where(conditions, params...).Count(&count).Error
		if err != nil {
			return
		}
	}
	// TODO: Vulnerable to injections. Use query builder. (is it?)

	// build custom db query for performance reasons
	dbQuery := "SELECT * FROM " + config.TableName
	if conditions != "" {
		dbQuery = dbQuery + " WHERE " + conditions
	}
	/* This makes all queries take roughly the same amount of time (lots)...
	if strings.Contains(conditions, "torrent_name") && offset > 0 {
		dbQuery = "WITH t AS (SELECT * FROM torrents WHERE " + conditions + ") SELECT * FROM t"
	}*/

	if orderBy == "" { // default OrderBy
		orderBy = "torrent_id DESC"
	}
	dbQuery = dbQuery + " ORDER BY " + orderBy
	if limit != 0 || offset != 0 { // if limits provided
		dbQuery = dbQuery + " LIMIT " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(offset)
	}
	err = db.ORM.Preload("Comments").Raw(dbQuery, params...).Find(&torrents).Error
	return
}

// GetTorrents obtain a list of torrents matching 'parameters' from the
// database. The list will be of length 'limit' and in default order.
// GetTorrents returns the first records found. Later records may be retrieved
// by providing a positive 'offset'
func GetTorrents(parameters serviceBase.WhereParams, limit int, offset int) ([]model.Torrent, int, error) {
	return GetTorrentsOrderBy(&parameters, "", limit, offset)
}

// Get Torrents with where parameters but no limit and order by default (get all the torrents corresponding in the db)
func GetTorrentsDB(parameters serviceBase.WhereParams) ([]model.Torrent, int, error) {
	return GetTorrentsOrderBy(&parameters, "", 0, 0)
}

func GetAllTorrentsOrderBy(orderBy string, limit int, offset int) ([]model.Torrent, int, error) {
	return GetTorrentsOrderBy(nil, orderBy, limit, offset)
}

func GetAllTorrents(limit int, offset int) ([]model.Torrent, int, error) {
	return GetTorrentsOrderBy(nil, "", limit, offset)
}

func GetAllTorrentsDB() ([]model.Torrent, int, error) {
	return GetTorrentsOrderBy(nil, "", 0, 0)
}

func DeleteTorrent(id string) (int, error) {
	var torrent model.Torrent
	if db.ORM.Table(config.TableName).First(&torrent, id).RecordNotFound() {
		return http.StatusNotFound, errors.New("Torrent is not found.")
	}
	if db.ORM.Table(config.TableName).Delete(&torrent).Error != nil {
		return http.StatusInternalServerError, errors.New("Torrent is not deleted.")
	}
	return http.StatusOK, nil
}

func UpdateTorrent(torrent model.Torrent) (int, error) {
	if db.ORM.Table(config.TableName).Save(torrent).Error != nil {
		return http.StatusInternalServerError, errors.New("Torrent is not updated.")
	}

	return http.StatusOK, nil
}
