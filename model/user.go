package model

import (
	"time"
)

type User struct {
	Id              uint      `gorm:"column:user_id;primary_key"`
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
	LikingCount     int       `json:"likingCount" gorm:"-"`
	LikedCount      int       `json:"likedCount" gorm:"-"`
	Likings         []User    `gorm:"foreignkey:userId;associationforeignkey:follower_id;many2many:users_followers;"`
	Liked           []User    `gorm:"foreignkey:follower_id;associationforeignkey:userId;many2many:users_followers;"`

	Md5             string     `json:"md5"`
	Torrents        []Torrents `gorm:"ForeignKey:UploaderId"`
}

type PublicUser struct {
	User      *User
}

// UsersFollowers is a relation table to relate users each other.
type UsersFollowers struct {
	UserID     uint `gorm:"column:user_id"`
	FollowerID uint `gorm:"column:following"`
}

func (c UsersFollowers) TableName() string {
	return "user_follows"
}
