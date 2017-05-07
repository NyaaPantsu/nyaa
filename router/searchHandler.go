package router

import (
	"encoding/json"
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

	should_json := r.URL.Query().Get("t") == "json"

	b := []model.TorrentsJson{}

	var err error

	if should_json {
		_, torrents := search.SearchByQueryNoCount(r, pagenum)
		for i := range torrents {
			res := torrents[i].ToJson()
			b = append(b, res)
		}
		w.Header().Set("Content-Type", "text/json; encoding=UTF-8")

		err = json.NewEncoder(w).Encode(b)
	} else {
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

		err = searchTemplate.ExecuteTemplate(w, "index.html", htv)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
