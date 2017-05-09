package db

import (
	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/util/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var ORM *gorm.DB

// GormInit init gorm ORM.
func GormInit(conf *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(conf.DBType, conf.DBParams)

	db.DB().Ping()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	// TODO: Enable Gorm initialization for non-development builds
	if config.Environment == "DEVELOPMENT" {
		db.LogMode(true)
		db.AutoMigrate(&model.Torrent{}, &model.UsersFollowers{}, &model.User{}, &model.Comment{}, &model.OldComment{})
	}
	log.CheckError(err)

	return db, err
}
