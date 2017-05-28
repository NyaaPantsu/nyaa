package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/NyaaPantsu/nyaa/config"
)

const (
	// UserStatusBanned : Int for User status banned
	UserStatusBanned = -1
	// UserStatusMember : Int for User status member
	UserStatusMember = 0
	// UserStatusTrusted : Int for User status trusted
	UserStatusTrusted = 1
	// UserStatusModerator : Int for User status moderator
	UserStatusModerator = 2
)

// User model
type User struct {
	ID             uint      `gorm:"column:user_id;primary_key"`
	Username       string    `gorm:"column:username"`
	Password       string    `gorm:"column:password"`
	Email          string    `gorm:"column:email"`
	Status         int       `gorm:"column:status"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
	APIToken       string    `gorm:"column:api_token"`
	APITokenExpiry time.Time `gorm:"column:api_token_expiry"`
	Language       string    `gorm:"column:language"`
	Theme          string    `gorm:"column:theme"`
	UserSettings   string    `gorm:"column:settings"`

	// TODO: move this to PublicUser
	Followers []User // Don't work `gorm:"foreignkey:user_id;associationforeignkey:follower_id;many2many:user_follows"`
	Likings   []User // Don't work `gorm:"foreignkey:follower_id;associationforeignkey:user_id;many2many:user_follows"`

	MD5           string         `json:"md5" gorm:"column:md5"` // Hash of email address, used for Gravatar
	Torrents      []Torrent      `gorm:"ForeignKey:UploaderID"`
	Notifications []Notification `gorm:"ForeignKey:UserID"`

	UnreadNotifications int          `gorm:"-"` // We don't want to loop every notifications when accessing user unread notif
	Settings            UserSettings `gorm:"-"` // We don't want to load settings everytime, stock it as a string, parse it when needed
}

// UserJSON : User model conversion in JSON
type UserJSON struct {
	ID          uint   `json:"user_id"`
	Username    string `json:"username"`
	Status      int    `json:"status"`
	CreatedAt   string `json:"created_at"`
	LikingCount int    `json:"liking_count"`
	LikedCount  int    `json:"liked_count"`
}

// Size : Returns the total size of memory recursively allocated for this struct
func (u User) Size() (s int) {
	s += 4 + // ints
		6*2 + // string pointers
		4*3 + //time.Time
		3*2 + // arrays
		// string arrays
		len(u.Username) + len(u.Password) + len(u.Email) + len(u.APIToken) + len(u.MD5) + len(u.Language) + len(u.Theme)
	s *= 8

	// Ignoring foreign key users. Fuck them.

	return
}

// IsBanned : Return true if user is banned
func (u *User) IsBanned() bool {
	return u.Status == UserStatusBanned
}

// IsMember : Return true if user is member
func (u *User) IsMember() bool {
	return u.Status == UserStatusMember
}

// IsTrusted : Return true if user is tusted
func (u *User) IsTrusted() bool {
	return u.Status == UserStatusTrusted
}

// IsModerator : Return true if user is moderator
func (u *User) IsModerator() bool {
	return u.Status == UserStatusModerator
}

// GetUnreadNotifications : Get unread notifications from a user
func (u *User) GetUnreadNotifications() int {
	if u.UnreadNotifications == 0 {
		for _, notif := range u.Notifications {
			if !notif.Read {
				u.UnreadNotifications++
			}
		}
	}
	return u.UnreadNotifications
}

// PublicUser : Is it Deprecated?
type PublicUser struct {
	User *User
}

// UserFollows association table : different users following eachother
type UserFollows struct {
	UserID     uint `gorm:"column:user_id"`
	FollowerID uint `gorm:"column:following"`
}

// UserUploadsOld model : Is it deprecated?
type UserUploadsOld struct {
	Username  string `gorm:"column:username"`
	TorrentID uint   `gorm:"column:torrent_id"`
}

// UserSettings : Struct for user settings, not a model
type UserSettings struct {
	Settings map[string]bool `json:"settings"`
}

// TableName : Return the name of OldComment table
func (c UserUploadsOld) TableName() string {
	// is this needed here?
	return config.UploadsOldTableName
}

// ToJSON : Conversion of a user model to json
func (u *User) ToJSON() UserJSON {
	json := UserJSON{
		ID:          u.ID,
		Username:    u.Username,
		Status:      u.Status,
		CreatedAt:   u.CreatedAt.Format(time.RFC3339),
		LikingCount: len(u.Followers),
		LikedCount:  len(u.Likings),
	}
	return json
}

/* User Settings */

// Get a user setting by keyname
func (s *UserSettings) Get(key string) bool {
	if val, ok := s.Settings[key]; ok {
		return val
	}
	return config.DefaultUserSettings[key]
}

// GetSettings : get all user settings
func (s *UserSettings) GetSettings() map[string]bool {
	return s.Settings
}

// Set a user setting by keyname
func (s *UserSettings) Set(key string, val bool) {
	if s.Settings == nil {
		s.Settings = make(map[string]bool)
	}
	s.Settings[key] = val
}

// ToDefault : Set user settings to default
func (s *UserSettings) ToDefault() {
	s.Settings = config.DefaultUserSettings
}

func (s *UserSettings) initialize() {
	s.Settings = make(map[string]bool)
}

// SaveSettings : Format settings into a json string for preparing before user insertion
func (u *User) SaveSettings() {
	byteArray, err := json.Marshal(u.Settings)

	if err != nil {
		fmt.Print(err)
	}
	u.UserSettings = string(byteArray)
}

// ParseSettings : Function to parse json string into usersettings struct, only parse if necessary
func (u *User) ParseSettings() {
	if len(u.Settings.GetSettings()) == 0 && u.UserSettings != "" {
		u.Settings.initialize()
		json.Unmarshal([]byte(u.UserSettings), &u.Settings)
	} else if len(u.Settings.GetSettings()) == 0 && u.UserSettings != "" {
		u.Settings.initialize()
		u.Settings.ToDefault()
	}
}
