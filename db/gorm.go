package db

import (
	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/util/log"
	"github.com/ewhal/nyaa/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	// _ "github.com/go-sql-driver/mysql"
)

var ORM, Errs = GormInit()

// GormInit init gorm ORM.
func GormInit() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", config.DbName)
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
		// db.DropTable(&model.User{}, "UserFollower")
		db.AutoMigrate(&model.Torrents{}, &model.Categories{}, &model.Sub_Categories{}, &model.Statuses{})
		// db.AutoMigrate(&model.User{}, &model.Role{}, &model.Connection{}, &model.Language{}, &model.Article{}, &model.Location{}, &model.Comment{}, &model.File{})
		// db.Model(&model.User{}).AddIndex("idx_user_token", "token")

	}
	log.CheckError(err)

	// relation := gorm.Relationship{}
	// relation.Kind = "many2many"
	// relation.ForeignFieldNames = []string{"id"}            //(M1 pkey)
	// relation.ForeignDBNames = []string{"user_id"}          //(M1 fkey in m1m2join)
	// relation.AssociationForeignFieldNames = []string{"id"} //(M2 pkey)
	// // relation.AssociationForeignStructFieldNames = []string{"id", "ID"} //(m2 pkey name in m2 struct?)
	// relation.AssociationForeignDBNames = []string{"follower_id"} //(m2 fkey in m1m2join)
	// m1Type := reflect.TypeOf(model.User{})
	// m2Type := reflect.TypeOf(model.User{})
	// handler := gorm.JoinTableHandler{}
	// // ORDER BELOW MATTERS
	// // Install handler
	// db.SetJoinTableHandler(&model.User{}, "Likings", &handler)
	// // Configure handler to use the relation that we've defined
	// handler.Setup(&relation, "users_followers", m1Type, m2Type)

	return db, err
}
