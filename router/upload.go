package router

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/ewhal/nyaa/cache"
	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/service/captcha"
	"github.com/ewhal/nyaa/service/upload"
	"github.com/ewhal/nyaa/util"
	"github.com/ewhal/nyaa/util/metainfo"
	"github.com/microcosm-cc/bluemonday"
	"github.com/zeebo/bencode"
)

// UploadForm serializing HTTP form for torrent upload
type UploadForm struct {
	Name        string
	Magnet      string
	Category    string
	Remake      bool
	Description string
	Status      int
	captcha.Captcha

	Infohash      string
	CategoryID    int
	SubCategoryID int
	Filesize      int64
	Filepath      string
}

// TODO: these should be in another package (?)

// form names
const UploadFormName = "name"
const UploadFormTorrent = "torrent"
const UploadFormMagnet = "magnet"
const UploadFormCategory = "c"
const UploadFormRemake = "remake"
const UploadFormDescription = "desc"
const UploadFormStatus = "status"

// error indicating that you can't send both a magnet link and torrent
var ErrTorrentPlusMagnet = errors.New("upload either a torrent file or magnet link, not both")

// error indicating a torrent is private
var ErrPrivateTorrent = errors.New("torrent is private")

// error indicating a problem with its trackers
var ErrTrackerProblem = errors.New("torrent does not have any (working) trackers: https://" + config.WebAddress + "/faq#trackers")

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
	f.Status, _ = strconv.Atoi(r.FormValue(UploadFormStatus))
	f.Magnet = r.FormValue(UploadFormMagnet)
	f.Remake = r.FormValue(UploadFormRemake) == "on"
	f.Captcha = captcha.Extract(r)

	if !captcha.Authenticate(f.Captcha) {
		// TODO: Prettier passing of mistyped Captcha errors
		return errors.New(captcha.ErrInvalidCaptcha.Error())
	}

	// trim whitespace
	f.Name = util.TrimWhitespaces(f.Name)
	f.Description = p.Sanitize(util.TrimWhitespaces(f.Description))
	f.Magnet = util.TrimWhitespaces(f.Magnet)
	cache.Clear()

	catsSplit := strings.Split(f.Category, "_")
	// need this to prevent out of index panics
	if len(catsSplit) == 2 {
		CatID, err := strconv.Atoi(catsSplit[0])
		if err != nil {
			return ErrInvalidTorrentCategory
		}
		SubCatID, err := strconv.Atoi(catsSplit[1])
		if err != nil {
			return ErrInvalidTorrentCategory
		}

		f.CategoryID = CatID
		f.SubCategoryID = SubCatID
	} else {
		return ErrInvalidTorrentCategory
	}

	// first: parse torrent file (if any) to fill missing information
	tfile, _, err := r.FormFile(UploadFormTorrent)
	if err == nil {
		var torrent metainfo.TorrentFile

		// decode torrent
		_, seekErr := tfile.Seek(0, io.SeekStart)
		if seekErr != nil {
			return seekErr
		}
		err = bencode.NewDecoder(tfile).Decode(&torrent)
		if err != nil {
			return metainfo.ErrInvalidTorrentFile
		}

		// check a few things
		if torrent.IsPrivate() {
			return ErrPrivateTorrent
		}
		trackers := torrent.GetAllAnnounceURLS()
		if !uploadService.CheckTrackers(trackers) {
			return ErrTrackerProblem
		}

		// Name
		if len(f.Name) == 0 {
			f.Name = torrent.TorrentName()
		}

		// Magnet link: if a file is provided it should be empty
		if len(f.Magnet) != 0 {
			return ErrTorrentPlusMagnet
		}
		binInfohash, err := torrent.Infohash()
		if err != nil {
			return err
		}
		f.Infohash = strings.ToUpper(hex.EncodeToString(binInfohash[:]))
		f.Magnet = util.InfoHashToMagnet(f.Infohash, f.Name, trackers...)

		// extract filesize
		f.Filesize = int64(torrent.TotalSize())
	} else {
		// No torrent file provided
		magnetURL, parseErr := url.Parse(f.Magnet)
		if parseErr != nil {
			return metainfo.ErrInvalidTorrentFile
		}
		exactTopic := magnetURL.Query().Get("xt")
		if !strings.HasPrefix(exactTopic, "urn:btih:") {
			return metainfo.ErrInvalidTorrentFile
		}
		exactTopic = strings.SplitAfter(exactTopic, ":")[2]
		f.InfoHash = strings.ToUpper(strings.Split(exactTopic, "&")[0])
		matched, err := regexp.MatchString("^[0-9A-Z]+$", f.Infohash) //ffuuuuuuck
		if err != nil || !matched {
			return metainfo.ErrInvalidTorrentFile
		}

		f.Filesize = 0
		f.Filepath = ""
	}

	// then actually check that we have everything we need
	if len(f.Name) == 0 {
		return ErrInvalidTorrentName
	}

	// after data has been checked & extracted, write it to disk
	if len(config.TorrentFileStorage) > 0 {
		err := WriteTorrentToDisk(tfile, f.Infohash+".torrent", &f.Filepath)
		if err != nil {
			return err
		}
	} else {
		f.Filepath = ""
	}

	return nil
}

func (f *UploadForm) ExtractEditInfo(r *http.Request) error {
	f.Name = r.FormValue(UploadFormName)
	f.Category = r.FormValue(UploadFormCategory)
	f.Description = r.FormValue(UploadFormDescription)
	f.Status, _ = strconv.Atoi(r.FormValue(UploadFormStatus))

	// trim whitespace
	f.Name = util.TrimWhitespaces(f.Name)
	f.Description = p.Sanitize(util.TrimWhitespaces(f.Description))

	catsSplit := strings.Split(f.Category, "_")
	// need this to prevent out of index panics
	if len(catsSplit) == 2 {
		CatID, err := strconv.Atoi(catsSplit[0])
		if err != nil {
			return ErrInvalidTorrentCategory
		}
		SubCatID, err := strconv.Atoi(catsSplit[1])
		if err != nil {
			return ErrInvalidTorrentCategory
		}

		f.CategoryID = CatID
		f.SubCategoryID = SubCatID
	} else {
		return ErrInvalidTorrentCategory
	}
	return nil
}

func WriteTorrentToDisk(file multipart.File, name string, fullpath *string) error {
	_, seekErr := file.Seek(0, io.SeekStart)
	if seekErr != nil {
		return seekErr
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	*fullpath = fmt.Sprintf("%s%c%s", config.TorrentFileStorage, os.PathSeparator, name)
	return ioutil.WriteFile(*fullpath, b, 0644)
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
