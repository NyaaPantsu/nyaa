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
	config.Configpaths[1] = path.Join("..", config.Configpaths[1])
	config.Configpaths[0] = path.Join("..", config.Configpaths[0])
	config.Reload()
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

	config.Get().DBType = SqliteType
	config.Get().DBParams = ":memory:?cache=shared&mode=memory"
	config.Get().DBLogMode = "detailed"

	db, err := GormInit(config.Get(), &errorLogger{t})
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

	config.Get().DBType = "postgres"
	config.Get().DBParams = "host=localhost user=nyaapantsu dbname=nyaapantsu sslmode=disable password=nyaapantsu"
	config.Get().DBLogMode = "detailed"
	config.Get().Models.CommentsTableName = "comments"
	config.Get().Models.FilesTableName = "files"
	config.Get().Models.NotificationsTableName = "notifications"
	config.Get().Models.ReportsTableName = "torrent_reports"
	config.Get().Models.TorrentsTableName = "torrents"
	config.Get().Models.UploadsOldTableName = "user_uploads_old"
	config.Get().Models.LastOldTorrentID = 90000

	db, err := GormInit(config.Get(), &errorLogger{t})
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
