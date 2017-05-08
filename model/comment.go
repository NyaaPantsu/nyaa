package model

import (
	"time"
)

// Comment is a comment model.
type Comment struct {
	Id        int    `json:"id"`
	Content   string `json:"content"`
	UserId    int    `json:"userId"`
	Username  string `json:"username"` // this is duplicate but it'll be faster rite?
	TorrentId int
	// LikingCount int       `json:"likingCount"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt""`
	User      User       `json:"user"`
}
