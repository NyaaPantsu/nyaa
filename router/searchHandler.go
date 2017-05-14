package router

import (
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/util"
	"github.com/ewhal/nyaa/util/languages"
	"github.com/ewhal/nyaa/util/log"
	"github.com/ewhal/nyaa/util/search"
	"github.com/gorilla/mux"
	"html"
	"net/http"
	"strconv"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]

	// db params url
	var err error
	pagenum := 1
	if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if !log.CheckError(err) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	searchParam, torrents, nbTorrents, err := search.SearchByQuery(r, pagenum)
	if err != nil {
		util.SendError(w, err, 400)
		return
	}

	b := model.TorrentsToJSON(torrents)

	navigationTorrents := Navigation{nbTorrents, int(searchParam.Max), pagenum, "search_page"}
	// Convert back to strings for now.
	searchForm := SearchForm{
		SearchParam:      searchParam,
		Category:         searchParam.Category.String(),
		ShowItemsPerPage: true,
	}
	htv := HomeTemplateVariables{b, searchForm, navigationTorrents, GetUser(r), r.URL, mux.CurrentRoute(r)}

	languages.SetTranslationFromRequest(searchTemplate, r, "en-us")
	err = searchTemplate.ExecuteTemplate(w, "index.html", htv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
