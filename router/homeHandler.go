package router

import (
	"html"
	"net/http"
	"strconv"

	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/service/torrent"
	"github.com/ewhal/nyaa/util"
	"github.com/ewhal/nyaa/util/languages"
	"github.com/ewhal/nyaa/util/log"
	"github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]

	// db params url
	maxPerPage, errConv := strconv.Atoi(r.URL.Query().Get("max"))
	if errConv != nil {
		maxPerPage = 50 // default Value maxPerPage
	}

	nbTorrents := 0
	pagenum, _ := strconv.Atoi(html.EscapeString(page))
	if pagenum == 0 {
		pagenum = 1
	}

	torrents, nbTorrents, err := torrentService.GetAllTorrents(maxPerPage, maxPerPage*(pagenum-1))
	if err != nil {
		util.SendError(w, err, 400)
		return
	}

	b := model.TorrentsToJSON(torrents)

	navigationTorrents := Navigation{nbTorrents, maxPerPage, pagenum, "search_page"}

	languages.SetTranslationFromRequest(homeTemplate, r, "en-us")
	htv := HomeTemplateVariables{b, NewSearchForm(), navigationTorrents, GetUser(r), r.URL, mux.CurrentRoute(r)}

	err = homeTemplate.ExecuteTemplate(w, "index.html", htv)
	if err != nil {
		log.Errorf("HomeHandler(): %s", err)
	}
}
