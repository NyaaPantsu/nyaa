package db

import (
	"fmt"
	"testing"

	"github.com/azhao12345/gorm"
	"github.com/NyaaPantsu/nyaa/config"
)

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
	conf := config.New()
	conf.DBType = "sqlite3"
	conf.DBParams = ":memory:?cache=shared&mode=memory"
	conf.DBLogMode = "detailed"

	db, err := GormInit(conf, &errorLogger{t})
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

	conf := config.New()
	conf.DBType = "postgres"
	conf.DBParams = "host=localhost user=nyaapantsu dbname=nyaapantsu sslmode=disable password=nyaapantsu"
	conf.DBLogMode = "detailed"

	db, err := GormInit(conf, &errorLogger{t})
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
