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
	"github.com/ewhal/nyaa/service/torrent"
	"github.com/ewhal/nyaa/util"
	"github.com/gorilla/mux"
)

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]

	maxPerPage, errConv := strconv.Atoi(r.URL.Query().Get("max"))
	if errConv != nil {
		maxPerPage = 50 // default Value maxPerPage
	}

	pagenum, _ := strconv.Atoi(html.EscapeString(page))
	if pagenum == 0 {
		pagenum = 1
	}

	torrents, nbTorrents, err := torrentService.GetAllTorrents(maxPerPage, maxPerPage*(pagenum-1))
	if err != nil {
		util.SendError(w, err, 400)
		return
	}

	b := model.ApiResultJSON{
		Torrents: model.TorrentsToJSON(torrents),
	}
	b.QueryRecordCount = maxPerPage
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
	b := model.ApiResultJSON{Torrents: []model.TorrentJSON{}}

	torrent, err := torrentService.GetTorrentById(id)
	res := torrent.ToJSON()
	b.Torrents = append(b.Torrents, res)

	b.QueryRecordCount = 1
	b.TotalRecordCount = 1
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

	defer r.Body.Close()

	// TODO: verify token
	//token := r.Header.Get("Authorization")

	decoder := json.NewDecoder(r.Body)
	torrentJSON := model.TorrentJSON{}
	err := decoder.Decode(&torrentJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	category, subCategory, err := ValidateJSON(&torrentJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	// TODO: Interface changed. Structure is incomplete.
	// TODO: Model package should provide conversion utility
	torrent := model.Torrent{
		ID:          0,
		Name:        torrentJSON.Name,
		Hash:        torrentJSON.Hash,
		Category:    category,
		SubCategory: subCategory,
		Status:      1,
		Date:        time.Now(),
		UploaderID:  0,
		Downloads:   0,
		Stardom:     0,
		Filesize:    0,
		Description: string(torrentJSON.Description),
		WebsiteLink: "N/A",
		Uploader:    nil,
		OldComments: nil,
		Comments:    nil}

	db.ORM.Create(&torrent)
	fmt.Printf("%+v\n", torrent)
}
