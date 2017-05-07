package router

import (
	"encoding/hex"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"regexp"

	"github.com/ewhal/nyaa/service/captcha"
	"github.com/ewhal/nyaa/util"
	"github.com/ewhal/nyaa/util/metainfo"
	"github.com/microcosm-cc/bluemonday"
	"github.com/zeebo/bencode"
)

// UploadForm serializing HTTP form for torrent upload
type UploadForm struct {
	Name          string
	Magnet        string
	Infohash      string
	Category      string
	CategoryId    int
	SubCategoryId int
	Description   string
	captcha.Captcha
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

// error indicating a problem with its trackers
// FIXME: hardcoded link
var ErrTrackerProblem = errors.New("torrent does not have any (working) trackers: https://nyaa.pantsu.cat/faq#trackers")

// error indicating a torrent's name is invalid
var ErrInvalidTorrentName = errors.New("torrent name is invalid")

// error indicating a torrent's description is invalid
var ErrInvalidTorrentDescription = errors.New("torrent description is invalid")

// error indicating a torrent's category is invalid
var ErrInvalidTorrentCategory = errors.New("torrent category is invalid")

var p = bluemonday.UGCPolicy()

/**
UploadForm.ExtractInfo takes an http request and computes all fields for this form
*/
func (f *UploadForm) ExtractInfo(r *http.Request) error {

	f.Name = r.FormValue(UploadFormName)
	f.Category = r.FormValue(UploadFormCategory)
	f.Description = r.FormValue(UploadFormDescription)
	f.Magnet = r.FormValue(UploadFormMagnet)
	f.Captcha = captcha.Extract(r)

	// trim whitespaces
	f.Name = util.TrimWhitespaces(f.Name)
	f.Description = p.Sanitize(util.TrimWhitespaces(f.Description))
	f.Magnet = util.TrimWhitespaces(f.Magnet)

	if len(f.Name) == 0 {
		return ErrInvalidTorrentName
	}

	//if len(f.Description) == 0 {
	//	return ErrInvalidTorrentDescription
	//}

	catsSplit := strings.Split(f.Category, "_")
	// need this to prevent out of index panics
	if len(catsSplit) == 2 {
		CatId, err := strconv.Atoi(catsSplit[0])
		if err != nil {
			return ErrInvalidTorrentCategory
		}
		SubCatId, err := strconv.Atoi(catsSplit[1])
		if err != nil {
			return ErrInvalidTorrentCategory
		}

		f.CategoryId = CatId
		f.SubCategoryId = SubCatId
	} else {
		return ErrInvalidTorrentCategory
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

		// check trackers
		trackers := torrent.GetAllAnnounceURLS()
		if !CheckTrackers(trackers) {
			return ErrTrackerProblem
		}

		// generate magnet
		binInfohash := torrent.Infohash()
		f.Infohash = hex.EncodeToString(binInfohash[:])
		f.Magnet = util.InfoHashToMagnet(f.Infohash, f.Name)
		f.Infohash = strings.ToUpper(f.Infohash)
	} else {
		magnetUrl, parseErr := url.Parse(f.Magnet)
		if parseErr != nil {
			return metainfo.ErrInvalidTorrentFile
		}
		exactTopic := magnetUrl.Query().Get("xt")
		if !strings.HasPrefix(exactTopic, "urn:btih:") {
			return metainfo.ErrInvalidTorrentFile
		} else {
			f.Infohash = strings.ToUpper(strings.TrimPrefix(exactTopic, "urn:btih:"))
			matched, err := regexp.MatchString("^[0-9A-F]{40}$", f.Infohash)
			if err != nil || !matched {
				return metainfo.ErrInvalidTorrentFile
			}
		}
	}

	return nil
}

var dead_trackers = []string{ // substring matches!
	"://open.nyaatorrents.info:6544",
	"://tracker.openbittorrent.com:80",
	"://tracker.publicbt.com:80",
	"://stats.anisource.net:2710",
	"://exodus.desync.com",
	"://open.demonii.com:1337",
	"://tracker.istole.it:80",
	"://tracker.ccc.de:80",
	"://bt2.careland.com.cn:6969",
	"://announce.torrentsmd.com:8080"}

func CheckTrackers(trackers []string) bool {
	var numGood int
	for _, t := range trackers {
		var good bool = true
		for _, check := range dead_trackers {
			if strings.Contains(t, check) {
				good = false
			}
		}
		if good {
			numGood += 1
		}
	}
	return numGood > 0
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
