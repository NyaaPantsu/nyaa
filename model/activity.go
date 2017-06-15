package model

import (
	"strings"

	"github.com/NyaaPantsu/nyaa/config"
)

// Activity model
type Activity struct {
	ID         uint
	Content    string
	Identifier string
	Filter     string
	UserID     uint
	User       *User
}

// NewActivity : Create a new activity log
func NewActivity(identifier string, filter string, c ...string) Activity {
	return Activity{Identifier: identifier, Content: strings.Join(c, ","), Filter: filter}
}

// TableName : Return the name of activity table
func (a *Activity) TableName() string {
	return config.Conf.Models.ActivityTableName
}
