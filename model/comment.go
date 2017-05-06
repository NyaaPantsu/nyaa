package model

import (
	"time"
)

// Comment is a comment model.
type Comment struct {
	Id          int      `json:"id"`
	Content     string    `json:"content"`
	UserId      int      `json:"userId"`
	TorrentId	int
	// LikingCount int       `json:"likingCount"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   time.Time `json:"deletedAt"`
	User        User      `json:"user"`
}
