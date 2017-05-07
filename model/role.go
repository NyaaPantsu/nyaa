package model

// Role is a role model for user permission.
type Role struct {
	Id          uint   `json:"id"`
	Name        string `json:"name",sql:"size:255"`
	Description string `json:"description",sql:"size:255"`
}
