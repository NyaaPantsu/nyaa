package torrents

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	elastic "gopkg.in/olivere/elastic.v5"

	"time"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/utils/cache"
)

// Query : Interface to pass for torrents query
type Query interface {
	String() string
	ToDBQuery() (string, []interface{})
	ToESQuery(client *elastic.Client) (*elastic.SearchService, error)
	Append(string, ...interface{})
	Prepend(string, ...interface{})
}

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

	tmp := models.ORM.Where("torrent_id = ?", id).Preload("Scrape").Preload("Uploader").Preload("Comments")
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

// FindUnscopeByID : Get torrent with ID deleted or not
func FindUnscopeByID(id uint) (torrent models.Torrent, err error) {
	err = nil
	if models.ORM.Unscoped().Where("torrent_id = ?", id).Preload("Uploader").Find(&torrent).RecordNotFound() {
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
func FindOrderByNoCount(parameters Query, orderBy string, limit int, offset int) (torrents []models.Torrent, err error) {
	torrents, _, err = findOrderBy(parameters, orderBy, limit, offset, false, false, false)
	return
}

// FindOrderBy : Get torrents based on search without user
func FindOrderBy(parameters Query, orderBy string, limit int, offset int) (torrents []models.Torrent, count int, err error) {
	torrents, count, err = findOrderBy(parameters, orderBy, limit, offset, true, false, false)
	return
}

// FindWithUserOrderBy : Get torrents based on search with user
func FindWithUserOrderBy(parameters Query, orderBy string, limit int, offset int) (torrents []models.Torrent, count int, err error) {
	torrents, count, err = findOrderBy(parameters, orderBy, limit, offset, true, true, false)
	return
}

func findOrderBy(parameters Query, orderBy string, limit int, offset int, countAll bool, withUser bool, deleted bool) (
	torrents []models.Torrent, count int, err error,
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
	if !deleted {
		conditionArray = append(conditionArray, "deleted_at IS NULL")
	} else {
		conditionArray = append(conditionArray, "deleted_at IS NOT NULL")
	}

	conditions := strings.Join(conditionArray, " AND ")

	if countAll {
		err = models.ORM.Unscoped().Model(&torrents).Where(conditions, params...).Count(&count).Error
		if err != nil {
			return
		}
	}

	// build custom db query for performance reasons
	dbQuery := models.ORM.Unscoped().Joins("LEFT JOIN " + config.Get().Models.ScrapeTableName + " ON " + config.Get().Models.TorrentsTableName + ".torrent_id = " + config.Get().Models.ScrapeTableName + ".torrent_id").Preload("Scrape").Preload("FileList")
	if withUser {
		dbQuery = dbQuery.Preload("Uploader")
	}
	if countAll {
		dbQuery = dbQuery.Preload("Comments").Preload("OldComments")
	}

	if conditions != "" {
		dbQuery = dbQuery.Where(conditions, params...)
	}

	if orderBy == "" { // default OrderBy
		orderBy = "date DESC"
	}
	dbQuery = dbQuery.Order(orderBy)
	if limit != 0 || offset != 0 { // if limits provided
		dbQuery = dbQuery.Limit(strconv.Itoa(limit)).Offset(strconv.Itoa(offset))
	}

	err = dbQuery.Find(&torrents).Error
	return
}

// Find obtain a list of torrents matching 'parameters' from the
// database. The list will be of length 'limit' and in default order.
// GetTorrents returns the first records found. Later records may be retrieved
// by providing a positive 'offset'
func Find(parameters Query, limit int, offset int) ([]models.Torrent, int, error) {
	return FindOrderBy(parameters, "", limit, offset)
}

// FindDB : Get Torrents with where parameters but no limit and order by default (get all the torrents corresponding in the db)
func FindDB(parameters Query) ([]models.Torrent, int, error) {
	return FindOrderBy(parameters, "", 0, 0)
}

// FindAllOrderBy : Get all torrents ordered by parameters
func FindAllOrderBy(orderBy string, limit int, offset int) ([]models.Torrent, int, error) {
	return FindOrderBy(nil, orderBy, limit, offset)
}

// FindAllForAdminsOrderBy : Get all torrents ordered by parameters
func FindAllForAdminsOrderBy(orderBy string, limit int, offset int) ([]models.Torrent, int, error) {
    return findOrderBy(nil, orderBy, limit, offset, true, true, false)
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
func FindDeleted(parameters Query, orderBy string, limit int, offset int) (torrents []models.Torrent, count int, err error) {
	torrents, count, err = findOrderBy(parameters, orderBy, limit, offset, true, true, true)
	return
}

// GetIDs : returns an array of id
func GetIDs(parameters Query) ([]uint, error) {
	var ids []uint
	err := models.ORM.Select("torrent_id").Where(parameters.ToDBQuery()).Find(&ids).Error
	return ids, err
}
