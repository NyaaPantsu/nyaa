package activity

import (
	"html/template"
	"strings"

	"github.com/NyaaPantsu/nyaa/db"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/util/publicSettings"
)

// Log : log an activity from a user to his own id (System user id is 0)
func Log(user *model.User, name string, filter string, msg ...string) {
	activity := model.NewActivity(name, filter, msg...)
	activity.UserID = user.ID
	db.ORM.Create(&activity)
}

// DeleteAll : Erase aticities from a user (System user id is 0)
func DeleteAll(id uint) {
	db.ORM.Where("user_id = ?", id).Delete(&model.Activity{})
}

// ToLocale : Convert list of parameters to message in local language
func ToLocale(a *model.Activity, T publicSettings.TemplateTfunc) template.HTML {
	c := strings.Split(a.Content, ",")
	d := make([]interface{}, len(c)-1)
	for i, s := range c[1:] {
		d[i] = s
	}
	return T(c[0], d...)
}

// GetAllActivities : Get All activities
func GetAllActivities(limit int, offset int, conditions string, values ...interface{}) ([]model.Activity, int) {
	var activities []model.Activity
	var nbActivities int
	db.ORM.Model(&activities).Where(conditions, values...).Count(&nbActivities)
	db.ORM.Preload("User").Limit(limit).Offset(offset).Order("id DESC").Where(conditions, values...).Find(&activities)
	return activities, nbActivities
}
