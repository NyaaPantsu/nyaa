package router

import (
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/util"
	"github.com/NyaaPantsu/nyaa/util/languages"
	"github.com/NyaaPantsu/nyaa/util/log"
	"github.com/NyaaPantsu/nyaa/util/search"
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
		if pagenum <= 0 {
			NotFoundHandler(w, r)
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
	T := languages.GetTfuncFromRequest(r)
	htv := HomeTemplateVariables{b, searchForm, navigationTorrents, T, GetUser(r), r.URL, mux.CurrentRoute(r)}

	err = searchTemplate.ExecuteTemplate(w, "index.html", htv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
