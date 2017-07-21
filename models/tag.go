package models

// Tag model for a torrent vote system
type Tag struct {
	TorrentID uint    `gorm:"column:torrent_id"`
	UserID    uint    `gorm:"column:user_id"`
	Tag       string  `gorm:"column:tag"`
	Type      string  `gorm:"column:type"`
	Weight    float64 `gorm:"column:weight"`
	Accepted  bool    `gorm:"column:accepted"`
}
