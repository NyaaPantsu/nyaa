package torrents

import (
	"errors"
	"fmt"
	"net/http"
	"nyaa-master/service/userService/userPermission"
	"strconv"
	"strings"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/service"
	"github.com/NyaaPantsu/nyaa/util/log"
	"github.com/NyaaPantsu/nyaa/util/publicSettings"
)

/* Function to interact with Models
 *
 * Get the torrents with where clause
 *
 */

// FindByID : get a torrent with its id
func FindByID(id uint) (torrent models.Torrent, err error) {
	// Postgres DB integer size is 32-bit
	idInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return
	}

	tmp := models.ORM.Where("torrent_id = ?", id).Preload("Scrape").Preload("Comments")
	if idInt > int64(config.Conf.Models.LastOldTorrentID) {
		tmp = tmp.Preload("FileList")
	}
	if idInt <= int64(config.Conf.Models.LastOldTorrentID) && !config.IsSukebei() {
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
	torrent.Uploader = new(models.User)
	models.ORM.Where("user_id = ?", torrent.UploaderID).Find(torrent.Uploader)
	torrent.OldUploader = ""
	if torrent.ID <= config.Conf.Models.LastOldTorrentID && torrent.UploaderID == 0 {
		var tmp models.UserUploadsOld
		if !models.ORM.Where("torrent_id = ?", torrent.ID).Find(&tmp).RecordNotFound() {
			torrent.OldUploader = tmp.Username
		}
	}
	for i := range torrent.Comments {
		torrent.Comments[i].User = new(models.User)
		err = models.ORM.Where("user_id = ?", torrent.Comments[i].UserID).Find(torrent.Comments[i].User).Error
		if err != nil {
			return
		}
	}

	return
}

// FindRawByID : Get torrent with id without user or comments
// won't fetch user or comments
func FindRawByID(id uint) (torrent models.Torrent, err error) {
	err = nil
	if models.ORM.Where("torrent_id = ?", id).Find(&torrent).RecordNotFound() {
		err = errors.New("Torrent is not found")
	}
	return
}

// FindRawByHash : Get torrent with id without user or comments
// won't fetch user or comments
func FindRawByHash(hash string) (torrent models.Torrent, err error) {
	err = nil
	if models.ORM.Where("torrent_hash = ?", hash).Find(&torrent).RecordNotFound() {
		err = errors.New("Torrent is not found")
	}
	return
}

// FindOrderByNoCount : Get torrents based on search without counting and user
func FindOrderByNoCount(parameters *serviceBase.WhereParams, orderBy string, limit int, offset int) (torrents []models.Torrent, err error) {
	torrents, _, err = findOrderBy(parameters, orderBy, limit, offset, false, false, false)
	return
}

// FindOrderBy : Get torrents based on search without user
func FindOrderBy(parameters *serviceBase.WhereParams, orderBy string, limit int, offset int) (torrents []models.Torrent, count int, err error) {
	torrents, count, err = findOrderBy(parameters, orderBy, limit, offset, true, false, false)
	return
}

// FindWithUserOrderBy : Get torrents based on search with user
func FindWithUserOrderBy(parameters *serviceBase.WhereParams, orderBy string, limit int, offset int) (torrents []models.Torrent, count int, err error) {
	torrents, count, err = findOrderBy(parameters, orderBy, limit, offset, true, true, false)
	return
}

func findOrderBy(parameters *serviceBase.WhereParams, orderBy string, limit int, offset int, countAll bool, withUser bool, deleted bool) (
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

	if countAll {
		err = models.ORM.Unscoped().Model(&torrents).Where(conditions, params...).Count(&count).Error
		if err != nil {
			return
		}
	}

	// build custom db query for performance reasons
	dbQuery := "SELECT * FROM " + config.Conf.Models.TorrentsTableName
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
	return
}

// FindRawByID obtain a list of torrents matching 'parameters' from the
// database. The list will be of length 'limit' and in default order.
// GetTorrents returns the first records found. Later records may be retrieved
// by providing a positive 'offset'
func FindRawByID(parameters serviceBase.WhereParams, limit int, offset int) ([]models.Torrent, int, error) {
	return findOrderBy(&parameters, "", limit, offset)
}

// FindDB : Get Torrents with where parameters but no limit and order by default (get all the torrents corresponding in the db)
func FindDB(parameters serviceBase.WhereParams) ([]models.Torrent, int, error) {
	return findOrderBy(&parameters, "", 0, 0)
}

// FindAllOrderBy : Get all torrents ordered by parameters
func FindAllOrderBy(orderBy string, limit int, offset int) ([]models.Torrent, int, error) {
	return findOrderBy(nil, orderBy, limit, offset)
}

