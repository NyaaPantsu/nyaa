package model

import (
	"time"

	"github.com/NyaaPantsu/nyaa/config"
)

const (
	UserStatusBanned    = -1
	UserStatusMember    = 0
	UserStatusTrusted   = 1
	UserStatusModerator = 2
)

type User struct {
	ID             uint      `gorm:"column:user_id;primary_key"`
	Username       string    `gorm:"column:username"`
	Password       string    `gorm:"column:password"`
	Email          string    `gorm:"column:email"`
	Status         int       `gorm:"column:status"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
	ApiToken       string    `gorm:"column:api_token"`
	ApiTokenExpiry time.Time `gorm:"column:api_token_expiry"`
	Language       string    `gorm:"column:language"`
	UserSettings   string    `gorm:"column:settings"`

	// TODO: move this to PublicUser
	Likings     []User // Don't work `gorm:"foreignkey:user_id;associationforeignkey:follower_id;many2many:user_follows"`
	Liked       []User // Don't work `gorm:"foreignkey:follower_id;associationforeignkey:user_id;many2many:user_follows"`

	MD5      string    `json:"md5" gorm:"column:md5"` // Hash of email address, used for Gravatar
	Torrents []Torrent `gorm:"ForeignKey:UploaderID"`
	Notifications []Notification `gorm:"ForeignKey:UserID"`

	UnreadNotifications int `gorm:"-"` // We don't want to loop every notifications when accessing user unread notif
	Settings UserSettings `gorm:"-"` // We don't want to load settings everytime, stock it as a string, parse it when needed
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
		len(u.Username) + len(u.Password) + len(u.Email) + len(u.ApiToken) + len(u.MD5) + len(u.Language)
	s *= 8

	// Ignoring foreign key users. Fuck them.

	return
}

func (u User) IsBanned() bool {
	return u.Status == UserStatusBanned
}
func (u User) IsMember() bool {
	return u.Status == UserStatusMember
}
func (u User) IsTrusted() bool {
	return u.Status == UserStatusTrusted
}
func (u User) IsModerator() bool {
	return u.Status == UserStatusModerator
}

func (u User) GetUnreadNotifications() int {
	if u.UnreadNotifications == 0 { 
		for _, notif := range u.Notifications {
			if !notif.Read {
				u.UnreadNotifications++
			}
		}
	}
	return u.UnreadNotifications
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

type UserSettings struct {
	settings map[string]interface{} `json:"settings"`
}

func (c UserUploadsOld) TableName() string {
	// is this needed here?
	return config.UploadsOldTableName
}

func (u *User) ToJSON() UserJSON {
	json := UserJSON{
		ID:          u.ID,
		Username:    u.Username,
		Status:      u.Status,
		CreatedAt:   u.CreatedAt.Format(time.RFC3339),
		LikingCount: len(u.Likings),
		LikedCount:  len(u.Liked),
	}
	return json
}

/* User Settings */

func(s UserSettings) Get(key string) interface{} {
	if (s.settings[key] != nil) {
	return s.settings[key]
	} else {
		return config.DefaultUserSettings[key]
	}
}

func (s UserSettings) GetSettings() {
	return s.settings
}

func (s UserSettings) Set(key string, val interface{}) {
	s.settings[key] = val
}

func (s UserSettings) ToDefault() {
	s.settings = config.DefaultUserSettings
}

func (u User) SaveSettings() {
	u.UserSettings , _ = json.Marshal(u.Settings.GetSettings())
}

func (u User) ParseSettings() {
	if len(u.Settings.GetSettings()) == 0 && u.UserSettings != "" {
		json.Unmarshal([]byte(u.UserSettings), u.Settings)
	}
}