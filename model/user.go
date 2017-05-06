package model

import (
	"time"
)

// omit is the bool type for omitting a field of struct.
type omit bool

// User is a user model
type User struct {
	Id uint `json:"id"`

	Email           string    `json:"email",sql:"size:255;unique"`
	Password        string    `json:"password",sql:"size:255"`
	Username        string    `json:"username",sql:"size:255;unique"`
	Description     string    `json:"description",sql:"size:100"`
	Token           string    `json:"token"`
	TokenExpiration time.Time `json:"tokenExperiation"`

	// email md5 for gravatar
	Md5 string `json:"md5"`

	// admin
	Activation         bool       `json:"activation"`
	PasswordResetToken string     `json:"passwordResetToken"`
	ActivationToken    string     `json:"activationToken"`
	PasswordResetUntil time.Time  `json:"passwordResetUntil"`
	ActivateUntil      time.Time  `json:"activateUntil"`
	ActivatedAt        time.Time  `json:"activatedAt"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
	DeletedAt          *time.Time `json:"deletedAt"`
	LastLoginAt        time.Time  `json:"lastLoginAt"`
	CurrentLoginAt     time.Time  `json:"currentLoginAt"`
	LastLoginIp        string     `json:"lastLoginIp",sql:"size:100"`
	CurrentLoginIp     string     `json:"currentLoginIp",sql:"size:100"`

	// Liking
	LikingCount int `json:"likingCount"`
	LikedCount  int `json:"likedCount"`
	Likings     []User     `gorm:"foreignkey:userId;associationforeignkey:follower_id;many2many:users_followers;"`
	Liked      []User     `gorm:"foreignkey:follower_id;associationforeignkey:userId;many2many:users_followers;"`

	Connections []Connection
	Languages   []Language `gorm:"many2many:user_languages;"` // Many To Many, user_languages is the join table
	Roles       []Role     `gorm:"many2many:users_roles;"`    // Many To Many, users_roles
	Torrents    []Torrents
}

// UsersFollowers is a relation table to relate users each other.
type UsersFollowers struct {
	UserID     uint `json:"user_id"`
	FollowerID uint `json:"follower_id"`
}

// PublicUser is a public user model that contains only a few information for everyone.
type PublicUser struct {
	*User
	Email           omit `json:"email,omitempty",sql:"size:255;unique"`
	Password        omit `json:"password,omitempty",sql:"size:255"`
	Token           omit `json:"token,omitempty"`
	TokenExpiration omit `json:"tokenExperiation,omitempty"`

	// admin
	Activation         omit `json:"activation,omitempty"`
	PasswordResetToken omit `json:"passwordResetToken,omitempty"`
	ActivationToken    omit `json:"activationToken,omitempty"`
	PasswordResetUntil omit `json:"passwordResetUntil,omitempty"`
	ActivateUntil      omit `json:"activateUntil,omitempty"`
	ActivatedAt        omit `json:"activatedAt,omitempty"`
	UpdatedAt          omit `json:"updatedAt,omitempty"`
	DeletedAt          omit `json:"deletedAt,omitempty"`
	LastLoginAt        omit `json:"lastLoginAt,omitempty"`
	CurrentLoginAt     omit `json:"currentLoginAt,omitempty"`
	LastLoginIp        omit `json:"lastLoginIp,omitempty",sql:"size:100"`
	CurrentLoginIp     omit `json:"currentLoginIp,omitempty",sql:"size:100"`

	Connections omit `json:"connections,omitempty"`
	Languages   omit `json:"languages,omitempty"`
	Roles       omit `json:"roles,omitempty"`
	Torrents    omit `json:"articles,omitempty"`
}

// Connection is a connection model for oauth.
type Connection struct {
	Id             uint   `json:"id"`
	UserId         uint   `json:"userId"`
	ProviderId     uint   `gorm:"column:provider_id", json:"providerId"`
	ProviderUserId string `gorm:"column:provider_user_id", json:"providerUserId"`
	AccessToken    string `json:"accessToken"`
	ProfileUrl     string `gorm:"column:profile_url", json:"profileUrl"`
	ImageUrl       string `gorm:"column:image_url", json:"imageUrl"`
}
