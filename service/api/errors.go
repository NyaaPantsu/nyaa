package apiService

import "errors"

var ErrShortName = errors.New("file name should be at least 100 characters long")
var ErrCategory = errors.New("this category doesn't exist")
var ErrSubCategory = errors.New("this sub category doesn't exist")
var ErrMagnet = errors.New("incorrect magnet")
var ErrHash = errors.New("incorrect hash")
var ErrApiKey = errors.New("incorrect api key")
var ErrTorrentId = errors.New("torrent with requested id doesn't exist")
var ErrRights = errors.New("not enough rights for this request")
