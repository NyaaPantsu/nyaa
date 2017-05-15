package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ewhal/nyaa/model"
)

const queryGetAllTorrents = "GetAllTorrents"
const queryGetTorrentByID = "GetTorrentByID"
const queryInsertComment = "InsertComment"
const queryInsertUser = "InsertUser"
const queryInsertTorrentReport = "InsertTorrentReport"
const queryUserFollows = "UserFollows"
const queryDeleteTorrentReportByID = "DeleteTorrentReportByID"
const queryInsertScrape = "InsertScrape"
const queryGetUserByApiToken = "GetUserByApiToken"
const queryGetUserByEmail = "GetUserByEmail"
const queryGetUserByName = "GetUserByName"
const queryGetUserByID = "GetUserByID"
const queryUpdateUser = "UpdateUser"
const queryDeleteUserByID = "DeleteUserByID"
const queryDeleteUserByEmail = "DeleteUserByEmail"
const queryDeleteUserByToken = "DeleteUserByToken"
const queryUserFollowsUpsert = "UserFollowsUpsert"
const queryDeleteUserFollowing = "DeleteUserFollowing"

const torrentSelectColumnsFull = `torrent_id, torrent_name, torrent_hash, category, sub_category, status, date, uploader, downloads, stardom, description, website_link, deleted_at, seeders, leechers, completed, last_scrape`

func scanTorrentColumnsFull(rows *sql.Rows, t *model.Torrent) (err error) {
	var downloads, leechers, seeders, completed, description, websitelink, lastscrape, date, deleted_at interface{}
	err = rows.Scan(&t.ID, &t.Name, &t.Hash, &t.Category, &t.SubCategory, &t.Status, &date, &t.UploaderID, &downloads, &t.Stardom, &description, &websitelink, &deleted_at, &seeders, &leechers, &completed, &lastscrape)
	if downloads != nil {
		t.Downloads = int(downloads.(int64))
	}
	if seeders != nil {
		t.Seeders = uint32(seeders.(int64))
	}
	if leechers != nil {
		t.Leechers = uint32(leechers.(int64))
	}
	if completed != nil {
		t.Completed = uint32(completed.(int64))
	}
	if description != nil {
		t.Description = description.(string)
	}
	if websitelink != nil {
		t.WebsiteLink = websitelink.(string)
	}
	if date != nil {
		t.Date = date.(time.Time)
	}
	if lastscrape != nil {
		t.LastScrape = lastscrape.(time.Time)
	}
	if deleted_at != nil {
		d := deleted_at.(time.Time)
		t.DeletedAt = &d
	}

	return
}

const commentSelectColumnsFull = `comment_id, torrent_id, user_id, content, created_at, updated_at, deleted_at`

func scanCommentColumnsFull(rows *sql.Rows, c *model.Comment) {

}

const torrentReportSelectColumnsFull = `torrent_report_id, type, torrent_id, user_id, created_at`

func scanTorrentReportColumnsFull(rows *sql.Rows, r *model.TorrentReport) {
	rows.Scan(&r.ID, &r.Description, &r.TorrentID, &r.UserID, &r.CreatedAt)
}

const userSelectColumnsFull = `user_id, username, password, email, status, created_at, updated_at, api_token, api_token_expiry, language, md5`

func scanUserColumnsFull(rows *sql.Rows, u *model.User) {
	rows.Scan(&u.ID, &u.Username, &u.Password, &u.Email, &u.Status, &u.CreatedAt, &u.UpdatedAt, &u.ApiToken, &u.ApiTokenExpiry, &u.Language, &u.MD5)

}

var statements = map[string]string{
	queryGetTorrentByID:          fmt.Sprintf("SELECT %s FROM %s WHERE torrent_id = $1 LIMIT 1", torrentSelectColumnsFull, tableTorrents),
	queryGetAllTorrents:          fmt.Sprintf("SELECT %s FROM %s LIMIT $2 OFFSET $1", torrentSelectColumnsFull, tableTorrents),
	queryInsertComment:           fmt.Sprintf("INSERT INTO %s (comment_id, torrent_id, content, created_at) VALUES ($1, $2, $3, $4)", tableComments),
	queryInsertTorrentReport:     fmt.Sprintf("INSERT INTO %s (type, torrent_id, user_id, created_at) VALUES ($1, $2, $3, $4)", tableTorrentReports),
	queryUserFollows:             fmt.Sprintf("SELECT user_id, following FROM %s WHERE user_id = $1 AND following = $1 LIMIT 1", tableUserFollows),
	queryDeleteTorrentReportByID: fmt.Sprintf("DELETE FROM %s WHERE torrent_report_id = $1", tableTorrentReports),
	queryInsertScrape:            fmt.Sprintf("UPDATE %s SET seeders = $1, leechers = $2, completed = $3, last_scrape = $4 WHERE torrent_id = $5", tableTorrents),
	queryGetUserByApiToken:       fmt.Sprintf("SELECT %s FROM %s WHERE api_token = $1", userSelectColumnsFull, tableUsers),
	queryGetUserByEmail:          fmt.Sprintf("SELECT %s FROM %s WHERE email = $1", userSelectColumnsFull, tableUsers),
	queryGetUserByName:           fmt.Sprintf("SELECT %s FROM %s WHERE username = $1", userSelectColumnsFull, tableUsers),
	queryGetUserByID:             fmt.Sprintf("SELECT %s FROM %s WHERE user_id = $1", userSelectColumnsFull, tableUsers),
	queryInsertUser:              fmt.Sprintf("INSERT INTO %s (username, password, email, status, created_at, updated_at, api_token, api_token_expiry, language, md5 ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", tableUsers),
	queryUpdateUser:              fmt.Sprintf("UPDATE %s SET username = $2, password = $3 , email = $4, status = $5, updated_at = $6 , api_token = $7 , api_token_expiry = $8 , language = $9 , md5 = $10  WHERE user_id = $1", tableUsers),
	queryDeleteUserByID:          fmt.Sprintf("DELETE FROM %s WHERE user_id = $1", tableUsers),
	queryDeleteUserByEmail:       fmt.Sprintf("DELETE FROM %s WHERE email = $1", tableUsers),
	queryDeleteUserByToken:       fmt.Sprintf("DELETE FROM %s WHERE api_token = $1", tableUsers),
	queryUserFollowsUpsert:       fmt.Sprintf("INSERT INTO %s (user_id, following) VALUES($1, $2) ON CONFLICT DO NOTHING", tableUserFollows),
	queryDeleteUserFollowing:     fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND following = $2", tableUserFollows),
}
