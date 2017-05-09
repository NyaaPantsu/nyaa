package router

import (
	"html"
	"net/http"
	"strconv"

	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/util"
	"github.com/ewhal/nyaa/util/languages"
	"github.com/ewhal/nyaa/util/search"
	"github.com/gorilla/mux"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]

	// db params url
	pagenum, _ := strconv.Atoi(html.EscapeString(page))
	if pagenum == 0 {
		pagenum = 1
	}

	search_param, torrents, nbTorrents, err := search.SearchByQuery(r, pagenum)
	if err != nil {
		util.SendError(w, err, 400)
		return
	}

	b := model.TorrentsToJSON(torrents)

	navigationTorrents := Navigation{nbTorrents, int(search_param.Max), pagenum, "search_page"}
	// Convert back to strings for now.
	searchForm := SearchForm{
		SearchParam:        search_param,
		Category:           search_param.Category.String(),
		HideAdvancedSearch: false,
	}
	htv := HomeTemplateVariables{b, searchForm, navigationTorrents, GetUser(r), r.URL, mux.CurrentRoute(r)}

	languages.SetTranslationFromRequest(searchTemplate, r, "en-us")
	err = searchTemplate.ExecuteTemplate(w, "index.html", htv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
