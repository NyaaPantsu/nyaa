package commentService

import (
	"net/http"

	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/service/user/permission"
	"github.com/ewhal/nyaa/util/modelHelper"

)

func GetAllComments(limit int, offset int) {
	var comments []model.Comment
	db.ORM.Limit(limit).Offset(offset).Preload("Uploader").Find(&comments)
	return comments
}
