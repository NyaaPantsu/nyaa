package apiService

import (
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/service/torrent"
)

type torrentsQuery struct {
	Category    int `json:"category"`
	SubCategory int `json:"sub_category"`
	Status      int `json:"status"`
	Uploader    int `json:"uploader"`
	Downloads   int `json:"downloads"`
}

type TorrentsRequest struct {
	Query      torrentsQuery `json:"search"`
	Page       int           `json:"page"`
	MaxPerPage int           `json:"limit"`
}

//accept torrent files?
type TorrentRequest struct {
	Name        string `json:"name"`
	Hash        string `json:"hash"`
	Magnet      string `json:"magnet"`
	Category    int    `json:"category"`
	SubCategory int    `json:"sub_category"`
	Description string `json:"description"`
}

type UpdateRequest struct {
	Id     int            `json:"id"`
	Update TorrentRequest `json:"update"`
}

func (r *TorrentsRequest) ToParams() torrentService.WhereParams {
	res := torrentService.WhereParams{}
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

var ErrShortName = errors.New("file name should be at least 100 characters long")
var ErrCategory = errors.New("this category doesn't exist")
var ErrSubCategory = errors.New("this sub category doesn't exist")
var ErrMagnet = errors.New("incorrect magnet")
var ErrHash = errors.New("incorrect hash")

//rewrite validators!!!

func (r *TorrentRequest) ValidateUpload() (error, int) {
	if len(r.Name) < 100 {
		return ErrShortName, http.StatusNotAcceptable
	}
	if r.Category == 0 {
		return ErrCategory, http.StatusNotAcceptable
	}
	if r.SubCategory == 0 {
		return ErrSubCategory, http.StatusNotAcceptable
	}

	if r.Hash == "" {
		magnetUrl, err := url.Parse(string(r.Magnet)) //?
		if err != nil {
			return err, http.StatusInternalServerError
		}
		exactTopic := magnetUrl.Query().Get("xt")
		if !strings.HasPrefix(exactTopic, "urn:btih:") {
			return ErrMagnet, http.StatusNotAcceptable
		}
		r.Hash = strings.ToUpper(strings.TrimPrefix(exactTopic, "urn:btih:"))
	}

	matched, err := regexp.MatchString("^[0-9A-F]{40}$", r.Hash)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	if !matched {
		return ErrHash, http.StatusNotAcceptable
	}

	return nil, http.StatusOK
}

func (r *TorrentRequest) ValidateUpdate() (error, int) {
	if len(r.Name) < 100 && len(r.Name) != 0 {
		return ErrShortName, http.StatusNotAcceptable
	}
	/*if r.Category == 0 {
		return ErrCategory, http.StatusNotAcceptable
	}
	if r.SubCategory == 0 {
		return ErrSubCategory, http.StatusNotAcceptable
	}*/

	if r.Magnet != "" || r.Hash != "" {
		if r.Hash == "" {
			magnetUrl, err := url.Parse(string(r.Magnet)) //?
			if err != nil {
				return err, http.StatusInternalServerError
			}
			exactTopic := magnetUrl.Query().Get("xt")
			if !strings.HasPrefix(exactTopic, "urn:btih:") {
				return ErrMagnet, http.StatusNotAcceptable
			}
			r.Hash = strings.ToUpper(strings.TrimPrefix(exactTopic, "urn:btih:"))
		}

		matched, err := regexp.MatchString("^[0-9A-F]{40}$", r.Hash)
		if err != nil {
			return err, http.StatusInternalServerError
		}
		if !matched {
			return ErrHash, http.StatusNotAcceptable
		}
	}

	return nil, http.StatusOK
}

//rewrite with reflect ?
func (r *UpdateRequest) UpdateTorrent(t *model.Torrents) {
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
		t.Sub_Category = r.Update.SubCategory
	}
	if r.Update.Description != "" {
		t.Description = r.Update.Description
	}
}
