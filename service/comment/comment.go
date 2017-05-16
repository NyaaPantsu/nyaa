package commentService

import (
	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
)

func GetAllComments(limit int, offset int, conditions string, values ...interface{}) ([]model.Comment, int) {
	var comments []model.Comment
	var nbComments int
	db.ORM.Table(config.CommentsTableName).Model(&comments).Where(conditions, values...).Count(&nbComments)
	db.ORM.Preload("User").Limit(limit).Offset(offset).Where(conditions, values...).Find(&comments)
	return comments, nbComments
}
