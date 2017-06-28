package models

import (
	"html/template"
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

// TemplateTfunc : Used to prevent cyclic import
type TemplateTfunc func(string, ...interface{}) template.HTML

// NewActivity : Create a new activity log
func NewActivity(identifier string, filter string, c ...string) Activity {
	return Activity{Identifier: identifier, Content: strings.Join(c, ","), Filter: filter}
}

// TableName : Return the name of activity table
func (a *Activity) TableName() string {
	return config.Conf.Models.ActivityTableName
}

// ToLocale : Convert list of parameters to message in local language
func (a *Activity) ToLocale(T TemplateTfunc) template.HTML {
	c := strings.Split(a.Content, ",")
	d := make([]interface{}, len(c)-1)
	for i, s := range c[1:] {
		d[i] = s
	}
	return T(c[0], d...)
}
