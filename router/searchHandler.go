package router

import (
	"github.com/ewhal/nyaa/model"
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

	b := []model.TorrentsJson{}

	search_param, torrents, nbTorrents := search.SearchByQuery(r, pagenum)

	for i := range torrents {
		res := torrents[i].ToJson()
		b = append(b, res)
	}

	navigationTorrents := Navigation{nbTorrents, search_param.Max, pagenum, "search_page"}
	searchForm := SearchForm{
		search_param.Query,
		search_param.Status,
		search_param.Category,
		search_param.Sort,
		search_param.Order,
		false,
	}
	htv := HomeTemplateVariables{b, searchForm, navigationTorrents, r.URL, mux.CurrentRoute(r)}

	err := searchTemplate.ExecuteTemplate(w, "index.html", htv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
