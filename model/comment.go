package model

import (
	"time"
	
	"github.com/NyaaPantsu/nyaa/config"
)

type Comment struct {
	ID        uint      `gorm:"column:comment_id;primary_key"`
	TorrentID uint      `gorm:"column:torrent_id"`
	UserID    uint      `gorm:"column:user_id"`
	Content   string    `gorm:"column:content"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	DeletedAt *time.Time

	Torrent *Torrent `gorm:"AssociationForeignKey:TorrentID;ForeignKey:torrent_id"`
	User    *User    `gorm:"AssociationForeignKey:UserID;ForeignKey:user_id"`
}

// Returns the total size of memory recursively allocated for this struct
func (c Comment) Size() int {
	return (3 + 3*3 + 2 + 2 + len(c.Content)) * 8
}

func (c Comment) TableName() string {
	return config.CommentsTableName
}

func (c *Comment) Identifier() string { // We Can personalize the identifier but we would have to handle toggle read in that case
	return c.Torrent.Identifier()
}

type OldComment struct {
	TorrentID uint      `gorm:"column:torrent_id"`
	Username  string    `gorm:"column:username"`
	Content   string    `gorm:"column:content"`
	Date      time.Time `gorm:"column:date"`

	Torrent *Torrent `gorm:"ForeignKey:torrent_id"`
}

// Returns the total size of memory recursively allocated for this struct
func (c OldComment) Size() int {
	return (1 + 2*2 + len(c.Username) + len(c.Content) + 3 + 1) * 8
}

func (c OldComment) TableName() string {
	// cba to rename this in the db
	// TODO: Update database schema to fix this hack
	//       I find this odd considering how often the schema changes already
	return "comments_old"
}
