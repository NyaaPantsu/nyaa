package commentService

import (
	"errors"
	"net/http"

	"github.com/NyaaPantsu/nyaa/db"
	"github.com/NyaaPantsu/nyaa/model"
)

// GetAllComments : Get all comments based on conditions
func GetAllComments(limit int, offset int, conditions string, values ...interface{}) ([]model.Comment, int) {
	var comments []model.Comment
	var nbComments int
	db.ORM.Model(&comments).Where(conditions, values...).Count(&nbComments)
	db.ORM.Preload("User").Limit(limit).Offset(offset).Where(conditions, values...).Find(&comments)
	return comments, nbComments
}

// DeleteComment : Delete a comment
// FIXME : move this to comment service
func DeleteComment(id string) (int, error) {
	var comment model.Comment
	if db.ORM.First(&comment, id).RecordNotFound() {
		return http.StatusNotFound, errors.New("Comment is not found")
	}
	if db.ORM.Delete(&comment).Error != nil {
		return http.StatusInternalServerError, errors.New("Comment is not deleted")
	}
	return http.StatusOK, nil
}
