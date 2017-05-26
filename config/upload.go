package config

const (
	// TorrentFileStorage = "/var/www/wherever/you/want"
	// TorrentStorageLink = "https://your.site/somewhere/%s.torrent"

	// TorrentFileStorage : Path to default torrent storage location
	TorrentFileStorage = ""
	// TorrentStorageLink : Url of torrent file download location
	TorrentStorageLink = ""

	// TODO: deprecate this and move all files to the same server

	// TorrentCacheLink : Url of torrent site cache
	TorrentCacheLink = "http://anicache.com/torrent/%s.torrent"
	// UploadsDisabled : Disable uploads for everyone except below
	UploadsDisabled = false
	// AdminsAreStillAllowedTo : Enable admin torrent upload even if UploadsDisabled is true
	AdminsAreStillAllowedTo = true
	// TrustedUsersAreStillAllowedTo : Enable trusted users torrent upload even if UploadsDisabled is true
	TrustedUsersAreStillAllowedTo = true
)
