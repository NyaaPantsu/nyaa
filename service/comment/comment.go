package commentService

import (
	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"

)

func GetAllComments(limit int, offset int) []model.Comment{
	var comments []model.Comment
	db.ORM.Limit(limit).Offset(offset).Preload("Uploader").Find(&comments)
	return comments
}
