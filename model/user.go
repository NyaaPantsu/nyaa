package model

import (
	"time"
)

// omit is the bool type for omitting a field of struct.
type omit bool

type User struct {
	Id              uint      `gorm:"column:user_id;primary_key"`
	Username        string    `gorm:"column:username"`
	Password        string    `gorm:"column:password"`
	Email           string    `gorm:"column:email"`
	Status          int       `gorm:"column:status"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
	/*Api*/Token    string    `gorm:"column:api_token"`
	//ApiTokenExpiry
	TokenExpiration time.Time `gorm:"column:api_token_expiry"`
	Language        string    `gorm:"column:language"`
}

type PublicUser struct {
	User *User
}

type UserFollows struct {
	User      User `gorm:"ForeignKey:user_id"`
	Following User `gorm:"ForeignKey:following"`
}


