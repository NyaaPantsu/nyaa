package config

const (
	// TorrentFileStorage = "/var/www/wherever/you/want"
	// TorrentStorageLink = "https://your.site/somewhere/%s.torrent"
	TorrentFileStorage = ""
	TorrentStorageLink = ""

	// TODO: deprecate this and move all files to the same server
	TorrentCacheLink = "http://anicache.com/torrent/%s.torrent"
	UploadsDisabled  = true
)
