package common

import (
	"time"
)

type ReportParam struct {
	Limit     uint32
	Offset    uint32
	AllTime   bool
	ID        uint32
	TorrentID uint32
	Before    time.Time
	After     time.Time
}
