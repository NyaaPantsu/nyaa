package router

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/db"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/service"
	"github.com/NyaaPantsu/nyaa/service/api"
	"github.com/NyaaPantsu/nyaa/service/upload"
	"github.com/NyaaPantsu/nyaa/service/torrent"
	"github.com/NyaaPantsu/nyaa/util"
	"github.com/NyaaPantsu/nyaa/util/log"
	"github.com/gorilla/mux"
)

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]
	whereParams := serviceBase.WhereParams{}
	req := apiService.TorrentsRequest{}

	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" {
		d := json.NewDecoder(r.Body)
		if err := d.Decode(&req); err != nil {
			util.SendError(w, err, 502)
		}

		if req.MaxPerPage == 0 {
			req.MaxPerPage = config.TorrentsPerPage
		}
		if req.Page <= 0 {
			req.Page = 1
		}

		whereParams = req.ToParams()
	} else {
		var err error
		maxString := r.URL.Query().Get("max")
		if maxString != "" {
			req.MaxPerPage, err = strconv.Atoi(maxString)
			if !log.CheckError(err) {
				req.MaxPerPage = config.TorrentsPerPage
			}
		} else {
			req.MaxPerPage = config.TorrentsPerPage
		}

		req.Page = 1
		if page != "" {
			req.Page, err = strconv.Atoi(html.EscapeString(page))
			if !log.CheckError(err) {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if req.Page <= 0 {
				NotFoundHandler(w, r)
				return
			}
		}
	}

	torrents, nbTorrents, err := torrentService.GetTorrents(whereParams, req.MaxPerPage, req.MaxPerPage*(req.Page-1))
	if err != nil {
		util.SendError(w, err, 400)
		return
	}

	b := model.ApiResultJSON{
		Torrents: model.TorrentsToJSON(torrents),
	}
	b.QueryRecordCount = req.MaxPerPage
	b.TotalRecordCount = nbTorrents
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ApiViewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	torrent, err := torrentService.GetTorrentById(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	b := torrent.ToJSON()
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(b)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ApiUploadHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	user := model.User{}
	db.ORM.Where("api_token = ?", token).First(&user) //i don't like this

	if !uploadService.IsUploadEnabled(user) {
		http.Error(w, "Error uploads are disabled", http.StatusBadRequest)
		return
	}

	if user.ID == 0 {
		http.Error(w, apiService.ErrApiKey.Error(), http.StatusUnauthorized)
		return
	}

	upload := apiService.TorrentRequest{}
	var filesize int64

	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" {
		defer r.Body.Close()

		d := json.NewDecoder(r.Body)
		if err := d.Decode(&upload); err != nil {
			decodeError := fmt.Errorf("Unable to decode upload data: %s", err).Error()
			http.Error(w, decodeError, http.StatusInternalServerError)
			return
		}
		err, code := upload.ValidateUpload()
		if err != nil {
			http.Error(w, err.Error(), code)
			return
		}
	} else if strings.HasPrefix(contentType, "multipart/form-data") {
		upload.Name = r.FormValue("name")
		upload.Category, _ = strconv.Atoi(r.FormValue("category"))
		upload.SubCategory, _ = strconv.Atoi(r.FormValue("sub_category"))
		upload.Description = r.FormValue("description")

		var err error
		var code int

		filesize, err, code = upload.ValidateMultipartUpload(r)
		if err != nil {
			http.Error(w, err.Error(), code)
			return
		}
	} else {
		// TODO What should we do here ? upload is empty so we shouldn't
		// create a torrent from it
		err := fmt.Errorf("Please provide either of Content-Type: application/json header or multipart/form-data").Error()
		http.Error(w, err, http.StatusInternalServerError)
	}
	var sameTorrents int

	db.ORM.Model(&model.Torrent{}).Where("torrent_hash = ?", upload.Hash).Count(&sameTorrents)

	if sameTorrents == 0 {
		torrent := model.Torrent{
			Name:        upload.Name,
			Category:    upload.Category,
			SubCategory: upload.SubCategory,
			Status:      1,
			Hash:        upload.Hash,
			Date:        time.Now(),
			Filesize:    filesize,
			Description: upload.Description,
			UploaderID:  user.ID,
			Uploader:    &user,
		}
		db.ORM.Create(&torrent)
		/*if err != nil {
			util.SendError(w, err, 500)
			return
		}*/
		fmt.Printf("%+v\n", torrent)
	} else {
		http.Error(w, "torrent already exists", http.StatusBadRequest)
		return
	}
}

// FIXME Impossible to update a torrent uploaded by user 0
func ApiUpdateHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	user := model.User{}
	db.ORM.Where("api_token = ?", token).First(&user) //i don't like this

	if !uploadService.IsUploadEnabled(user) {
		http.Error(w, "Error uploads are disabled", http.StatusBadRequest)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" {
		if user.ID == 0 {
			http.Error(w, apiService.ErrApiKey.Error(), http.StatusForbidden)
			return
		}

		defer r.Body.Close()

		update := apiService.UpdateRequest{}
		d := json.NewDecoder(r.Body)
		if err := d.Decode(&update); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id := update.ID
		torrent := model.Torrent{}
		db.ORM.Where("torrent_id = ?", id).First(&torrent)
		if torrent.ID == 0 {
			http.Error(w, apiService.ErrTorrentId.Error(), http.StatusBadRequest)
			return
		}
		if torrent.UploaderID != 0 && torrent.UploaderID != user.ID { //&& user.Status != mod
			http.Error(w, apiService.ErrRights.Error(), http.StatusForbidden)
			return
		}

		err, code := update.Update.ValidateUpdate()
		if err != nil {
			http.Error(w, err.Error(), code)
			return
		}
		update.UpdateTorrent(&torrent)

		db.ORM.Save(&torrent)
		if err != nil {
			util.SendError(w, err, 500)
			return
		}
		fmt.Printf("%+v\n", torrent)
	}
}
