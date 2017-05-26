package apiService

import (
	"encoding/base32"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/service"
	"github.com/NyaaPantsu/nyaa/service/upload"
	"github.com/NyaaPantsu/nyaa/util/metainfo"
	"github.com/zeebo/bencode"
)

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

// TorrentRequest struct
//accept torrent files?
type TorrentRequest struct {
	Name        string `json:"name"`
	Category    int    `json:"category"`
	SubCategory int    `json:"sub_category"`
	Magnet      string `json:"magnet"`
	Hash        string `json:"hash"`
	Description string `json:"description"`
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

func validateName(r *TorrentRequest) (error, int) {
	/*if len(r.Name) < 100 { //isn't this too much?
		return ErrShortName, http.StatusNotAcceptable
	}*/
	return nil, http.StatusOK
}

// TODO Check category is within accepted range
func validateCategory(r *TorrentRequest) (error, int) {
	if r.Category == 0 {
		return ErrCategory, http.StatusNotAcceptable
	}
	return nil, http.StatusOK
}

// TODO Check subCategory is within accepted range
func validateSubCategory(r *TorrentRequest) (error, int) {
	if r.SubCategory == 0 {
		return ErrSubCategory, http.StatusNotAcceptable
	}
	return nil, http.StatusOK
}

func validateMagnet(r *TorrentRequest) (error, int) {
	magnetURL, err := url.Parse(string(r.Magnet)) //?
	if err != nil {
		return err, http.StatusInternalServerError
	}
	xt := magnetURL.Query().Get("xt")
	if !strings.HasPrefix(xt, "urn:btih:") {
		return ErrMagnet, http.StatusNotAcceptable
	}
	xt = strings.SplitAfter(xt, ":")[2]
	r.Hash = strings.ToUpper(strings.Split(xt, "&")[0])
	fmt.Println(r.Hash)
	return nil, http.StatusOK
}

func validateHash(r *TorrentRequest) (error, int) {
	r.Hash = strings.ToUpper(r.Hash)
	isBase32, err := regexp.MatchString("^[2-7A-Z]{32}$", r.Hash)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	if !isBase32 {
		isBase16, err := regexp.MatchString("^[0-9A-F]{40}$", r.Hash)
		if err != nil {
			return err, http.StatusInternalServerError
		}
		if !isBase16 {
			return ErrHash, http.StatusNotAcceptable
		}
	} else {
		//convert to base16
		data, err := base32.StdEncoding.DecodeString(r.Hash)
		if err != nil {
			return err, http.StatusInternalServerError
		}
		hash16 := make([]byte, hex.EncodedLen(len(data)))
		hex.Encode(hash16, data)
		r.Hash = strings.ToUpper(string(hash16))
	}
	return nil, http.StatusOK
}

// ValidateUpload : Check if an upload is valid
func (r *TorrentRequest) ValidateUpload() (err error, code int) {
	validators := []func(r *TorrentRequest) (error, int){
		validateName,
		validateCategory,
		validateSubCategory,
		validateMagnet,
		validateHash,
	}

	for i, validator := range validators {
		if r.Hash != "" && i == 3 {
			continue
		}
		err, code = validator(r)
		if err != nil {
			break
		}
	}
	return err, code
}

// ValidateMultipartUpload : Check if multipart upload is valid
func (r *TorrentRequest) ValidateMultipartUpload(req *http.Request) (int64, error, int) {
	tfile, _, err := req.FormFile("torrent")
	if err == nil {
		var torrent metainfo.TorrentFile

		// decode torrent
		if _, err = tfile.Seek(0, io.SeekStart); err != nil {
			return 0, err, http.StatusInternalServerError
		}
		if err = bencode.NewDecoder(tfile).Decode(&torrent); err != nil {
			return 0, err, http.StatusInternalServerError
		}
		// check a few things
		if torrent.IsPrivate() {
			return 0, errors.New("private torrents not allowed"), http.StatusNotAcceptable
		}
		trackers := torrent.GetAllAnnounceURLS()
		if !uploadService.CheckTrackers(trackers) {
			return 0, errors.New("tracker(s) not allowed"), http.StatusNotAcceptable
		}
		if r.Name == "" {
			r.Name = torrent.TorrentName()
		}

		binInfohash, err := torrent.Infohash()
		if err != nil {
			return 0, err, http.StatusInternalServerError
		}
		r.Hash = strings.ToUpper(hex.EncodeToString(binInfohash[:]))

		// extract filesize
		filesize := int64(torrent.TotalSize())
		err, code := r.ValidateUpload()
		return filesize, err, code
	}
	return 0, err, http.StatusInternalServerError
}

// ValidateUpdate : Check if an update is valid
func (r *TorrentRequest) ValidateUpdate() (err error, code int) {
	validators := []func(r *TorrentRequest) (error, int){
		validateName,
		validateCategory,
		validateSubCategory,
		validateMagnet,
		validateHash,
	}

	//don't update not requested values
	//rewrite with reflect?
	for i, validator := range validators {
		if (r.Name == "" && i == 0) || (r.Category == 0 && i == 1) ||
			(r.SubCategory == 0 && i == 2) ||
			(r.Hash != "" || r.Magnet == "" && i == 3) || (r.Hash == "" && i == 4) {
			continue
		}
		err, code = validator(r)
		if err != nil {
			break
		}
	}

	return err, code
}

// UpdateTorrent : Update torrent model
//rewrite with reflect ?
func (r *UpdateRequest) UpdateTorrent(t *model.Torrent) {
	if r.Update.Name != "" {
		t.Name = r.Update.Name
	}
	if r.Update.Hash != "" {
		t.Hash = r.Update.Hash
	}
	if r.Update.Category != 0 {
		t.Category = r.Update.Category
	}
	if r.Update.SubCategory != 0 {
		t.SubCategory = r.Update.SubCategory
	}
	if r.Update.Description != "" {
		t.Description = r.Update.Description
	}
}
