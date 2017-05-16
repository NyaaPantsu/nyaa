package router

import (
	"html"
	"net/http"
	"strconv"

	"github.com/ewhal/nyaa/cache"
	"github.com/ewhal/nyaa/common"
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
	var err error
	maxPerPage := 50
	maxString := r.URL.Query().Get("max")
	if maxString != "" {
		maxPerPage, err = strconv.Atoi(maxString)
		if !log.CheckError(err) {
			maxPerPage = 50 // default Value maxPerPage
		}
	}

	pagenum := 1
	if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if !log.CheckError(err) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	search := common.SearchParam{
		Max:  uint(maxPerPage),
		Page: pagenum,
	}

	torrents, nbTorrents, err := cache.Impl.Get(search, func() ([]model.Torrent, int, error) {
		torrents, nbTorrents, err := torrentService.GetAllTorrents(maxPerPage, maxPerPage*(pagenum-1))
		if !log.CheckError(err) {
			util.SendError(w, err, 400)
		}
		return torrents, nbTorrents, err
	})

	b := model.TorrentsToJSON(torrents)

	navigationTorrents := Navigation{nbTorrents, maxPerPage, pagenum, "search_page"}

	languages.SetTranslationFromRequest(homeTemplate, r)
	htv := HomeTemplateVariables{b, NewSearchForm(), navigationTorrents, GetUser(r), r.URL, mux.CurrentRoute(r)}

	err = homeTemplate.ExecuteTemplate(w, "index.html", htv)
	if err != nil {
		log.Errorf("HomeHandler(): %s", err)
	}
}
