package upload

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"reflect"
	"strings"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/utils/sanitize"
	"github.com/NyaaPantsu/nyaa/utils/search/structs"
	"github.com/NyaaPantsu/nyaa/utils/validator/torrent"
	"github.com/gin-gonic/gin"
)

// form names
const uploadFormTorrent = "torrent"

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

// APIResultJSON for torrents in json for api
type APIResultJSON struct {
	Torrents         []models.TorrentJSON `json:"torrents"`
	QueryRecordCount int                  `json:"queryRecordCount"`
	TotalRecordCount int                  `json:"totalRecordCount"`
}

// ToParams : Convert a torrentsrequest to searchparams
func (r *TorrentsRequest) ToParams() structs.WhereParams {
	res := structs.WhereParams{}
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

// ExtractEditInfo : takes an http request and computes all fields for this form
func ExtractEditInfo(c *gin.Context, r *torrentValidator.TorrentRequest) error {
	err := ExtractBasicValue(c, r)
	if err != nil {
		return err
	}

	err = r.ValidateName()
	if err != nil {
		return err
	}

	err = r.ExtractCategory()
	if err != nil {
		return err
	}

	err = r.ExtractLanguage()
	return err
}

// ExtractBasicValue : takes an http request and computes all basic fields for this form
func ExtractBasicValue(c *gin.Context, r *torrentValidator.TorrentRequest) error {
	c.Bind(r)
	// trim whitespace
	r.Name = strings.TrimSpace(r.Name)
	r.Description = sanitize.Sanitize(strings.TrimSpace(r.Description), "default")
	r.WebsiteLink = strings.TrimSpace(r.WebsiteLink)
	r.Magnet = strings.TrimSpace(r.Magnet)

	if len(r.Languages) == 0 { // Shouldn't have to do that since c.Bind actually bind arrays, but better off adding it in case gin doesn't do his work
		r.Languages = c.PostFormArray("languages")
	}
	// then actually check that we have everything we need

	err := r.ValidateDescription()
	if err != nil {
		return err
	}

	err = r.ValidateWebsiteLink()
	return err
}

// ExtractInfo : takes an http request and computes all fields for this form
func ExtractInfo(c *gin.Context, r *torrentValidator.TorrentRequest) error {
	err := ExtractBasicValue(c, r)
	if err != nil {
		return err
	}

	err = r.ExtractCategory()
	if err != nil {
		return err
	}

	err = r.ExtractLanguage()
	if err != nil {
		return err
	}

	tfile, err := r.ValidateMultipartUpload(c, uploadFormTorrent)
	if err != nil {
		return err
	}

	// We check name only here, reason: we can try to retrieve them from the torrent file
	err = r.ValidateName()
	if err != nil {
		return err
	}

	// after data has been checked & extracted, write it to disk
	if len(config.Get().Torrents.FileStorage) > 0 {
		err := writeTorrentToDisk(tfile, r.Infohash+".torrent", &r.Filepath)
		if err != nil {
			return err
		}
	} else {
		r.Filepath = ""
	}

	return nil
}

// UpdateTorrent : Update torrent model
//rewrite with reflect ?
func UpdateTorrent(r *torrentValidator.UpdateRequest, t *models.Torrent, currentUser *models.User) *models.Torrent {
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
	t.Languages = r.Update.Languages
	status := models.TorrentStatusNormal
	if r.Update.Remake { // overrides trusted
		status = models.TorrentStatusRemake
	} else if currentUser.IsTrusted() {
		status = models.TorrentStatusTrusted
	}
	t.Status = status

	t.Hidden = r.Update.Hidden

	return t
}

// UpdateUnscopeTorrent : Update a torrent model without scoping
func UpdateUnscopeTorrent(r *torrentValidator.UpdateRequest, t *models.Torrent, currentUser *models.User) *models.Torrent {
	t = UpdateTorrent(r, t, currentUser)
	t.Status = r.Update.Status
	return t
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
	*fullpath = fmt.Sprintf("%s%c%s", config.Get().Torrents.FileStorage, os.PathSeparator, name)
	return ioutil.WriteFile(*fullpath, b, 0644)
}

// NewTorrentRequest : creates a new torrent request struc with some default value
func NewTorrentRequest(params ...string) *torrentValidator.TorrentRequest {
	torrentRequest := &torrentValidator.TorrentRequest{}
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
	return torrentRequest
}
