package db

import (
	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/util/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	// _ "github.com/go-sql-driver/mysql"
)

var ORM *gorm.DB

// GormInit init gorm ORM.
func GormInit(conf *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(conf.DBType, conf.DBParams)
	// db, err := gorm.Open("mysql", config.MysqlDSL())
	//db, err := gorm.Open("sqlite3", "/tmp/gorm.db")

	// Get database connection handle [*sql.DB](http://golang.org/pkg/database/sql/#DB)
	db.DB()

	// Then you could invoke `*sql.DB`'s functions with it
	db.DB().Ping()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	// Disable table name's pluralization
	// db.SingularTable(true)
	if config.Environment == "DEVELOPMENT" {
		db.LogMode(true)
		db.AutoMigrate(&model.Torrents{}, &model.UserFollows{}, &model.User{}, &model.Comment{}, &model.OldComment{})
		// db.Model(&model.User{}).AddIndex("idx_user_token", "token")

	}
	log.CheckError(err)

	return db, err
}
