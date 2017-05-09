package router

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	//"sort"
	"strconv"
	"time"

	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/service/api"
	"github.com/ewhal/nyaa/service/torrent"
	"github.com/ewhal/nyaa/util"
	"github.com/gorilla/mux"
)

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]
	whereParams := torrentService.WhereParams{}
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
		var errConv error
		req.MaxPerPage, errConv = strconv.Atoi(r.URL.Query().Get("max"))
		if errConv != nil || req.MaxPerPage == 0 {
			req.MaxPerPage = 50 // default Value maxPerPage
		}

		req.Page, _ = strconv.Atoi(html.EscapeString(page))
		if req.Page == 0 {
			req.Page = 1
		}
	}

	nbTorrents := 0
	torrents, nbTorrents, err := torrentService.GetTorrents(whereParams, req.MaxPerPage, req.MaxPerPage*(req.Page-1))
	if err != nil {
		util.SendError(w, err, 400)
		return
	}

	b := model.ApiResultJson{
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
	b := torrent.ToJson()

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(b)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ApiUploadHandler(w http.ResponseWriter, r *http.Request) {
	if config.UploadsDisabled == 1 {
		http.Error(w, "Error uploads are disabled", http.StatusInternalServerError)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" {
		token := r.Header.Get("Authorization")
		user := model.User{}
		db.ORM.Where("api_token = ?", token).First(&user) //i don't like this
		if user.Id == 0 {
			http.Error(w, "incorrect api key", http.StatusForbidden)
			return
		}

		defer r.Body.Close()

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

		torrent := model.Torrents{
			Name:         upload.Name,
			Category:     upload.Category,
			Sub_Category: upload.SubCategory,
			Status:       1,
			Hash:         upload.Hash,
			Date:         time.Now(),
			Filesize:     0, //?
			Description:  upload.Description,
			UploaderId:   user.Id,
			Uploader:     &user,
		}

		db.ORM.Create(&torrent)
		fmt.Printf("%+v\n", torrent)
	}
}

func ApiUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if config.UploadsDisabled == 1 {
		http.Error(w, "Error uploads are disabled", http.StatusInternalServerError)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" {
		token := r.Header.Get("Authorization")
		user := model.User{}
		db.ORM.Where("api_token = ?", token).First(&user) //i don't like this
		if user.Id == 0 {
			http.Error(w, "incorrect api key", http.StatusForbidden)
			return
		}

		defer r.Body.Close()

		update := apiService.UpdateRequest{}
		d := json.NewDecoder(r.Body)
		if err := d.Decode(&update); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id := update.Id
		torrent := model.Torrents{}
		db.ORM.Where("torrent_id = ?", id).First(&torrent)
		if torrent.Id == 0 {
			http.Error(w, "incorrect id", http.StatusBadRequest)
			return
		}
		if torrent.UploaderId != 0 && torrent.UploaderId != user.Id { //&& user is not mod
			http.Error(w, "not enough rights to edit torrent", http.StatusForbidden)
			return
		}
		err, code := update.Update.ValidateUpdate()
		if err != nil {
			http.Error(w, err.Error(), code)
			return
		}
		fmt.Printf("%+v\n", torrent)
		update.UpdateTorrent(&torrent)

		db.ORM.Save(&torrent)
		fmt.Printf("%+v\n", torrent)
	}
}
