package db

import (
	"github.com/azhao12345/gorm"
	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/util/log"
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

	if config.Environment == "DEVELOPMENT" {
		db.LogMode(true)
	}

	db.AutoMigrate(&model.User{}, &model.UserFollows{}, &model.UserUploadsOld{})
	db.AutoMigrate(&model.Torrent{}, &model.TorrentReport{})
	db.AutoMigrate(&model.Comment{}, &model.OldComment{})

	return db, nil
}
