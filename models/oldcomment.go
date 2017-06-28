package models

import "time"

// OldComment model from old nyaa
type OldComment struct {
	TorrentID uint      `gorm:"column:torrent_id"`
	Username  string    `gorm:"column:username"`
	Content   string    `gorm:"column:content"`
	Date      time.Time `gorm:"column:date"`

	Torrent *Torrent `gorm:"ForeignKey:torrent_id"`
}

// Size : Returns the total size of memory recursively allocated for this struct
func (c OldComment) Size() int {
	return (1 + 2*2 + len(c.Username) + len(c.Content) + 3 + 1) * 8
}

// TableName : Return the name of OldComment table
func (c OldComment) TableName() string {
	// cba to rename this in the db
	// TODO: Update database schema to fix this hack
	//       I find this odd considering how often the schema changes already
	return "comments_old"
}
