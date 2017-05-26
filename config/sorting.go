package config

// TODO: Perform sorting configuration at runtime
//       Future hosts shouldn't have to rebuild the binary to update a setting

const (
	// TorrentOrder : Default sorting field for torrents
	TorrentOrder = "torrent_id"
	// TorrentSort : Default sorting order for torrents
	TorrentSort = "DESC"
)