// FindAll : Get all torrents without order
func FindAll(limit int, offset int) ([]models.Torrent, int, error) {
	return findOrderBy(nil, "", limit, offset)
}

// GetAllInDB : Get all torrents
func GetAllInDB() ([]models.Torrent, int, error) {
	return findOrderBy(nil, "", 0, 0)
}

// ToggleBlock ; Lock/Unlock a torrent based on id
func ToggleBlock(id uint) (models.Torrent, int, error) {
	var torrent models.Torrent
	if models.ORM.Unscoped().Model(&torrent).First(&torrent, id).RecordNotFound() {
		return torrent, http.StatusNotFound, errors.New("Torrent is not found")
	}
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

// Update : Update a torrent based on model
func Update(torrent *models.Torrent) (int, error) {
	if models.ORM.Model(torrent).UpdateColumn(torrent).Error != nil {
		return http.StatusInternalServerError, errors.New("Torrent was not updated")
	}

	// TODO Don't create a new client for each request
	if models.ElasticSearchClient != nil {
		err := torrent.AddToESIndex(models.ElasticSearchClient)
		if err == nil {
			log.Infof("Successfully updated torrent to ES index.")
		} else {
			log.Errorf("Unable to update torrent to ES index: %s", err)
		}
	}

	return http.StatusOK, nil
}

// UpdateUnscope : Update a torrent based on model
func UpdateUnscope(torrent *models.Torrent) (int, error) {
	if models.ORM.Unscoped().Model(torrent).UpdateColumn(torrent).Error != nil {
		return http.StatusInternalServerError, errors.New("Torrent was not updated")
	}

	// TODO Don't create a new client for each request
	if models.ElasticSearchClient != nil {
		err := torrent.AddToESIndex(models.ElasticSearchClient)
		if err == nil {
			log.Infof("Successfully updated torrent to ES index.")
		} else {
			log.Errorf("Unable to update torrent to ES index: %s", err)
		}
	}

	return http.StatusOK, nil
}

// FindDeleted : Gets deleted torrents based on search params
func FindDeleted(parameters *serviceBase.WhereParams, orderBy string, limit int, offset int) (torrents []models.Torrent, count int, err error) {
	torrents, count, err = findOrderBy(parameters, orderBy, limit, offset, true, true, true)
	return
}

// ExistOrDelete : Check if a torrent exist with the same hash and if it can be replaced, it is replaced
func ExistOrDelete(hash string, user *models.User) error {
	torrentIndb := models.Torrent{}
	models.ORM.Unscoped().Model(&models.Torrent{}).Where("torrent_hash = ?", hash).First(&torrentIndb)
	if torrentInmodels.ID > 0 {
		if userPermission.CurrentUserIdentical(user, torrentInmodels.UploaderID) && torrentInmodels.IsDeleted() && !torrentInmodels.IsBlocked() { // if torrent is not locked and is deleted and the user is the actual owner
			torrent, _, err := DefinitelyDeleteTorrent(strconv.Itoa(int(torrentInmodels.ID)))
			if err != nil {
				return err
			}
			activity.Log(&models.User{}, torrent.Identifier(), "delete", "torrent_deleted_by", strconv.Itoa(int(torrent.ID)), torrent.Uploader.Username, user.Username)
		} else {
			return errors.New("Torrent already in database")
		}
	}
	return nil
}

// NewTorrentEvent : Should be called when you create a new torrent
func NewTorrentEvent(user *models.User, torrent *models.Torrent) error {
	url := "/view/" + strconv.FormatUint(uint64(torrent.ID), 10)
	if user.ID > 0 && config.Conf.Users.DefaultUserSettings["new_torrent"] { // If we are a member and notifications for new torrents are enabled
		userService.GetFollowers(user) // We populate the liked field for users
		if len(user.Followers) > 0 {   // If we are followed by at least someone
			for _, follower := range user.Followers {
				follower.ParseSettings() // We need to call it before checking settings
				if follower.Settings.Get("new_torrent") {
					T, _, _ := publicSettings.TfuncAndLanguageWithFallback(follower.Language, follower.Language) // We need to send the notification to every user in their language
					notifierService.NotifyUser(&follower, torrent.Identifier(), fmt.Sprintf(T("new_torrent_uploaded"), torrent.Name, user.Username), url, follower.Settings.Get("new_torrent_email"))
				}
			}
		}
	}
	return nil
}

// HideTorrentUser : hides a torrent user for hidden torrents
func HideTorrentUser(uploaderID uint, uploaderName string, torrentHidden bool) (uint, string) {
	if torrentHidden {
		return 0, "れんちょん"
	}
	if uploaderID == 0 {
		return 0, uploaderName
	}
	return uploaderID, uploaderName
}
