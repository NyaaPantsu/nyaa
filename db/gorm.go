package db

import (
	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/util/log"
	"github.com/azhao12345/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var ORM *gorm.DB

// GormInit init gorm ORM.
func GormInit(conf *config.Config) (*gorm.DB, error) {
	db, openErr := gorm.Open(conf.DBType, conf.DBParams)
	if openErr != nil {
		log.CheckError(openErr)
		return nil, openErr
	}

	connectionErr := db.DB().Ping()
	if connectionErr != nil {
		log.CheckError(connectionErr)
		return nil, connectionErr
	}
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	// TODO: Enable Gorm initialization for non-development builds
	if config.Environment == "DEVELOPMENT" {
		db.LogMode(true)
		db.AutoMigrate(&model.Torrent{}, &model.UserFollows{}, &model.User{}, &model.Comment{}, &model.OldComment{}, &model.TorrentReport{})
	}

	return db, nil
}
