package model

import (
	"time"
)

// omit is the bool type for omitting a field of struct.
type omit bool

type User struct {
	Id        uint      `gorm:"column:user_id;primary_key"`
	Username  string    `gorm:"column:username"`
	Password  string    `gorm:"column:password"`
	Email     string    `gorm:"column:email"`
	Status    int       `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	/*Api*/ Token string `gorm:"column:api_token"`
	//ApiTokenExpiry

	// Liking
	LikingCount int    `json:"likingCount"`
	LikedCount  int    `json:"likedCount"`
	Likings     []User `gorm:"foreignkey:userId;associationforeignkey:follower_id;many2many:users_followers;"`
	Liked       []User `gorm:"foreignkey:follower_id;associationforeignkey:userId;many2many:users_followers;"`

	Md5             string     `json:"md5"`
	TokenExpiration time.Time  `gorm:"column:api_token_expiry"`
	Language        string     `gorm:"column:language"`
	Torrents        []Torrents `gorm:"ForeignKey:UploaderId"`
}

// Returns the total size of memory recursively allocated for this struct
func (u User) Size() (s int) {
	s += 4 + // ints
		6*2 + // string pointers
		4*3 + //time.Time
		3*2 + // arrays
		// string arrays
		len(u.Username) + len(u.Password) + len(u.Email) + len(u.Token) + len(u.Md5) + len(u.Language)
	s *= 8

	// Ignoring foreign key users. Fuck them.

	return
}

type PublicUser struct {
	User *User
}

// UsersFollowers is a relation table to relate users each other.
type UsersFollowers struct {
	UserID     uint `gorm:"column:userId"`
	FollowerID uint `gorm:"column:follower_id"`
}
