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
func ToggleReadNotification(identifier string, user *models.User) { //
	models.ORM.Model(&models.Notification{}).Where("identifier = ? AND user_id = ?", identifier, user.ID).Updates(models.Notification{Read: true})
	for i, notif := range user.Notifications {
		if notif.Identifier == identifier {
			user.Notifications[i].Read = true
		}
	}
	//Need to update both DB and variable, otherwise when the function is called the user still needs to do an additional refresh to see the notification gone/read
}

// MarkAllNotificationsAsRead : Force every notification as read
func MarkAllNotificationsAsRead(user *models.User) { //
	models.ORM.Model(&models.Notification{}).Where("user_id = ?", user.ID).Updates(models.Notification{Read: true})
	for i := range user.Notifications {
			user.Notifications[i].Read = true
	}
	//Need to update both DB and variable, otherwise when the function is called the user still needs to do an additional refresh to see the notification gone/read
}

// DeleteNotifications : Erase notifications from a user
func DeleteNotifications(user *models.User, all bool) { //
	if all {
		models.ORM.Where("user_id = ?", user.ID).Delete(&models.Notification{})
		user.Notifications = []models.Notification{}
	} else {
		models.ORM.Where("user_id = ? AND read = ?", user.ID, true).Delete(&models.Notification{})
		NewNotifications := []models.Notification{}
		for _, notif := range user.Notifications {
			if !notif.Read {
				NewNotifications = append(NewNotifications, notif)
			}
		}
		user.Notifications = NewNotifications
	}
}
