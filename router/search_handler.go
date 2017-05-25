package router

import (
	"html"
	"net/http"
	"strconv"

	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/util"
	"github.com/NyaaPantsu/nyaa/util/log"
	msg "github.com/NyaaPantsu/nyaa/util/messages"
	"github.com/NyaaPantsu/nyaa/util/search"
	"github.com/gorilla/mux"
)

// SearchHandler : Controller for displaying search result page, accepting common search arguments
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]
	messages := msg.GetMessages(r)

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

	common := newCommonVariables(r)
	common.Navigation = navigation{nbTorrents, int(searchParam.Max), pagenum, "search_page"}
	// Convert back to strings for now.
	common.Search = searchForm{
		SearchParam:      searchParam,
		Category:         searchParam.Category.String(),
		ShowItemsPerPage: true,
	}
	htv := modelListVbs{common, b, messages.GetAllErrors(), messages.GetAllInfos()}

	err = searchTemplate.ExecuteTemplate(w, "index.html", htv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
