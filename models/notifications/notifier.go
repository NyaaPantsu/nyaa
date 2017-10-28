package notifications

import (
	"github.com/NyaaPantsu/nyaa/models"
)

// NotifyUser : Notify a user with a notification according to his settings
func NotifyUser(user *models.User, name string, msg string, url string, email bool) {
	if user.ID > 0 {
		notification := models.NewNotification(name, msg, url)
		notification.UserID = user.ID
		models.ORM.Create(&notification)
		// TODO: Email notification
		/*		if email {

				}*/
	}
}

// ToggleReadNotification : Make a notification as read according to its identifier
func ToggleReadNotification(identifier string, id uint) { //
	models.ORM.Model(&models.Notification{}).Where("identifier = ? AND user_id = ?", identifier, id).Updates(models.Notification{Read: true})
}

// DeleteNotifications : Erase notifications from a user
func DeleteNotifications(id uint, all bool) { //
	if all {
		models.ORM.Where("user_id = ?", id).Delete(&models.Notification{})
	} else {
		models.ORM.Where("user_id = ? AND read = ?", id, true).Delete(&models.Notification{})
	}
}
