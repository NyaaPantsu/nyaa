package model

import (
	"github.com/NyaaPantsu/nyaa/config"
)

type Notification struct {
	ID uint
	Content string
	Read bool
	Identifier string
	Url string
	UserID uint
//	User *User `gorm:"AssociationForeignKey:UserID;ForeignKey:user_id"` // Don't think that we need it here
}

func NewNotification(identifier string, c string, url string) Notification {
	return Notification{Identifier: identifier, Content: c, Url: url}
}

func (n *Notification) TableName() string {
	return config.NotificationTableName
}

