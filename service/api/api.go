package apiService

import (
	"encoding/base32"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/NyaaPantsu/nyaa/cache"
	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/service"
	"github.com/NyaaPantsu/nyaa/service/upload"
	"github.com/NyaaPantsu/nyaa/util"
	"github.com/NyaaPantsu/nyaa/util/categories"
	"github.com/NyaaPantsu/nyaa/util/metainfo"
	"github.com/NyaaPantsu/nyaa/util/torrentLanguages"
	"github.com/zeebo/bencode"
)

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
const uploadFormLanguage = "language"

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

// error indicating a torrent's language is invalid
var errInvalidTorrentLanguage = errors.New("Torrent language is invalid")

// error indicating that a non-english torrent was uploaded to a english category
var errNonEnglishLanguageInEnglishCategory = errors.New("Torrent's category is for English translations, but torrent language isn't English.")

// error indicating that a english torrent was uploaded to a non-english category
var errEnglishLanguageInNonEnglishCategory = errors.New("Torrent's category is for non-English translations, but torrent language is English.")

type torrentsQuery struct {
	Category    int `json:"category"`
	SubCategory int `json:"sub_category"`
	Status      int `json:"status"`
	Uploader    int `json:"uploader"`
	Downloads   int `json:"downloads"`
}

// TorrentsRequest struct
type TorrentsRequest struct {
	Query      torrentsQuery `json:"search"`
	Page       int           `json:"page"`
	MaxPerPage int           `json:"limit"`
}

// Use this, because we seem to avoid using models, and we would need
// the torrent ID to create the File in the DB
type uploadedFile struct {
	Path     []string `json:"path"`
	Filesize int64    `json:"filesize"`
}

// TorrentRequest struct
// Same json name as the constant!
type TorrentRequest struct {
	Name        string `json:"name,omitempty"`
	Magnet      string `json:"magnet,omitempty"`
	Category    string `json:"c"`
	Remake      bool   `json:"remake,omitempty"`
	Description string `json:"desc,omitempty"`
	Status      int    `json:"status,omitempty"`
	Hidden      bool   `json:"hidden,omitempty"`
	CaptchaID   string `json:"-"`
	WebsiteLink string `json:"website_link,omitempty"`
	SubCategory int    `json:"sub_category,omitempty"`
	Language    string `json:"language,omitempty"`

	Infohash      string         `json:"hash,omitempty"`
	CategoryID    int            `json:"-"`
	SubCategoryID int            `json:"-"`
	Filesize      int64          `json:"filesize,omitempty"`
	Filepath      string         `json:"-"`
	FileList      []uploadedFile `json:"filelist,omitempty"`
	Trackers      []string       `json:"trackers,omitempty"`
}

// UpdateRequest struct
type UpdateRequest struct {
	ID     int            `json:"id"`
	Update TorrentRequest `json:"update"`
}

// ToParams : Convert a torrentsrequest to searchparams
func (r *TorrentsRequest) ToParams() serviceBase.WhereParams {
	res := serviceBase.WhereParams{}
	conditions := ""
	v := reflect.ValueOf(r.Query)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Interface() != reflect.Zero(field.Type()).Interface() {
			if i != 0 {
				conditions += " AND "
			}
			conditions += v.Type().Field(i).Tag.Get("json") + " = ?"
			res.Params = append(res.Params, field.Interface())
		}
	}
	res.Conditions = conditions
	return res
}

func (r *TorrentRequest) validateName() error {
	// then actually check that we have everything we need
	if len(r.Name) == 0 {
		return errInvalidTorrentName
	}
	return nil
}

func (r *TorrentRequest) validateDescription() error {
	if len(r.Description) > 500 {
		return errInvalidTorrentDescription
	}
	return nil
}

func (r *TorrentRequest) validateWebsiteLink() error {
	if r.WebsiteLink != "" {
		// WebsiteLink
		urlRegexp, _ := regexp.Compile(`^(https?:\/\/|ircs?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$`)
		if !urlRegexp.MatchString(r.WebsiteLink) {
			return errInvalidWebsiteLink
		}
	}
	return nil
}

func (r *TorrentRequest) validateMagnet() error {
	magnetURL, err := url.Parse(string(r.Magnet)) //?
	if err != nil {
		return err
	}
	xt := magnetURL.Query().Get("xt")
	if !strings.HasPrefix(xt, "urn:btih:") {
		return ErrMagnet
	}
	xt = strings.SplitAfter(xt, ":")[2]
	r.Infohash = strings.ToUpper(strings.Split(xt, "&")[0])

	return nil
}

