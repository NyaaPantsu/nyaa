package comments

import (
	"errors"
	"net/http"

	"github.com/NyaaPantsu/nyaa/models"
)

// FindAll : Find all comments based on conditions
func FindAll(limit int, offset int, conditions string, values ...interface{}) ([]models.Comment, int) {
	var comments []models.Comment
	var nbComments int
	models.ORM.Model(&comments).Where(conditions, values...).Count(&nbComments)
	models.ORM.Limit(limit).Offset(offset).Where(conditions, values...).Preload("User").Find(&comments)
	return comments, nbComments
}

// Delete : Delete a comment
func Delete(id uint) (*models.Comment, int, error) {
	var comment models.Comment
	if models.ORM.Where("comment_id = ?", id).Preload("User").Preload("Torrent").Find(&comment).RecordNotFound() {
		return &comment, http.StatusNotFound, errors.New("Comment is not found")
	}
	if models.ORM.Delete(&comment).Error != nil {
		return &comment, http.StatusInternalServerError, errors.New("Comment is not deleted")
	}
	return &comment, http.StatusOK, nil
}
