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

	MD5      string    `json:"md5" gorm:"column:md5"` // Hash of email address, used for Gravatar
	Torrents []Torrent `gorm:"ForeignKey:UploaderID"`
}

type UserJSON struct {
	ID          uint   `json:"user_id"`
	Username    string `json:"username"`
	Status      int    `json:"status"`
	CreatedAt   string `json:"created_at"`
	LikingCount int    `json:"liking_count"`
	LikedCount  int    `json:"liked_count"`
}

// Returns the total size of memory recursively allocated for this struct
func (u User) Size() (s int) {
	s += 4 + // ints
		6*2 + // string pointers
		4*3 + //time.Time
		3*2 + // arrays
		// string arrays
		len(u.Username) + len(u.Password) + len(u.Email) + len(u.Token) + len(u.MD5) + len(u.Language)
	s *= 8

	// Ignoring foreign key users. Fuck them.

	return
}

type PublicUser struct {
	User *User
}

// different users following eachother
type UserFollows struct {
	UserID     uint `gorm:"column:user_id"`
	FollowerID uint `gorm:"column:following"`
}

type UserUploadsOld struct {
	Username  string `gorm:"column:username"`
	TorrentId uint   `gorm:"column:torrent_id"`
}

func (c UserUploadsOld) TableName() string {
	// TODO: rename this in db
	return "user_uploads_old"
}

func (u *User) ToJSON() UserJSON {
	json := UserJSON{
		ID:          u.ID,
		Username:    u.Username,
		Status:      u.Status,
		CreatedAt:   u.CreatedAt.Format(time.RFC3339),
		LikingCount: u.LikingCount,
		LikedCount:  u.LikedCount,
	}
	return json
}
