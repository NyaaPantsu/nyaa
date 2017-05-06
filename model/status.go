package model

//user status e.g. verified, filtered, etc
type Status struct {
	Id   int    `json:"id"`
	Name string `json:"name",sql:"size:255"`
}
