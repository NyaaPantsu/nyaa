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
	db.ORM.Limit(limit).Offset(offset).Where(conditions, values...).Preload("User").Find(&comments)
	return comments, nbComments
}

// DeleteComment : Delete a comment
// FIXME : move this to comment service
func DeleteComment(id string) (*model.Comment, int, error) {
	var comment model.Comment
	if db.ORM.Where("comment_id = ?", id).Preload("User").Preload("Torrent").Find(&comment).RecordNotFound() {
		return &comment, http.StatusNotFound, errors.New("Comment is not found")
	}
	if db.ORM.Delete(&comment).Error != nil {
		return &comment, http.StatusInternalServerError, errors.New("Comment is not deleted")
	}
	return &comment, http.StatusOK, nil
}
