package models

import (
	"errors"
	"time"

	"github.com/NyaaPantsu/nyaa/config"
)

// Notification model
type Notification struct {
	ID         uint
	Content    string
	Read       bool
	Identifier string
	URL        string
	Expire     time.Time
	Date       time.Time
	UserID     uint
	//	User *User `gorm:"AssociationForeignKey:UserID;ForeignKey:user_id"` // Don't think that we need it here
}

// NewNotification : Create a new notification
func NewNotification(identifier string, c string, url string) Notification {
	return Notification{Identifier: identifier, Content: c, URL: url, Date: time.Now()}
}

// TableName : Return the name of notification table
func (n *Notification) TableName() string {
	return config.Get().Models.NotificationsTableName
}

// Delete a notification
func (n *Notification) Delete() error {
	if n.ID == 0 {
		return errors.New("Can't delete a non existent notification")
	}
	ORM.Where("id = ?", n.ID).Delete(n)
	return nil
}
