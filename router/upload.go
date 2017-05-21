package router

import (
	"encoding/base32"
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

	"github.com/NyaaPantsu/nyaa/cache"
	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/service/upload"
	"github.com/NyaaPantsu/nyaa/util"
	"github.com/NyaaPantsu/nyaa/util/metainfo"
	"github.com/microcosm-cc/bluemonday"
	"github.com/zeebo/bencode"
)

// Use this, because we seem to avoid using models, and we would need
// the torrent ID to create the File in the DB
type UploadedFile struct {
	Path     []string
	Filesize int64
}

// UploadForm serializing HTTP form for torrent upload
type UploadForm struct {
	Name        string
	Magnet      string
	Category    string
	Remake      bool
	Description string
	Status      int
	CaptchaID   string
	WebsiteLink string

	Infohash      string
	CategoryID    int
	SubCategoryID int
	Filesize      int64
	Filepath      string
	FileList      []UploadedFile
}

// TODO: these should be in another package (?)

// form names
const UploadFormName = "name"
const UploadFormTorrent = "torrent"
const UploadFormMagnet = "magnet"
const UploadFormCategory = "c"
const UploadFormRemake = "remake"
const UploadFormDescription = "desc"
const UploadFormWebsiteLink = "website_link"
const UploadFormStatus = "status"

// error indicating that you can't send both a magnet link and torrent
var ErrTorrentPlusMagnet = errors.New("Upload either a torrent file or magnet link, not both")

// error indicating a torrent is private
var ErrPrivateTorrent = errors.New("Torrent is private")

// error indicating a problem with its trackers
var ErrTrackerProblem = errors.New("Torrent does not have any (working) trackers: https://" + config.WebAddress + "/faq#trackers")

// error indicating a torrent's name is invalid
var ErrInvalidTorrentName = errors.New("Torrent name is invalid")

// error indicating a torrent's description is invalid
var ErrInvalidTorrentDescription = errors.New("Torrent description is invalid")

// error indicating a torrent's description is invalid
var ErrInvalidWebsiteLink = errors.New("Website url or IRC link is invalid")

// error indicating a torrent's category is invalid
var ErrInvalidTorrentCategory = errors.New("Torrent category is invalid")

var p = bluemonday.UGCPolicy()

/**
UploadForm.ExtractInfo takes an http request and computes all fields for this form
*/
func (f *UploadForm) ExtractInfo(r *http.Request) error {

	f.Name = r.FormValue(UploadFormName)
	f.Category = r.FormValue(UploadFormCategory)
	f.Description = r.FormValue(UploadFormDescription)
	f.WebsiteLink = r.FormValue(UploadFormWebsiteLink)
	f.Status, _ = strconv.Atoi(r.FormValue(UploadFormStatus))
	f.Magnet = r.FormValue(UploadFormMagnet)
	f.Remake = r.FormValue(UploadFormRemake) == "on"

	// trim whitespace
	f.Name = util.TrimWhitespaces(f.Name)
	f.Description = p.Sanitize(util.TrimWhitespaces(f.Description))
	f.WebsiteLink = util.TrimWhitespaces(f.WebsiteLink)
	f.Magnet = util.TrimWhitespaces(f.Magnet)
	cache.Impl.ClearAll()

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
	if f.WebsiteLink != "" {
		// WebsiteLink
		urlRegexp, _ := regexp.Compile(`^(https?:\/\/|ircs?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$`)
		if !urlRegexp.MatchString(f.WebsiteLink) {
			return ErrInvalidWebsiteLink
		}
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

		// extract filelist
		fileInfos := torrent.Info.GetFiles()
		for _, fileInfo := range fileInfos {
			f.FileList = append(f.FileList, UploadedFile{
				Path:     fileInfo.Path,
				Filesize: int64(fileInfo.Length),
			})
		}
	} else {
		// No torrent file provided
		magnetUrl, err := url.Parse(string(f.Magnet)) //?
		if err != nil {
			return err
		}
		xt := magnetUrl.Query().Get("xt")
		if !strings.HasPrefix(xt, "urn:btih:") {
			return errors.New("Incorrect magnet")
		}
		xt = strings.SplitAfter(xt, ":")[2]
		f.Infohash = strings.ToUpper(strings.Split(xt, "&")[0])
		isBase32, err := regexp.MatchString("^[2-7A-Z]{32}$", f.Infohash)
		if err != nil {
			return err
		}
		if !isBase32 {
			isBase16, err := regexp.MatchString("^[0-9A-F]{40}$", f.Infohash)
			if err != nil {
				return err
			}
			if !isBase16 {
				return errors.New("Incorrect hash")
			}
		} else {
			//convert to base16
			data, err := base32.StdEncoding.DecodeString(f.Infohash)
			if err != nil {
				return err
			}
			hash16 := make([]byte, hex.EncodedLen(len(data)))
			hex.Encode(hash16, data)
			f.Infohash = strings.ToUpper(string(hash16))
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
	f.WebsiteLink = r.FormValue(UploadFormWebsiteLink)
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
