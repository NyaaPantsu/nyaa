package models

import (
	"fmt"
	"path"
	"testing"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/azhao12345/gorm"
)

// run before config/parse.go:init()
var _ = func() (_ struct{}) {
	config.ConfigPath = path.Join("..", config.ConfigPath)
	config.DefaultConfigPath = path.Join("..", config.DefaultConfigPath)
	config.Parse()
	return
}()

type errorLogger struct {
	t *testing.T
}

func (logger *errorLogger) Print(values ...interface{}) {
	if len(values) > 1 {
		message := gorm.LogFormatter(values...)
		level := values[0]
		if level == "log" {
			logger.t.Error(message...)
		}

		fmt.Println(message...)
	}
}

func TestGormInitSqlite(t *testing.T) {

	config.Conf.DBType = SqliteType
	config.Conf.DBParams = ":memory:?cache=shared&mode=memory"
	config.Conf.DBLogMode = "detailed"

	db, err := GormInit(config.Conf, &errorLogger{t})
	if err != nil {
		t.Errorf("failed to initialize database: %v", err)
		return
	}

	if db == nil {
		return
	}

	err = db.Close()
	if err != nil {
		t.Errorf("failed to close database: %v", err)
	}
}

// This test requires a running postgres instance. To run it in CI build add these settings in the .travis.yml
// services:
// - postgresql
// before_script:
// - psql -c "CREATE DATABASE nyaapantsu;" -U postgres
// - psql -c "CREATE USER nyaapantsu WITH PASSWORD 'nyaapantsu';" -U postgres
//
// Then enable the test by setting this variable to "true" via ldflags:
// go test ./... -v -ldflags="-X github.com/NyaaPantsu/nyaa/db.testPostgres=true"
var testPostgres = "false"

func TestGormInitPostgres(t *testing.T) {
	if testPostgres != "true" {
		t.Skip("skip", testPostgres)
	}

	config.Conf.DBType = "postgres"
	config.Conf.DBParams = "host=localhost user=nyaapantsu dbname=nyaapantsu sslmode=disable password=nyaapantsu"
	config.Conf.DBLogMode = "detailed"
	config.Conf.Models.CommentsTableName = "comments"
	config.Conf.Models.FilesTableName = "files"
	config.Conf.Models.NotificationsTableName = "notifications"
	config.Conf.Models.ReportsTableName = "torrent_reports"
	config.Conf.Models.TorrentsTableName = "torrents"
	config.Conf.Models.UploadsOldTableName = "user_uploads_old"
	config.Conf.Models.LastOldTorrentID = 90000

	db, err := GormInit(config.Conf, &errorLogger{t})
	if err != nil {
		t.Errorf("failed to initialize database: %v", err)
	}

	if db == nil {
		return
	}

	err = db.Close()
	if err != nil {
		t.Errorf("failed to close database: %v", err)
	}
}
