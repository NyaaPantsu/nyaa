package database

import (
	"github.com/ewhal/nyaa/common"
	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/database/postgres"
	//"github.com/ewhal/nyaa/db/sqlite"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/util/log"

	"database/sql"
	"errors"
)

// Database obstraction layer
type Database interface {

	// Initialize internal state
	Init() error

	// return true if we need to call MigrateNext again
	NeedsMigrate() (bool, error)
	// migrate to next database revision
	MigrateNext() error

	// get torrents given parameters
	GetTorrentsWhere(param *common.TorrentParam) ([]model.Torrent, error)

	// insert new comment
	InsertComment(comment *model.Comment) error
	// new torrent report
	InsertTorrentReport(report *model.TorrentReport) error

	// check if user A follows B (by id)
	UserFollows(a, b uint32) (bool, error)

	// delete reports given params
	DeleteTorrentReportsWhere(param *common.ReportParam) (uint32, error)

	// get reports given params
	GetTorrentReportsWhere(param *common.ReportParam) ([]model.TorrentReport, error)

	// bulk record scrape events in 1 transaction
	RecordScrapes(scrapes []common.ScrapeResult) error

	// insert new user
	InsertUser(u *model.User) error

	// update existing user info
	UpdateUser(u *model.User) error

	// get users given paramteters
	GetUsersWhere(param *common.UserParam) ([]model.User, error)

	// delete many users given parameters
	DeleteUsersWhere(param *common.UserParam) (uint32, error)

	// get comments by given parameters
	GetCommentsWhere(param *common.CommentParam) ([]model.Comment, error)

	// delete comment by given parameters
	DeleteCommentsWhere(param *common.CommentParam) (uint32, error)

	// add user A following B
	AddUserFollowing(a, b uint32) error

	// delete user A following B
	DeleteUserFollowing(a, b uint32) (bool, error)

	// insert/update torrent
	UpsertTorrent(t *model.Torrent) error

	// delete torrents by given parameters
	// DeleteTorrentsWhere(param *common.TorrentParam) (uint32, error)


	// DO NOT USE ME kthnx
	Query(query string, params ...interface{}) (*sql.Rows, error)
}

var ErrInvalidDatabaseDialect = errors.New("invalid database dialect")
var ErrSqliteSucksAss = errors.New("sqlite3 sucks ass so it's not supported yet")

var Impl Database

func Configure(conf *config.Config) (err error) {
	switch conf.DBType {
	case "postgres":
		Impl, err = postgres.New(conf.DBParams)
		break
	case "sqlite3":
		err = ErrSqliteSucksAss
		// Impl, err = sqlite.New(conf.DBParams)
		break
	default:
		err = ErrInvalidDatabaseDialect
	}
	if err == nil {
		log.Infof("Init %s database", conf.DBType)
		err = Impl.Init()
	}
	return
}

// Migrate migrates the database to latest revision, call after Configure
func Migrate() (err error) {
	next := true
	for err == nil && next {
		next, err = Impl.NeedsMigrate()
		if next {
			err = Impl.MigrateNext()
		}

	}
	return
}
