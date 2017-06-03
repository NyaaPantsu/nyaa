package router

import (
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
	"github.com/NyaaPantsu/nyaa/util/categories"
	"github.com/NyaaPantsu/nyaa/util/metainfo"
	"github.com/zeebo/bencode"
)

// Use this, because we seem to avoid using models, and we would need
// the torrent ID to create the File in the DB
type uploadedFile struct {
	Path     []string
	Filesize int64
}

// uploadForm serializing HTTP form for torrent upload
type uploadForm struct {
	Name        string
	Magnet      string
	Category    string
	Remake      bool
	Description string
	Status      int
	Hidden      bool
	CaptchaID   string
	WebsiteLink string

	Infohash      string
	CategoryID    int
	SubCategoryID int
	Filesize      int64
	Filepath      string
	FileList      []uploadedFile
	Trackers      []string
}

// TODO: these should be in another package (?)

// form names
const uploadFormName = "name"
const uploadFormTorrent = "torrent"
const uploadFormMagnet = "magnet"
const uploadFormCategory = "c"
const uploadFormRemake = "remake"
const uploadFormDescription = "desc"
const uploadFormWebsiteLink = "website_link"
const uploadFormStatus = "status"
const uploadFormHidden = "hidden"

// error indicating that you can't send both a magnet link and torrent
var errTorrentPlusMagnet = errors.New("Upload either a torrent file or magnet link, not both")

// error indicating a torrent is private
var errPrivateTorrent = errors.New("Torrent is private")

// error indicating a problem with its trackers
var errTrackerProblem = errors.New("Torrent does not have any (working) trackers: " + config.Conf.WebAddress.Nyaa + "/faq#trackers")

// error indicating a torrent's name is invalid
var errInvalidTorrentName = errors.New("Torrent name is invalid")

// error indicating a torrent's description is invalid
var errInvalidTorrentDescription = errors.New("Torrent description is invalid")

// error indicating a torrent's website link is invalid
var errInvalidWebsiteLink = errors.New("Website url or IRC link is invalid")

// error indicating a torrent's category is invalid
var errInvalidTorrentCategory = errors.New("Torrent category is invalid")

// var p = bluemonday.UGCPolicy()

