package router

import (
	"html"
	"net/http"
	"strconv"

	"github.com/NyaaPantsu/nyaa/cache"
	"github.com/NyaaPantsu/nyaa/common"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/service/torrent"
	"github.com/NyaaPantsu/nyaa/util"
	"github.com/NyaaPantsu/nyaa/util/languages"
	"github.com/NyaaPantsu/nyaa/util/log"
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

	navigationTorrents := Navigation{
		TotalItem:      nbTorrents,
		MaxItemPerPage: maxPerPage,
		CurrentPage:    pagenum,
		Route:          "search_page",
	}

	languages.SetTranslationFromRequest(homeTemplate, r)

	torrentsJson := model.TorrentsToJSON(torrents)
	htv := HomeTemplateVariables{
		ListTorrents: torrentsJson,
		Search:       NewSearchForm(),
		Navigation:   navigationTorrents,
		User:         GetUser(r),
		URL:          r.URL,
		Route:        mux.CurrentRoute(r),
	}

	err = homeTemplate.ExecuteTemplate(w, "index.html", htv)
	if err != nil {
		log.Errorf("HomeHandler(): %s", err)
	}
}