func (r *TorrentRequest) validateHash() error {
	isBase32, err := regexp.MatchString("^[2-7A-Z]{32}$", r.Infohash)
	if err != nil {
		return err
	}
	if !isBase32 {
		isBase16, err := regexp.MatchString("^[0-9A-F]{40}$", r.Infohash)
		if err != nil {
			return err
		}
		if !isBase16 {
			return ErrHash
		}
	} else {
		//convert to base16
		data, err := base32.StdEncoding.DecodeString(r.Infohash)
		if err != nil {
			return err
		}
		hash16 := make([]byte, hex.EncodedLen(len(data)))
		hex.Encode(hash16, data)
		r.Infohash = strings.ToUpper(string(hash16))
	}
	return nil
}

// ExtractEditInfo : takes an http request and computes all fields for this form
func (r *TorrentRequest) ExtractEditInfo(req *http.Request) error {
	err := r.ExtractBasicValue(req)
	if err != nil {
		return err
	}

	err = r.validateName()
	if err != nil {
		return err
	}
	defer req.Body.Close()

	err = r.ExtractCategory(req)
	if err != nil {
		return err
	}

	err = r.ExtractLanguage(req)
	return err
}

// ExtractCategory : takes an http request and computes category field for this form
func (r *TorrentRequest) ExtractCategory(req *http.Request) error {
	catsSplit := strings.Split(r.Category, "_")
	// need this to prevent out of index panics
	if len(catsSplit) != 2 {
		return errInvalidTorrentCategory
	}
	CatID, err := strconv.Atoi(catsSplit[0])
	if err != nil {
		return errInvalidTorrentCategory
	}
	SubCatID, err := strconv.Atoi(catsSplit[1])
	if err != nil {
		return errInvalidTorrentCategory
	}

	if !categories.CategoryExists(r.Category) {
		return errInvalidTorrentCategory
	}

	r.CategoryID = CatID
	r.SubCategoryID = SubCatID
	return nil
}

// ExtractLanguage : takes a http request, computes the torrent language from the form.
func (r *TorrentRequest) ExtractLanguage(req *http.Request) error {
	isEnglishCategory := false
	for _, cat := range config.Conf.Torrents.EnglishOnlyCategories {
		if cat == r.Category {
			isEnglishCategory = true
			break
		}
	}

	if r.Language == "other" || r.Language == "multiple" {
		// In this case, only check if it's on a English-only category.
		if isEnglishCategory {
			return errNonEnglishLanguageInEnglishCategory
		}
		return nil
	}

	if r.Language == "" && isEnglishCategory { // If no language, but in an English category, set to en-us.
		// FIXME Maybe this shouldn't be hard-coded?
		r.Language = "en-us"
	}

	if !torrentLanguages.LanguageExists(r.Language) {
		return errInvalidTorrentLanguage
	}

	if !strings.HasPrefix(r.Language, "en") {
		if isEnglishCategory {
			return errNonEnglishLanguageInEnglishCategory
		}
	} else {
		for _, cat := range config.Conf.Torrents.NonEnglishOnlyCategories {
			if cat == r.Category {
				return errEnglishLanguageInNonEnglishCategory
			}
		}
	}

	return nil
}

// ExtractBasicValue : takes an http request and computes all basic fields for this form
func (r *TorrentRequest) ExtractBasicValue(req *http.Request) error {
	if strings.HasPrefix(req.Header.Get("Content-type"), "multipart/form-data") || req.Header.Get("Content-Type") == "application/x-www-form-urlencoded" { // Multipart
		if strings.HasPrefix(req.Header.Get("Content-type"), "multipart/form-data") { // We parse the multipart form
			err := req.ParseMultipartForm(15485760)
			if err != nil {
				return err
			}
		}
		r.Name = req.FormValue(uploadFormName)
		r.Category = req.FormValue(uploadFormCategory)
		r.WebsiteLink = req.FormValue(uploadFormWebsiteLink)
		r.Description = req.FormValue(uploadFormDescription)
		r.Hidden = req.FormValue(uploadFormHidden) == "on"
		r.Status, _ = strconv.Atoi(req.FormValue(uploadFormStatus))
		r.Remake = req.FormValue(uploadFormRemake) == "on"
		r.Magnet = req.FormValue(uploadFormMagnet)
		r.Language = req.FormValue(uploadFormLanguage)
	} else { // JSON (no file upload then)
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&r)
		if err != nil {
			return err
		}
	}
	// trim whitespace
	r.Name = strings.TrimSpace(r.Name)
	r.Description = util.Sanitize(strings.TrimSpace(r.Description), "default")
	r.WebsiteLink = strings.TrimSpace(r.WebsiteLink)
	r.Magnet = strings.TrimSpace(r.Magnet)

	// then actually check that we have everything we need

	err := r.validateDescription()
	if err != nil {
		return err
	}

	err = r.validateWebsiteLink()
	return err
}

