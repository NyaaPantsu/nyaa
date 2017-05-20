package notifierService

import (
	"github.com/NyaaPantsu/nyaa/db"
	"github.com/NyaaPantsu/nyaa/model"
)


func NotifyUser(user *model.User, name string, msg string) {
	if (user.ID > 0) {
		user.Notifications = append(user.Notifications, model.NewNotification(name, msg))
		// TODO: Email notification
	}
}

func ToggleReadNotification(identifier string, id uint) { // 
	db.ORM.Model(&model.Notification{}).Where("identifier = ? AND user_id = ?", identifier, id).Updates(model.Notification{Read: true})
}