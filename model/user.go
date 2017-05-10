package model

import (
	"time"
)

type User struct {
	ID              uint      `gorm:"column:user_id;primary_key"`
	Username        string    `gorm:"column:username"`
	Password        string    `gorm:"column:password"`
	Email           string    `gorm:"column:email"`
	Status          int       `gorm:"column:status"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
	Token           string    `gorm:"column:api_token"`
	TokenExpiration time.Time `gorm:"column:api_token_expiry"`
	Language        string    `gorm:"column:language"`

	// TODO: move this to PublicUser
	LikingCount int    `json:"likingCount" gorm:"-"`
	LikedCount  int    `json:"likedCount" gorm:"-"`
	Likings     []User // Don't work `gorm:"foreignkey:user_id;associationforeignkey:follower_id;many2many:user_follows"`
	Liked       []User // Don't work `gorm:"foreignkey:follower_id;associationforeignkey:user_id;many2many:user_follows"`

	MD5      string    `json:"md5"` // Hash of email address, used for Gravatar
	Torrents []Torrent `gorm:"ForeignKey:UploaderID"`
}

type PublicUser struct {
	User *User
}

// different users following eachother
type UserFollows struct {
	UserID     uint `gorm:"column:user_id"`
	FollowerID uint `gorm:"column:following"`
}
