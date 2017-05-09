package router

import (
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/util/languages"
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
	pagenum, _ := strconv.Atoi(html.EscapeString(page))
	if pagenum == 0 {
		pagenum = 1
	}

	b := []model.TorrentJSON{}

	searchParam, torrents, nbTorrents := search.SearchByQuery(r, pagenum)

	for i := range torrents {
		res := torrents[i].ToJSON()
		b = append(b, res)
	}

	navigationTorrents := Navigation{nbTorrents, searchParam.Max, pagenum, "search_page"}
	searchForm := SearchForm{
		searchParam.Query,
		searchParam.Status,
		searchParam.Category,
		searchParam.Sort,
		searchParam.Order,
		false,
	}
	htv := HomeTemplateVariables{b, searchForm, navigationTorrents, GetUser(r), r.URL, mux.CurrentRoute(r)}

	languages.SetTranslationFromRequest(searchTemplate, r, "en-us")
	err := searchTemplate.ExecuteTemplate(w, "index.html", htv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
