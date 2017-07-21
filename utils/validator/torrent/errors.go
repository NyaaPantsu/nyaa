package torrentValidator

import "errors"

var errTorrentNameInvalid = errors.New("torrent_name_invalid")
var errTorrentDescInvalid = errors.New("torrent_desc_invalid")
var errTorrentMagnetInvalid = errors.New("torrent_magnet_invalid")
var errTorrentURIInvalid = errors.New("torrent_uri_invalid")
var errTorrentCatInvalid = errors.New("torrent_cat_invalid")
var errTorrentLangInvalid = errors.New("torrent_lang_invalid")
var errTorrentPrivate = errors.New("torrent_private")
var errTorrentNoTrackers = errors.New("torrent_no_working_trackers")
var errTorrentAndMagnet = errors.New("torrent_plus_magnet")
var errTorrentHashInvalid = errors.New("torrent_hash_invalid")
