package common

type UserParam struct {
	Full     bool // if true populate Uploads, UsersWeLiked and UsersLikingMe
	Email    string
	Name     string
	ApiToken string
	ID       uint32
	Max      uint32
	Offset   uint32
}
