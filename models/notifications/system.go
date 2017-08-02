package notifications

import (
	"time"

	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/utils/cache"
)

const identifierAnouncement = "system.anouncement"

// NotifyAll notify all users through an anouncement
func NotifyAll(msg string, expire time.Time) (*models.Notification, error) {
	anouncement := &models.Notification{
		Content:    msg,
		Expire:     expire,
		Identifier: identifierAnouncement,
		Read:       false,
		UserID:     0,
	}
	err := models.ORM.Create(anouncement).Error
	return anouncement, err
}

// CheckAnouncement check if there are any new anouncements
func CheckAnouncement() ([]models.Notification, error) {
	if retrieved, ok := cache.C.Get(identifierAnouncement); ok {
		return retrieved.([]models.Notification), nil
	}
	var anouncements []models.Notification
	err := models.ORM.Where("identifier = ? AND expire >= ?", identifierAnouncement, time.Now()).Find(&anouncements).Error
	if err == nil {
		cache.C.Set(identifierAnouncement, anouncements, time.Minute*5)
	}
	return anouncements, err
}
