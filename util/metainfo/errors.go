package metainfo

import (
	"errors"
)

// error for indicating we have an invalid torrent file
var ErrInvalidTorrentFile = errors.New("invalid bittorrent file")
