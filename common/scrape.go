package common

import (
	"time"
)

type ScrapeResult struct {
	Hash      string
	TorrentID uint32
	Seeders   uint32
	Leechers  uint32
	Completed uint32
	Date      time.Time
}
