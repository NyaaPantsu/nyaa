package models

import (
	"time"

	"github.com/NyaaPantsu/nyaa/config"
)

// Scrape model
type Scrape struct {
	TorrentID  uint      `gorm:"column:torrent_id;primary_key"`
	Seeders    uint32    `gorm:"column:seeders"`
	Leechers   uint32    `gorm:"column:leechers"`
	Completed  uint32    `gorm:"column:completed"`
	LastScrape time.Time `gorm:"column:last_scrape"`
}

// TableName : return the table name of the scrape table
func (t Scrape) TableName() string {
	return config.Get().Models.ScrapeTableName
}