// ExtractInfo : takes an http request and computes all fields for this form
func (r *TorrentRequest) ExtractInfo(req *http.Request) error {
	err := r.ExtractBasicValue(req)
	if err != nil {
		return err
	}

	cache.Impl.ClearAll()
	defer req.Body.Close()

	err = r.ExtractCategory(req)
	if err != nil {
		return err
	}

	err = r.ExtractLanguage(req)
	if err != nil {
		return err
	}

	tfile, err := r.ValidateMultipartUpload(req)
	if err != nil {
		return err
	}

	// We check name only here, reason: we can try to retrieve them from the torrent file
	err = r.validateName()
	if err != nil {
		return err
	}

	// after data has been checked & extracted, write it to disk
	if len(config.Conf.Torrents.FileStorage) > 0 {
		err := writeTorrentToDisk(tfile, r.Infohash+".torrent", &r.Filepath)
		if err != nil {
			return err
		}
	} else {
		r.Filepath = ""
	}

	return nil
}

// ValidateMultipartUpload : Check if multipart upload is valid
func (r *TorrentRequest) ValidateMultipartUpload(req *http.Request) (multipart.File, error) {
	// first: parse torrent file (if any) to fill missing information
	tfile, _, err := req.FormFile(uploadFormTorrent)
	if err == nil {
		var torrent metainfo.TorrentFile

		// decode torrent
		_, seekErr := tfile.Seek(0, io.SeekStart)
		if seekErr != nil {
			return tfile, seekErr
		}
		err = bencode.NewDecoder(tfile).Decode(&torrent)
		if err != nil {
			return tfile, metainfo.ErrInvalidTorrentFile
		}

		// check a few things
		if torrent.IsPrivate() {
			return tfile, errPrivateTorrent
		}
		trackers := torrent.GetAllAnnounceURLS()
		r.Trackers = uploadService.CheckTrackers(trackers)
		if len(r.Trackers) == 0 {
			return tfile, errTrackerProblem
		}

		// Name
		if len(r.Name) == 0 {
			r.Name = torrent.TorrentName()
		}

		// Magnet link: if a file is provided it should be empty
		if len(r.Magnet) != 0 {
			return tfile, errTorrentPlusMagnet
		}

		_, seekErr = tfile.Seek(0, io.SeekStart)
		if seekErr != nil {
			return tfile, seekErr
		}
		infohash, err := metainfo.DecodeInfohash(tfile)
		if err != nil {
			return tfile, metainfo.ErrInvalidTorrentFile
		}
		r.Infohash = infohash
		r.Magnet = util.InfoHashToMagnet(infohash, r.Name, trackers...)

		// extract filesize
		r.Filesize = int64(torrent.TotalSize())

		// extract filelist
		fileInfos := torrent.Info.GetFiles()
		for _, fileInfo := range fileInfos {
			r.FileList = append(r.FileList, uploadedFile{
				Path:     fileInfo.Path,
				Filesize: int64(fileInfo.Length),
			})
		}
	} else {
		err = r.validateMagnet()
		if err != nil {
			return tfile, err
		}
		err = r.validateHash()
		if err != nil {
			return tfile, err
		}
		// TODO: Get Trackers from magnet URL
		r.Filesize = 0
		r.Filepath = ""

		return tfile, nil
	}
	return tfile, err
}

// UpdateTorrent : Update torrent model
//rewrite with reflect ?
func (r *UpdateRequest) UpdateTorrent(t *model.Torrent, currentUser *model.User) {
	if r.Update.Name != "" {
		t.Name = r.Update.Name
	}
	if r.Update.Infohash != "" {
		t.Hash = r.Update.Infohash
	}
	if r.Update.CategoryID != 0 {
		t.Category = r.Update.CategoryID
	}
	if r.Update.SubCategoryID != 0 {
		t.SubCategory = r.Update.SubCategoryID
	}
	if r.Update.Description != "" {
		t.Description = r.Update.Description
	}
	if r.Update.WebsiteLink != "" {
		t.WebsiteLink = r.Update.WebsiteLink
	}
	status := model.TorrentStatusNormal
	if r.Update.Remake { // overrides trusted
		status = model.TorrentStatusRemake
	} else if currentUser.IsTrusted() {
		status = model.TorrentStatusTrusted
	}
	t.Status = status
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

// NewTorrentRequest : creates a new torrent request struc with some default value
func NewTorrentRequest(params ...string) (torrentRequest TorrentRequest) {
	if len(params) > 1 {
		torrentRequest.Category = params[0]
	} else {
		torrentRequest.Category = "3_12"
	}
	if len(params) > 2 {
		torrentRequest.Description = params[1]
	} else {
		torrentRequest.Description = "Description"
	}
	return
}
