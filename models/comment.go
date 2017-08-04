package models

import (
	"html/template"
	"time"

	"github.com/NyaaPantsu/nyaa/config"
)

// Comment model
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

// CommentJSON for comment model in json
type CommentJSON struct {
	Username   string        `json:"username"`
	UserID     int           `json:"user_id"`
	UserAvatar string        `json:"user_avatar"`
	Content    template.HTML `json:"content"`
	Date       time.Time     `json:"date"`
}

// Size : Returns the total size of memory recursively allocated for this struct
func (c Comment) Size() int {
	return (3 + 3*3 + 2 + 2 + len(c.Content)) * 8
}

// TableName : Return the name of comment table
func (c Comment) TableName() string {
	return config.Get().Models.CommentsTableName
}

// Identifier : Return the identifier of the comment
func (c *Comment) Identifier() string { // We Can personalize the identifier but we would have to handle toggle read in that case
	return c.Torrent.Identifier()
}
