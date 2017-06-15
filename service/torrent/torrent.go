package torrentService

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/db"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/service"
	"github.com/NyaaPantsu/nyaa/service/activity"
	"github.com/NyaaPantsu/nyaa/service/notifier"
	"github.com/NyaaPantsu/nyaa/service/user"
	"github.com/NyaaPantsu/nyaa/service/user/permission"
	"github.com/NyaaPantsu/nyaa/util/log"
	"github.com/NyaaPantsu/nyaa/util/publicSettings"
	"github.com/gorilla/mux"
)

/* Function to interact with Models
 *
 * Get the torrents with where clause
 *
 */

// GetTorrentByID : get a torrent with its id
func GetTorrentByID(id string) (torrent model.Torrent, err error) {
	// Postgres DB integer size is 32-bit
	idInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return
	}

	tmp := db.ORM.Where("torrent_id = ?", id).Preload("Comments")
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
	torrent.Uploader = new(model.User)
	db.ORM.Where("user_id = ?", torrent.UploaderID).Find(torrent.Uploader)
	torrent.OldUploader = ""
	if torrent.ID <= config.Conf.Models.LastOldTorrentID && torrent.UploaderID == 0 {
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
	dbQuery := "SELECT * FROM " + config.Conf.Models.TorrentsTableName
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
	err = dbQ.Preload("FileList").Raw(dbQuery, params...).Find(&torrents).Error
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
func DeleteTorrent(id string) (*model.Torrent, int, error) {
	var torrent model.Torrent
	if db.ORM.First(&torrent, id).RecordNotFound() {
		return &torrent, http.StatusNotFound, errors.New("Torrent is not found")
	}
	if db.ORM.Delete(&torrent).Error != nil {
		return &torrent, http.StatusInternalServerError, errors.New("Torrent was not deleted")
	}

	if db.ElasticSearchClient != nil {
		err := torrent.DeleteFromESIndex(db.ElasticSearchClient)
		if err == nil {
			log.Infof("Successfully deleted torrent to ES index.")
		} else {
			log.Errorf("Unable to delete torrent to ES index: %s", err)
		}
	}
	return &torrent, http.StatusOK, nil
}

// DefinitelyDeleteTorrent : deletes definitely a torrent based on id
func DefinitelyDeleteTorrent(id string) (*model.Torrent, int, error) {
	var torrent model.Torrent
	if db.ORM.Unscoped().Model(&torrent).First(&torrent, id).RecordNotFound() {
		return &torrent, http.StatusNotFound, errors.New("Torrent is not found")
	}
	if db.ORM.Unscoped().Model(&torrent).Delete(&torrent).Error != nil {
		return &torrent, http.StatusInternalServerError, errors.New("Torrent was not deleted")
	}

	if db.ElasticSearchClient != nil {
		err := torrent.DeleteFromESIndex(db.ElasticSearchClient)
		if err == nil {
			log.Infof("Successfully deleted torrent to ES index.")
		} else {
			log.Errorf("Unable to delete torrent to ES index: %s", err)
		}
	}
	return &torrent, http.StatusOK, nil
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
func UpdateTorrent(torrent *model.Torrent) (int, error) {
	if db.ORM.Model(torrent).UpdateColumn(torrent).Error != nil {
		return http.StatusInternalServerError, errors.New("Torrent was not updated")
	}

	// TODO Don't create a new client for each request
	if db.ElasticSearchClient != nil {
		err := torrent.AddToESIndex(db.ElasticSearchClient)
		if err == nil {
			log.Infof("Successfully updated torrent to ES index.")
		} else {
			log.Errorf("Unable to update torrent to ES index: %s", err)
		}
	}

	return http.StatusOK, nil
}

// UpdateUnscopeTorrent : Update a torrent based on model
func UpdateUnscopeTorrent(torrent *model.Torrent) (int, error) {
	if db.ORM.Unscoped().Model(torrent).UpdateColumn(torrent).Error != nil {
		return http.StatusInternalServerError, errors.New("Torrent was not updated")
	}

	// TODO Don't create a new client for each request
	if db.ElasticSearchClient != nil {
		err := torrent.AddToESIndex(db.ElasticSearchClient)
		if err == nil {
			log.Infof("Successfully updated torrent to ES index.")
		} else {
			log.Errorf("Unable to update torrent to ES index: %s", err)
		}
	}

	return http.StatusOK, nil
}

// GetDeletedTorrents : Gets deleted torrents based on search params
func GetDeletedTorrents(parameters *serviceBase.WhereParams, orderBy string, limit int, offset int) (torrents []model.Torrent, count int, err error) {
	torrents, count, err = getTorrentsOrderBy(parameters, orderBy, limit, offset, true, true, true)
	return
}

// ExistOrDelete : Check if a torrent exist with the same hash and if it can be replaced, it is replaced
func ExistOrDelete(hash string, user *model.User) error {
	torrentIndb := model.Torrent{}
	db.ORM.Unscoped().Model(&model.Torrent{}).Where("torrent_hash = ?", hash).First(&torrentIndb)
	if torrentIndb.ID > 0 {
		if userPermission.CurrentUserIdentical(user, torrentIndb.UploaderID) && torrentIndb.IsDeleted() && !torrentIndb.IsBlocked() { // if torrent is not locked and is deleted and the user is the actual owner
			torrent, _, err := DefinitelyDeleteTorrent(strconv.Itoa(int(torrentIndb.ID)))
			if err != nil {
				return err
			}
			activity.Log(&model.User{}, torrent.Identifier(), "delete", "torrent_deleted_by", strconv.Itoa(int(torrent.ID)), torrent.Uploader.Username, user.Username)
		} else {
			return errors.New("Torrent already in database")
		}
	}
	return nil
}

// NewTorrentEvent : Should be called when you create a new torrent
func NewTorrentEvent(router *mux.Router, user *model.User, torrent *model.Torrent) error {
	url, err := router.Get("view_torrent").URL("id", strconv.FormatUint(uint64(torrent.ID), 10))
	if err != nil {
		return err
	}
	if user.ID > 0 && config.Conf.Users.DefaultUserSettings["new_torrent"] { // If we are a member and notifications for new torrents are enabled
		userService.GetFollowers(user) // We populate the liked field for users
		if len(user.Followers) > 0 {   // If we are followed by at least someone
			for _, follower := range user.Followers {
				follower.ParseSettings() // We need to call it before checking settings
				if follower.Settings.Get("new_torrent") {
					T, _, _ := publicSettings.TfuncAndLanguageWithFallback(follower.Language, follower.Language) // We need to send the notification to every user in their language
					notifierService.NotifyUser(&follower, torrent.Identifier(), fmt.Sprintf(T("new_torrent_uploaded"), torrent.Name, user.Username), url.String(), follower.Settings.Get("new_torrent_email"))
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

// TorrentsToAPI : Map Torrents for API usage without reallocations
func TorrentsToAPI(t []model.Torrent) []model.TorrentJSON {
	json := make([]model.TorrentJSON, len(t))
	for i := range t {
		json[i] = t[i].ToJSON()
		uploaderID, username := HideTorrentUser(json[i].UploaderID, string(json[i].UploaderName), json[i].Hidden)
		json[i].UploaderName = template.HTML(username)
		json[i].UploaderID = uploaderID
	}
	return json
}
