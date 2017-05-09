package router

import (
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/service/torrent"
	"github.com/ewhal/nyaa/util/languages"
	"github.com/ewhal/nyaa/util/log"
	"github.com/gorilla/mux"
	"html"
	"net/http"
	"strconv"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]

	// db params url
	maxPerPage, errConv := strconv.Atoi(r.URL.Query().Get("max"))
	if errConv != nil {
		maxPerPage = 50 // default Value maxPerPage
	}

	pagenum, _ := strconv.Atoi(html.EscapeString(page))
	if pagenum == 0 {
		pagenum = 1
	}

	b := []model.TorrentJSON{}
	torrents, nbTorrents := torrentService.GetAllTorrents(maxPerPage, maxPerPage*(pagenum-1))

	for i := range torrents {
		res := torrents[i].ToJSON()
		b = append(b, res)
	}

	navigationTorrents := Navigation{nbTorrents, maxPerPage, pagenum, "search_page"}

	languages.SetTranslationFromRequest(homeTemplate, r, "en-us")
	htv := HomeTemplateVariables{b, NewSearchForm(), navigationTorrents, GetUser(r), r.URL, mux.CurrentRoute(r)}

	err := homeTemplate.ExecuteTemplate(w, "index.html", htv)
	if err != nil {
		log.Errorf("HomeHandler(): %s", err)
	}
}
