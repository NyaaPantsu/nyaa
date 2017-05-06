package router

import (
	"errors"
	"github.com/ewhal/nyaa/util"
	"github.com/ewhal/nyaa/util/metainfo"
	"github.com/zeebo/bencode"
	"net/http"
)

// UploadForm serializing HTTP form for torrent upload
type UploadForm struct {
	Name        string
	Magnet      string
	Category    string
	Description string
}

// TODO: these should be in another package (?)

// form value for torrent name
const UploadFormName = "name"

// form value for torrent file
const UploadFormTorrent = "torrent"

// form value for magnet uri (?)
const UploadFormMagnet = "magnet"

// form value for category
const UploadFormCategory = "c"

// form value for description
const UploadFormDescription = "desc"

// error indicating a torrent is private
var ErrPrivateTorrent = errors.New("torrent is private")

// error indicating a torrent's name is invalid
var ErrInvalidTorrentName = errors.New("torrent name is invalid")

// error indicating a torrent's description is invalid
var ErrInvalidTorrentDescription = errors.New("torrent description is invalid")

/**
UploadForm.ExtractInfo takes an http request and computes all fields for this form
*/
func (f *UploadForm) ExtractInfo(r *http.Request) error {

	f.Name = r.FormValue(UploadFormName)
	f.Category = r.FormValue(UploadFormCategory)
	f.Description = r.FormValue(UploadFormDescription)
	f.Magnet = r.FormValue(UploadFormMagnet)

	// trim whitespaces
	f.Name = util.TrimWhitespaces(f.Name)
	f.Description = util.TrimWhitespaces(f.Description)
	f.Magnet = util.TrimWhitespaces(f.Magnet)

	if len(f.Name) == 0 {
		return ErrInvalidTorrentName
	}

	if len(f.Description) == 0 {
		return ErrInvalidTorrentDescription
	}

	if len(f.Magnet) == 0 {
		// try parsing torrent file if provided if no magnet is specified
		tfile, _, err := r.FormFile(UploadFormTorrent)
		if err != nil {
			return err
		}

		var torrent metainfo.TorrentFile
		// decode torrent
		err = bencode.NewDecoder(tfile).Decode(&torrent)
		if err != nil {
			return metainfo.ErrInvalidTorrentFile
		}

		// check if torrent is private
		if torrent.IsPrivate() {
			return ErrPrivateTorrent
		}

		// generate magnet
		f.Magnet = util.InfoHashToMagnet(torrent.Infohash(), f.Name)
	}
	return nil
}

// NewUploadForm creates a new upload form given parameters as list
func NewUploadForm(params ...string) (uploadForm UploadForm) {
	if len(params) > 1 {
		uploadForm.Category = params[0]
	} else {
		uploadForm.Category = "3_12"
	}
	if len(params) > 2 {
		uploadForm.Description = params[1]
	} else {
		uploadForm.Description = "Description"
	}
	return
}
