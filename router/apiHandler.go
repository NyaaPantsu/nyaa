package router

import(
	"github.com/gorilla/mux"
	"net/http"
	"html"
	"strconv"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/service/torrent"
	"encoding/json"
)

func ApiHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	page := vars["page"]
	pagenum, _ := strconv.Atoi(html.EscapeString(page))

	b := model.CategoryJson{Torrents: []model.TorrentsJson{}}
	maxPerPage := 50
	nbTorrents := 0

	torrents, nbTorrents := torrentService.GetAllTorrents(maxPerPage, maxPerPage*(pagenum-1))
	for i, _ := range torrents {
		res := torrents[i].ToJson()
		b.Torrents = append(b.Torrents, res)
	}

	b.QueryRecordCount = maxPerPage
	b.TotalRecordCount = nbTorrents
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(b)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ApiViewHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	b := model.CategoryJson{Torrents: []model.TorrentsJson{}}

	torrent, err := torrentService.GetTorrentById(id)
	res := torrent.ToJson()
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