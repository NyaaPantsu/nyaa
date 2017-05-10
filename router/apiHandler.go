package router

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strconv"
	"time"

	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/service"
	"github.com/ewhal/nyaa/service/api"
	"github.com/ewhal/nyaa/service/torrent"
	"github.com/ewhal/nyaa/util"
	"github.com/ewhal/nyaa/util/log"
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
			req.MaxPerPage = 50
		}
		if req.Page == 0 {
			req.Page = 1
		}

		whereParams = req.ToParams()
	} else {
		var err error
		maxString := r.URL.Query().Get("max")
		if maxString != "" {
			req.MaxPerPage, err = strconv.Atoi(maxString)
			if !log.CheckError(err) {
				req.MaxPerPage = 50 // default Value maxPerPage
			}
		}

		req.Page = 1
		if page != "" {
			req.Page, err = strconv.Atoi(html.EscapeString(page))
			if !log.CheckError(err) {
				http.Error(w, err.Error(), http.StatusInternalServerError)
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
	b := torrent.ToJSON()
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(b)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ApiUploadHandler(w http.ResponseWriter, r *http.Request) {
	if config.UploadsDisabled {
		http.Error(w, "Error uploads are disabled", http.StatusInternalServerError)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" {
		token := r.Header.Get("Authorization")
		user := model.User{}
		db.ORM.Where("api_token = ?", token).First(&user) //i don't like this
		if user.ID == 0 {
			http.Error(w, apiService.ErrApiKey.Error(), http.StatusForbidden)
			return
		}

		defer r.Body.Close()

		//verify token
		//token := r.Header.Get("Authorization")

		upload := apiService.TorrentRequest{}
		d := json.NewDecoder(r.Body)
		if err := d.Decode(&upload); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err, code := upload.ValidateUpload()
		if err != nil {
			http.Error(w, err.Error(), code)
			return
		}

		torrent := model.Torrent{
			Name:        upload.Name,
			Category:    upload.Category,
			SubCategory: upload.SubCategory,
			Status:      1,
			Hash:        upload.Hash,
			Date:        time.Now(),
			Filesize:    0, //?
			Description: upload.Description,
			UploaderID:  user.ID,
			Uploader:    &user,
		}

		db.ORM.Create(&torrent)
		if err != nil {
			util.SendError(w, err, 500)
			return
		}
		fmt.Printf("%+v\n", torrent)
	}
}

func ApiUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if config.UploadsDisabled {
		http.Error(w, "Error uploads are disabled", http.StatusInternalServerError)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" {
		token := r.Header.Get("Authorization")
		user := model.User{}
		db.ORM.Where("api_token = ?", token).First(&user) //i don't like this
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
