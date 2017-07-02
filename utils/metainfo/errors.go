package metainfo

import (
	"errors"
)

// ErrInvalidTorrentFile : error for indicating we have an invalid torrent file
var ErrInvalidTorrentFile = errors.New("torrent_file_invalid")
