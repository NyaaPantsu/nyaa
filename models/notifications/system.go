package notifications

import (
	"errors"
	"time"

	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/utils/cache"
	"github.com/NyaaPantsu/nyaa/utils/validator/announcement"
)

const identifierAnnouncement = "system.announcement"

// NotifyAll notify all users through an announcement
func NotifyAll(msg string, expire time.Time) (*models.Notification, error) {
	announcement := &models.Notification{
		Content:    msg,
		Expire:     expire,
		Identifier: identifierAnnouncement,
		Read:       false,
		UserID:     0,
	}
	err := models.ORM.Create(announcement).Error
	return announcement, err
}

// UpdateAnnouncement updates an announcement
func UpdateAnnouncement(announcement *models.Notification, form *announcementValidator.CreateForm) error {
	announcement.Content = form.Message
	if form.Delay > 0 {
		announcement.Expire = time.Now().AddDate(0, 0, form.Delay)
	}
	if models.ORM.Model(announcement).UpdateColumn(announcement).Error != nil {
		return errors.New("Announcement was not updated")
	}
	return nil
}

// CheckAnnouncement check if there are any new announcements
func CheckAnnouncement() ([]models.Notification, error) {
	if retrieved, ok := cache.C.Get(identifierAnnouncement); ok {
		return retrieved.([]models.Notification), nil
	}
	var announcements []models.Notification
	err := models.ORM.Where("identifier = ? AND expire >= ?", identifierAnnouncement, time.Now().Format("2006-01-02")).Find(&announcements).Error
	if err == nil {
		cache.C.Set(identifierAnnouncement, announcements, time.Minute*5)
	}
	return announcements, err
}

// FindAll return all the announcements
func FindAll(limit int, offset int, conditions string, values ...interface{}) ([]models.Notification, int) {
	var announcements []models.Notification
	var nbAnnouncement int
	if conditions == "" {
		conditions += "identifier = ?"
	} else {
		conditions += "AND identifier = ?"
	}
	values = append(values, identifierAnnouncement)
	models.ORM.Model(&announcements).Where(conditions, values...).Count(&nbAnnouncement)
	models.ORM.Limit(limit).Offset(offset).Where(conditions, values...).Find(&announcements)
	return announcements, nbAnnouncement
}

// FindByID return the notification by its ID
func FindByID(id uint) (*models.Notification, error) {
	d := &models.Notification{}
	err := models.ORM.Where("id = ?", id).Find(d).Error
	if err != nil {
		return d, err
	}
	return d, nil
}
