package notifierService

import (
	"github.com/NyaaPantsu/nyaa/db"
	"github.com/NyaaPantsu/nyaa/model"
)

// NotifyUser : Notify a user with a notification according to his settings
func NotifyUser(user *model.User, name string, msg string, url string, email bool) {
	if user.ID > 0 {
		notification := model.NewNotification(name, msg, url)
		notification.UserID = user.ID
		db.ORM.Create(&notification)
		// TODO: Email notification
		/*		if email {

				}*/
	}
}

// ToggleReadNotification : Make a notification as read according to its identifier
func ToggleReadNotification(identifier string, id uint) { //
	db.ORM.Model(&model.Notification{}).Where("identifier = ? AND user_id = ?", identifier, id).Updates(model.Notification{Read: true})
}

// DeleteAllNotifications : Erase notifications from a user
func DeleteAllNotifications(id uint) { //
	db.ORM.Where("user_id = ?", id).Delete(&model.Notification{})
}
