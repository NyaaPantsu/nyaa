package apiService

import (
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/service"
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
	Category    int    `json:"category"`
	SubCategory int    `json:"sub_category"`
	Magnet      string `json:"magnet"`
	Hash        string `json:"hash"`
	Description string `json:"description"`
}

type UpdateRequest struct {
	ID     int            `json:"id"`
	Update TorrentRequest `json:"update"`
}

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
	if len(r.Name) < 100 { //isn't this too much?
		return ErrShortName, http.StatusNotAcceptable
	}
	return nil, http.StatusOK
}

func validateCategory(r *TorrentRequest) (error, int) {
	if r.Category == 0 {
		return ErrCategory, http.StatusNotAcceptable
	}
	return nil, http.StatusOK
}

func validateSubCategory(r *TorrentRequest) (error, int) {
	if r.SubCategory == 0 {
		return ErrSubCategory, http.StatusNotAcceptable
	}
	return nil, http.StatusOK
}

func validateMagnet(r *TorrentRequest) (error, int) {
	magnetUrl, err := url.Parse(string(r.Magnet)) //?
	if err != nil {
		return err, http.StatusInternalServerError
	}
	exactTopic := magnetUrl.Query().Get("xt")
	if !strings.HasPrefix(exactTopic, "urn:btih:") {
		return ErrMagnet, http.StatusNotAcceptable
	}
	r.Hash = strings.ToUpper(strings.TrimPrefix(exactTopic, "urn:btih:"))
	return nil, http.StatusOK
}

func validateHash(r *TorrentRequest) (error, int) {
	r.Hash = strings.ToUpper(r.Hash)
	matched, err := regexp.MatchString("^[0-9A-F]{40}$", r.Hash)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	if !matched {
		return ErrHash, http.StatusNotAcceptable
	}
	return nil, http.StatusOK
}

//rewrite validators!!!

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
