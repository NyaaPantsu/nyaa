package config


/*
 * Here we config the notifications options
 */
var EnableNotifications = map[string]bool {
	"new_torrent": true,
	"new_comment_owner": true,
	"new_comment_all": true,
	"new_follower": false,
	"followed": false,
}