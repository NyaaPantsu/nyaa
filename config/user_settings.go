package config


/*
 * Here we config the notifications options
 * Uses in user model for default setting
 * Be aware, default values in user update form are
 * in service/user/form/form_validator.go
 */

var DefaultUserSettings = map[string]bool {
	"new_torrent": true,
	"new_torrent_email": false,
	"new_comment": true,
	"new_comment_email": false,
	"new_responses": false,
	"new_responses_email": false,
	"new_follower": false,
	"new_follower_email": false,
	"followed": false,
	"followed_email": false,
}