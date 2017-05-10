package model

import (
	"time"
)

type Comment struct {
	ID        uint      `gorm:"column:comment_id;primary_key"`
	TorrentID uint      `gorm:"column:torrent_id"`
	UserID    uint      `gorm:"column:user_id"`
	Content   string    `gorm:"column:content"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	DeletedAt *time.Time

	Torrent *Torrent `gorm:"ForeignKey:torrent_id"`
	User    *User    `gorm:"ForeignKey:user_id"`
}

type OldComment struct {
	TorrentID uint      `gorm:"column:torrent_id"`
	Username  string    `gorm:"column:username"`
	Content   string    `gorm:"column:content"`
	Date      time.Time `gorm:"column:date"`

	Torrent *Torrent `gorm:"ForeignKey:torrent_id"`
}

func (c OldComment) TableName() string {
	// cba to rename this in the db
	// TODO: Update database schema to fix this hack
	//       I find this odd considering how often the schema changes already
	return "comments_old"
}
