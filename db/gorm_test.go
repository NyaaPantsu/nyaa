package db

import (
	"testing"

	"github.com/ewhal/nyaa/config"
)

func TestGormInit(t *testing.T) {
	conf := config.New()
	conf.DBType = "sqlite3"
	conf.DBParams = ":memory:?cache=shared&mode=memory"

	db, err := GormInit(conf)
	if err != nil {
		t.Errorf("failed to initialize database: %v", err)
		return
	}

	err = db.Close()
	if err != nil {
		t.Errorf("failed to close database: %v", err)
	}
}
