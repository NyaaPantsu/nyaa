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
	token := r.Header.Get("Authorization")
	user := model.User{}
	db.ORM.Where("api_token = ?", token).First(&user) //i don't like this
	if config.UploadsDisabled && config.AdminsAreStillAllowedTo && user.Status != 2 && config.TrustedUsersAreStillAllowedTo && user.Status != 1 {
		http.Error(w, "Error uploads are disabled", http.StatusBadRequest)
		return
	} else if config.UploadsDisabled && !config.AdminsAreStillAllowedTo && user.Status == 2 {
		http.Error(w, "Error uploads are disabled", http.StatusBadRequest)
		return
	} else if config.UploadsDisabled && !config.TrustedUsersAreStillAllowedTo && user.Status == 1 {
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
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

func ApiUpdateHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	user := model.User{}
	db.ORM.Where("api_token = ?", token).First(&user) //i don't like this
	if config.UploadsDisabled && config.AdminsAreStillAllowedTo && user.Status != 2 && config.TrustedUsersAreStillAllowedTo && user.Status != 1 {
		http.Error(w, "Error uploads are disabled", http.StatusInternalServerError)
		return
	} else if config.UploadsDisabled && !config.AdminsAreStillAllowedTo && user.Status == 2 {
		http.Error(w, "Error uploads are disabled", http.StatusInternalServerError)
		return
	} else if config.UploadsDisabled && !config.TrustedUsersAreStillAllowedTo && user.Status == 1 {
		http.Error(w, "Error uploads are disabled", http.StatusInternalServerError)
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