/**
uploadForm.ExtractInfo takes an http request and computes all fields for this form
*/
func (f *uploadForm) ExtractInfo(r *http.Request) error {

	f.Name = r.FormValue(uploadFormName)
	f.Category = r.FormValue(uploadFormCategory)
	f.Description = r.FormValue(uploadFormDescription)
	f.WebsiteLink = r.FormValue(uploadFormWebsiteLink)
	f.Status, _ = strconv.Atoi(r.FormValue(uploadFormStatus))
	f.Magnet = r.FormValue(uploadFormMagnet)
	f.Remake = r.FormValue(uploadFormRemake) == "on"
	f.Hidden = r.FormValue(uploadFormHidden) == "on"

	// trim whitespace
	f.Name = strings.TrimSpace(f.Name)
	f.Description = util.Sanitize(strings.TrimSpace(f.Description), "default")
	f.WebsiteLink = strings.TrimSpace(f.WebsiteLink)
	f.Magnet = strings.TrimSpace(f.Magnet)
	cache.Impl.ClearAll()
	defer r.Body.Close()

	catsSplit := strings.Split(f.Category, "_")
	// need this to prevent out of index panics
	if len(catsSplit) == 2 {
		CatID, err := strconv.Atoi(catsSplit[0])
		if err != nil {
			return errInvalidTorrentCategory
		}
		SubCatID, err := strconv.Atoi(catsSplit[1])
		if err != nil {
			return errInvalidTorrentCategory
		}

		if !categories.CategoryExists(f.Category) {
			return errInvalidTorrentCategory
		}

		f.CategoryID = CatID
		f.SubCategoryID = SubCatID
	} else {
		return errInvalidTorrentCategory
	}
	if f.WebsiteLink != "" {
		// WebsiteLink
		urlRegexp, _ := regexp.Compile(`^(https?:\/\/|ircs?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$`)
		if !urlRegexp.MatchString(f.WebsiteLink) {
			return errInvalidWebsiteLink
		}
	}

	// first: parse torrent file (if any) to fill missing information
	tfile, _, err := r.FormFile(uploadFormTorrent)
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
			return errPrivateTorrent
		}
		trackers := torrent.GetAllAnnounceURLS()
		f.Trackers = uploadService.CheckTrackers(trackers)
		if len(f.Trackers) == 0 {
			return errTrackerProblem
		}

		// Name
		if len(f.Name) == 0 {
			f.Name = torrent.TorrentName()
		}

		// Magnet link: if a file is provided it should be empty
		if len(f.Magnet) != 0 {
			return errTorrentPlusMagnet
		}

		_, seekErr = tfile.Seek(0, io.SeekStart)
		if seekErr != nil {
			return seekErr
		}
		infohash, err := metainfo.DecodeInfohash(tfile)
		if err != nil {
			return metainfo.ErrInvalidTorrentFile
		}
		f.Infohash = infohash
		f.Magnet = util.InfoHashToMagnet(infohash, f.Name, trackers...)

		// extract filesize
		f.Filesize = int64(torrent.TotalSize())

		// extract filelist
		fileInfos := torrent.Info.GetFiles()
		for _, fileInfo := range fileInfos {
			f.FileList = append(f.FileList, uploadedFile{
				Path:     fileInfo.Path,
				Filesize: int64(fileInfo.Length),
			})
		}
	} else {
		// No torrent file provided
		magnetURL, err := url.Parse(string(f.Magnet)) //?
		if err != nil {
			return err
		}
		xt := magnetURL.Query().Get("xt")
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
		}
		// TODO: Get Trackers from magnet URL
		f.Filesize = 0
		f.Filepath = ""

		return nil
	}

	// then actually check that we have everything we need
	if len(f.Name) == 0 {
		return errInvalidTorrentName
	}

	// after data has been checked & extracted, write it to disk
	if len(config.Conf.Torrents.FileStorage) > 0 {
		err := writeTorrentToDisk(tfile, f.Infohash+".torrent", &f.Filepath)
		if err != nil {
			return err
		}
	} else {
		f.Filepath = ""
	}

	return nil
}

func (f *uploadForm) ExtractEditInfo(r *http.Request) error {
	f.Name = r.FormValue(uploadFormName)
	f.Category = r.FormValue(uploadFormCategory)
	f.WebsiteLink = r.FormValue(uploadFormWebsiteLink)
	f.Description = r.FormValue(uploadFormDescription)
	f.Hidden = r.FormValue(uploadFormHidden) == "on"
	f.Status, _ = strconv.Atoi(r.FormValue(uploadFormStatus))

	// trim whitespace
	f.Name = strings.TrimSpace(f.Name)
	f.Description = util.Sanitize(strings.TrimSpace(f.Description), "default")
	defer r.Body.Close()

	catsSplit := strings.Split(f.Category, "_")
	// need this to prevent out of index panics
	if len(catsSplit) == 2 {
		CatID, err := strconv.Atoi(catsSplit[0])
		if err != nil {
			return errInvalidTorrentCategory
		}
		SubCatID, err := strconv.Atoi(catsSplit[1])
		if err != nil {
			return errInvalidTorrentCategory
		}

		if !categories.CategoryExists(f.Category) {
			return errInvalidTorrentCategory
		}

		f.CategoryID = CatID
		f.SubCategoryID = SubCatID
	} else {
		return errInvalidTorrentCategory
	}
	return nil
}

func writeTorrentToDisk(file multipart.File, name string, fullpath *string) error {
	_, seekErr := file.Seek(0, io.SeekStart)
	if seekErr != nil {
		return seekErr
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	*fullpath = fmt.Sprintf("%s%c%s", config.Conf.Torrents.FileStorage, os.PathSeparator, name)
	return ioutil.WriteFile(*fullpath, b, 0644)
}

// newUploadForm creates a new upload form given parameters as list
func newUploadForm(params ...string) (uploadForm uploadForm) {
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
