package model

import (
	"time"
)

type Comment struct {
	Id        uint      `gorm:"column:comment_id;primary_key"`
	TorrentId uint      `gorm:"column:torrent_id"`
	UserId    uint      `gorm:"column:user_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	Torrent   *Torrents `gorm:"ForeignKey:torrent_id"`
	User      *User     `gorm:"ForeignKey:user_id"`
	Content   string    `gorm:"column:content"`
}

// Returns the total size of memory recursively allocated for this struct
func (c Comment) Size() int {
	return (3 + 3*2 + 2 + 2 + len(c.Content)) * 8
}

type OldComment struct {
	TorrentId uint      `gorm:"column:torrent_id"`
	Date      time.Time `gorm:"column:date"`
	Torrent   *Torrents `gorm:"ForeignKey:torrent_id"`
	Username  string    `gorm:"column:username"`
	Content   string    `gorm:"column:content"`
}

// Returns the total size of memory recursively allocated for this struct
func (c OldComment) Size() int {
	return (4 + 2*2 + len(c.Username) + len(c.Content)) * 8
}

func (c OldComment) TableName() string {
	// cba to rename this in the db
	return "comments_old"
}
