package db

import (
	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/util/log"
	"github.com/azhao12345/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Need for postgres support
	_ "github.com/jinzhu/gorm/dialects/sqlite"   // Need for sqlite
)

// Logger interface
type Logger interface {
	Print(v ...interface{})
}

// DefaultLogger : use the default gorm logger that prints to stdout
var DefaultLogger Logger = nil

// ORM : Variable for interacting with database
var ORM *gorm.DB

// IsSqlite : Variable to know if we are in sqlite or postgres
var IsSqlite bool

// GormInit init gorm ORM.
func GormInit(conf *config.Config, logger Logger) (*gorm.DB, error) {

	db, openErr := gorm.Open(conf.DBType, conf.DBParams)
	if openErr != nil {
		log.CheckError(openErr)
		return nil, openErr
	}

	IsSqlite = conf.DBType == "sqlite"

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

	switch conf.DBLogMode {
	case "detailed":
		db.LogMode(true)
	case "silent":
		db.LogMode(false)
	}

	if logger != nil {
		db.SetLogger(logger)
	}

	db.AutoMigrate(&model.User{}, &model.UserFollows{}, &model.UserUploadsOld{}, &model.Notification{})
	if db.Error != nil {
		return db, db.Error
	}
	db.AutoMigrate(&model.Torrent{}, &model.TorrentReport{})
	if db.Error != nil {
		return db, db.Error
	}
	db.AutoMigrate(&model.File{})
	if db.Error != nil {
		return db, db.Error
	}
	db.AutoMigrate(&model.Comment{}, &model.OldComment{})
	if db.Error != nil {
		return db, db.Error
	}

	return db, nil
}
