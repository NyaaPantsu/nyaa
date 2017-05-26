package apiService

import "errors"

// ErrShortName : Error for invalid file name used by api
var ErrShortName = errors.New("file name should be at least 100 characters long")

// ErrCategory : Error for not found category used by api
var ErrCategory = errors.New("this category doesn't exist")

// ErrSubCategory : Error for not found sub category used by api
var ErrSubCategory = errors.New("this sub category doesn't exist")

// ErrMagnet : Error for incorrect magnet used by api
var ErrMagnet = errors.New("incorrect magnet")

// ErrHash : Error for incorrect hash used by api
var ErrHash = errors.New("incorrect hash")

// ErrAPIKey : Error for incorrect api key used by api
var ErrAPIKey = errors.New("incorrect api key")

// ErrTorrentID : Error for torrent id used by api
var ErrTorrentID = errors.New("torrent with requested id doesn't exist")

// ErrRights : Error for rights used by api
var ErrRights = errors.New("not enough rights for this request")
