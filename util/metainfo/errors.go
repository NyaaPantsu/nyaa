package metainfo

import (
	"errors"
)

// ErrInvalidTorrentFile : error for indicating we have an invalid torrent file
var ErrInvalidTorrentFile = errors.New("invalid bittorrent file")
