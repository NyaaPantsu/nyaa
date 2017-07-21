package activities

import (
	"github.com/NyaaPantsu/nyaa/models"
)

// Log : log an activity from a user to his own id (System user id is 0)
func Log(user *models.User, name string, filter string, msg ...string) {
	activity := models.NewActivity(name, filter, msg...)
	activity.UserID = user.ID
	models.ORM.Create(&activity)
}

// DeleteAll : Erase aticities from a user (System user id is 0)
func DeleteAll(id uint) {
	models.ORM.Where("user_id = ?", id).Delete(&models.Activity{})
}

// FindAll : Get All activities
func FindAll(limit int, offset int, conditions string, values ...interface{}) ([]models.Activity, int) {
	var activities []models.Activity
	var nbActivities int
	models.ORM.Model(&activities).Where(conditions, values...).Count(&nbActivities)
	models.ORM.Preload("User").Limit(limit).Offset(offset).Order("id DESC").Where(conditions, values...).Find(&activities)
	return activities, nbActivities
}
