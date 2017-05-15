package common

import (
	"errors"
)

var ErrInsufficientPermission = errors.New("you are not allowed to do that")
var ErrUserExists = errors.New("user already exists")
var ErrNoSuchUser = errors.New("no such user")
var ErrBadLogin = errors.New("bad login")
var ErrUserBanned = errors.New("banned")
var ErrBadEmail = errors.New("bad email")
var ErrNoSuchEntry = errors.New("no such entry")
var ErrNoSuchComment = errors.New("no such comment")
var ErrNotFollowing = errors.New("not following that user")
var ErrInvalidToken = errors.New("invalid token")
var ErrExpiredToken = errors.New("token is expired")
