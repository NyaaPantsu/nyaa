package torrents

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"time"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/utils/cache"
	"github.com/NyaaPantsu/nyaa/utils/search/structs"
)

/* Function to interact with Models
 *
 * Get the torrents with where clause
 *
 */

// FindByID : get a torrent with its id
func FindByID(id uint) (*models.Torrent, error) {
	torrent := &models.Torrent{ID: id}
	var err error
	if found, ok := cache.C.Get(torrent.Identifier()); ok {
		return found.(*models.Torrent), nil

	}

	tmp := models.ORM.Where("torrent_id = ?", id).Preload("Scrape").Preload("Comments")
	if id > config.Get().Models.LastOldTorrentID {
		tmp = tmp.Preload("FileList")
	}
	if id <= config.Get().Models.LastOldTorrentID && !config.IsSukebei() {
		// only preload old comments if they could actually exist
		tmp = tmp.Preload("OldComments")
	}
	err = tmp.Error
	if err != nil {
		return torrent, err
	}
	if tmp.Find(torrent).RecordNotFound() {
		err = errors.New("Article is not found")
		return torrent, err
	}
	torrent.ParseLanguages()
	// GORM relly likes not doing its job correctly
	// (or maybe I'm just retarded)
	torrent.OldUploader = ""
	if torrent.ID <= config.Get().Models.LastOldTorrentID && torrent.UploaderID == 0 {
		var tmp models.UserUploadsOld
		if !models.ORM.Where("torrent_id = ?", torrent.ID).Find(&tmp).RecordNotFound() {
			torrent.OldUploader = tmp.Username
		}
	}
	for i := range torrent.Comments {
		torrent.Comments[i].User = new(models.User)
		models.ORM.Where("user_id = ?", torrent.Comments[i].UserID).Find(torrent.Comments[i].User)
	}
	cache.C.Set(torrent.Identifier(), torrent, 5*time.Minute)
	return torrent, nil
}

// FindRawByID : Get torrent with id without user or comments
// won't fetch user or comments
func FindRawByID(id uint) (torrent models.Torrent, err error) {
	err = nil
	if models.ORM.Where("torrent_id = ?", id).Find(&torrent).RecordNotFound() {
		err = errors.New("Torrent is not found")
	}
	torrent.ParseLanguages()
	return
}

// FindRawByHash : Get torrent with id without user or comments
// won't fetch user or comments
func FindRawByHash(hash string) (torrent models.Torrent, err error) {
	err = nil
	if models.ORM.Where("torrent_hash = ?", hash).Find(&torrent).RecordNotFound() {
		err = errors.New("Torrent is not found")
	}
	torrent.ParseLanguages()
	return
}

// FindOrderByNoCount : Get torrents based on search without counting and user
func FindOrderByNoCount(parameters *structs.WhereParams, orderBy string, limit int, offset int) (torrents []models.Torrent, err error) {
	torrents, _, err = findOrderBy(parameters, orderBy, limit, offset, false, false, false)
	return
}

// FindOrderBy : Get torrents based on search without user
func FindOrderBy(parameters *structs.WhereParams, orderBy string, limit int, offset int) (torrents []models.Torrent, count int, err error) {
	torrents, count, err = findOrderBy(parameters, orderBy, limit, offset, true, false, false)
	return
}

// FindWithUserOrderBy : Get torrents based on search with user
func FindWithUserOrderBy(parameters *structs.WhereParams, orderBy string, limit int, offset int) (torrents []models.Torrent, count int, err error) {
	torrents, count, err = findOrderBy(parameters, orderBy, limit, offset, true, true, false)
	return
}

func findOrderBy(parameters *structs.WhereParams, orderBy string, limit int, offset int, countAll bool, withUser bool, deleted bool) (
	torrents []models.Torrent, count int, err error,
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
	/*if found, ok := cache.C.Get(fmt.Sprintf("%v", parameters)); ok {
		torrentCache := found.(*structs.TorrentCache)
		torrents = torrentCache.Torrents
		count = torrentCache.Count
		return
	}*/
	if countAll {
		err = models.ORM.Unscoped().Model(&torrents).Where(conditions, params...).Count(&count).Error
		if err != nil {
			return
		}
	}

	// build custom db query for performance reasons
	dbQuery := "SELECT * FROM " + config.Get().Models.TorrentsTableName
	if conditions != "" {
		dbQuery = dbQuery + " WHERE " + conditions
	}

	if orderBy == "" { // default OrderBy
		orderBy = "torrent_id DESC"
	}
	dbQuery = dbQuery + " ORDER BY " + orderBy
	if limit != 0 || offset != 0 { // if limits provided
		dbQuery = dbQuery + " LIMIT " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(offset)
	}
	dbQ := models.ORM.Preload("Scrape")
	if withUser {
		dbQ = dbQ.Preload("Uploader")
	}
	if countAll {
		dbQ = dbQ.Preload("Comments")
	}
	err = dbQ.Preload("FileList").Raw(dbQuery, params...).Find(&torrents).Error
	// cache.C.Set(fmt.Sprintf("%v", parameters), &structs.TorrentCache{torrents, count}, 5*time.Minute) // Cache shouldn't be done here but in search util
	return
}

// Find obtain a list of torrents matching 'parameters' from the
// database. The list will be of length 'limit' and in default order.
// GetTorrents returns the first records found. Later records may be retrieved
// by providing a positive 'offset'
func Find(parameters structs.WhereParams, limit int, offset int) ([]models.Torrent, int, error) {
	return FindOrderBy(&parameters, "", limit, offset)
}

// FindDB : Get Torrents with where parameters but no limit and order by default (get all the torrents corresponding in the db)
func FindDB(parameters structs.WhereParams) ([]models.Torrent, int, error) {
	return FindOrderBy(&parameters, "", 0, 0)
}

// FindAllOrderBy : Get all torrents ordered by parameters
func FindAllOrderBy(orderBy string, limit int, offset int) ([]models.Torrent, int, error) {
	return FindOrderBy(nil, orderBy, limit, offset)
}

// FindAll : Get all torrents without order
func FindAll(limit int, offset int) ([]models.Torrent, int, error) {
	return FindOrderBy(nil, "", limit, offset)
}

// GetAllInDB : Get all torrents
func GetAllInDB() ([]models.Torrent, int, error) {
	return FindOrderBy(nil, "", 0, 0)
}

// ToggleBlock ; Lock/Unlock a torrent based on id
func ToggleBlock(id uint) (models.Torrent, int, error) {
	var torrent models.Torrent
	if models.ORM.Unscoped().Model(&torrent).First(&torrent, id).RecordNotFound() {
		return torrent, http.StatusNotFound, errors.New("Torrent is not found")
	}
	torrent.ParseLanguages()
	if torrent.Status == models.TorrentStatusBlocked {
		torrent.Status = models.TorrentStatusNormal
	} else {
		torrent.Status = models.TorrentStatusBlocked
	}
	if models.ORM.Unscoped().Model(&torrent).UpdateColumn(&torrent).Error != nil {
		return torrent, http.StatusInternalServerError, errors.New("Torrent was not updated")
	}
	return torrent, http.StatusOK, nil
}

// FindDeleted : Gets deleted torrents based on search params
func FindDeleted(parameters *structs.WhereParams, orderBy string, limit int, offset int) (torrents []models.Torrent, count int, err error) {
	torrents, count, err = findOrderBy(parameters, orderBy, limit, offset, true, true, true)
	return
}
