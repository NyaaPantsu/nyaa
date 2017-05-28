package torrentService

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	elastic "gopkg.in/olivere/elastic.v5"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/db"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/service"
	"github.com/NyaaPantsu/nyaa/util"
	"github.com/NyaaPantsu/nyaa/util/log"
)

/* Function to interact with Models
 *
 * Get the torrents with where clause
 *
 */

// GetFeeds : don't need raw SQL once we get MySQL
func GetFeeds() (result []model.Feed, err error) {
	result = make([]model.Feed, 0, 50)
	rows, err := db.ORM.DB().
		Query(
			"SELECT `torrent_id` AS `id`, `torrent_name` AS `name`, `torrent_hash` AS `hash`, `timestamp` FROM `" + config.TorrentsTableName +
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

// GetTorrentByID : get a torrent with its id
func GetTorrentByID(id string) (torrent model.Torrent, err error) {
	// Postgres DB integer size is 32-bit
	idInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return
	}

	tmp := db.ORM.Where("torrent_id = ?", id).Preload("Comments")
	if idInt > config.LastOldTorrentID {
		tmp = tmp.Preload("FileList")
	}
	if idInt <= config.LastOldTorrentID && !config.IsSukebei() {
		// only preload old comments if they could actually exist
		tmp = tmp.Preload("OldComments")
	}
	err = tmp.Error
	if err != nil {
		return
	}
	if tmp.Find(&torrent).RecordNotFound() {
		err = errors.New("Article is not found")
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

// GetRawTorrentByID : Get torrent with id without user or comments
// won't fetch user or comments
func GetRawTorrentByID(id uint) (torrent model.Torrent, err error) {
	err = nil
	if db.ORM.Where("torrent_id = ?", id).Find(&torrent).RecordNotFound() {
		err = errors.New("Torrent is not found")
	}
	return
}

// GetRawTorrentByHash : Get torrent with id without user or comments
// won't fetch user or comments
func GetRawTorrentByHash(hash string) (torrent model.Torrent, err error) {
	err = nil
	if db.ORM.Where("torrent_hash = ?", hash).Find(&torrent).RecordNotFound() {
		err = errors.New("Torrent is not found")
	}
	return
}

// GetTorrentsOrderByNoCount : Get torrents based on search without counting and user
func GetTorrentsOrderByNoCount(parameters *serviceBase.WhereParams, orderBy string, limit int, offset int) (torrents []model.Torrent, err error) {
	torrents, _, err = getTorrentsOrderBy(parameters, orderBy, limit, offset, false, false, false)
	return
}

// GetTorrentsOrderBy : Get torrents based on search without user
func GetTorrentsOrderBy(parameters *serviceBase.WhereParams, orderBy string, limit int, offset int) (torrents []model.Torrent, count int, err error) {
	torrents, count, err = getTorrentsOrderBy(parameters, orderBy, limit, offset, true, false, false)
	return
}

// GetTorrentsWithUserOrderBy : Get torrents based on search with user
func GetTorrentsWithUserOrderBy(parameters *serviceBase.WhereParams, orderBy string, limit int, offset int) (torrents []model.Torrent, count int, err error) {
	torrents, count, err = getTorrentsOrderBy(parameters, orderBy, limit, offset, true, true, false)
	return
}

func getTorrentsOrderBy(parameters *serviceBase.WhereParams, orderBy string, limit int, offset int, countAll bool, withUser bool, deleted bool) (
	torrents []model.Torrent, count int, err error,
) {
	var conditionArray []string
	var params []interface{}
	if parameters != nil { // if there is where parameters
		if len(parameters.Conditions) > 0 {
			conditionArray = append(conditionArray, parameters.Conditions)
		}
		params = parameters.Params
	}
	if !deleted {
		conditionArray = append(conditionArray, "deleted_at IS NULL")
	} else {
		conditionArray = append(conditionArray, "deleted_at NOT NULL")
	}

	conditions := strings.Join(conditionArray, " AND ")

	if countAll {
		err = db.ORM.Unscoped().Model(&torrents).Where(conditions, params...).Count(&count).Error
		if err != nil {
			return
		}
	}

	// build custom db query for performance reasons
	dbQuery := "SELECT * FROM " + config.TorrentsTableName
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
	dbQ := db.ORM
	if withUser {
		dbQ = dbQ.Preload("Uploader")
	}
	if countAll {
		dbQ = dbQ.Preload("Comments")
	}
	err = dbQ.Raw(dbQuery, params...).Find(&torrents).Error
	return
}

// GetTorrents obtain a list of torrents matching 'parameters' from the
// database. The list will be of length 'limit' and in default order.
// GetTorrents returns the first records found. Later records may be retrieved
// by providing a positive 'offset'
func GetTorrents(parameters serviceBase.WhereParams, limit int, offset int) ([]model.Torrent, int, error) {
	return GetTorrentsOrderBy(&parameters, "", limit, offset)
}

// GetTorrentsDB : Get Torrents with where parameters but no limit and order by default (get all the torrents corresponding in the db)
func GetTorrentsDB(parameters serviceBase.WhereParams) ([]model.Torrent, int, error) {
	return GetTorrentsOrderBy(&parameters, "", 0, 0)
}

// GetAllTorrentsOrderBy : Get all torrents ordered by parameters
func GetAllTorrentsOrderBy(orderBy string, limit int, offset int) ([]model.Torrent, int, error) {
	return GetTorrentsOrderBy(nil, orderBy, limit, offset)
}

// GetAllTorrents : Get all torrents without order
func GetAllTorrents(limit int, offset int) ([]model.Torrent, int, error) {
	return GetTorrentsOrderBy(nil, "", limit, offset)
}

// GetAllTorrentsDB : Get all torrents
func GetAllTorrentsDB() ([]model.Torrent, int, error) {
	return GetTorrentsOrderBy(nil, "", 0, 0)
}

// DeleteTorrent : delete a torrent based on id
func DeleteTorrent(id string) (int, error) {
	var torrent model.Torrent
	if db.ORM.First(&torrent, id).RecordNotFound() {
		return http.StatusNotFound, errors.New("Torrent is not found")
	}
	if db.ORM.Delete(&torrent).Error != nil {
		return http.StatusInternalServerError, errors.New("Torrent was not deleted")
	}

	// TODO Don't create a new client for each request
	client, err := elastic.NewClient()
	if err == nil {
		err = torrent.DeleteFromESIndex(client)
		if err == nil {
			log.Infof("Successfully deleted torrent to ES index.")
		} else {
			log.Errorf("Unable to delete torrent to ES index: %s", err)
		}
	} else {
		log.Errorf("Unable to create elasticsearch client: %s", err)
	}
	return http.StatusOK, nil
}

// DefinitelyDeleteTorrent : deletes definitely a torrent based on id
func DefinitelyDeleteTorrent(id string) (int, error) {
	var torrent model.Torrent
	if db.ORM.Unscoped().Model(&torrent).First(&torrent, id).RecordNotFound() {
		return http.StatusNotFound, errors.New("Torrent is not found")
	}
	if db.ORM.Unscoped().Model(&torrent).Delete(&torrent).Error != nil {
		return http.StatusInternalServerError, errors.New("Torrent was not deleted")
	}

	// TODO Don't create a new client for each request
	client, err := elastic.NewClient()
	if err == nil {
		err = torrent.DeleteFromESIndex(client)
		if err == nil {
			log.Infof("Successfully deleted torrent to ES index.")
		} else {
			log.Errorf("Unable to delete torrent to ES index: %s", err)
		}
	} else {
		log.Errorf("Unable to create elasticsearch client: %s", err)
	}
	return http.StatusOK, nil
}

// ToggleBlockTorrent ; Lock/Unlock a torrent based on id
func ToggleBlockTorrent(id string) (model.Torrent, int, error) {
	var torrent model.Torrent
	if db.ORM.Unscoped().Model(&torrent).First(&torrent, id).RecordNotFound() {
		return torrent, http.StatusNotFound, errors.New("Torrent is not found")
	}
	if torrent.Status == model.TorrentStatusBlocked {
		torrent.Status = model.TorrentStatusNormal
	} else {
		torrent.Status = model.TorrentStatusBlocked
	}
	if db.ORM.Unscoped().Model(&torrent).UpdateColumn(&torrent).Error != nil {
		return torrent, http.StatusInternalServerError, errors.New("Torrent was not updated")
	}
	return torrent, http.StatusOK, nil
}

// UpdateTorrent : Update a torrent based on model
func UpdateTorrent(torrent model.Torrent) (int, error) {
	if db.ORM.Model(&torrent).UpdateColumn(&torrent).Error != nil {
		return http.StatusInternalServerError, errors.New("Torrent was not updated")
	}

	// TODO Don't create a new client for each request
	client, err := elastic.NewClient()
	if err == nil {
		err = torrent.AddToESIndex(client)
		if err == nil {
			log.Infof("Successfully updated torrent to ES index.")
		} else {
			log.Errorf("Unable to update torrent to ES index: %s", err)
		}
	} else {
		log.Errorf("Unable to create elasticsearch client: %s", err)
	}

	return http.StatusOK, nil
}

// GetDeletedTorrents : Gets deleted torrents based on search params
func GetDeletedTorrents(parameters *serviceBase.WhereParams, orderBy string, limit int, offset int) (torrents []model.Torrent, count int, err error) {
	torrents, count, err = getTorrentsOrderBy(parameters, orderBy, limit, offset, true, true, true)
	return
}
