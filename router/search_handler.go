package router

import (
	"github.com/gorilla/mux"
	"html"
	"net/http"
	"strconv"

	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/util"
	"github.com/NyaaPantsu/nyaa/util/log"
	msg "github.com/NyaaPantsu/nyaa/util/messages"
	"github.com/NyaaPantsu/nyaa/util/search"
)

// SearchHandler : Controller for displaying search result page, accepting common search arguments
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var err error
	// TODO Don't create a new client for each request
	messages := msg.GetMessages(r)
	// TODO Fallback to postgres search if es is down

	vars := mux.Vars(r)
	page := vars["page"]

	// db params url
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

	commonVar := newCommonVariables(r)
	commonVar.Navigation = navigation{int(nbTorrents), int(searchParam.Max), int(searchParam.Page), "search_page"}
	// Convert back to strings for now.
	// TODO Deprecate fully SearchParam and only use TorrentParam
	commonVar.Search = searchForm{
		SearchParam:      searchParam,
		Category:         searchParam.Category.String(),
		ShowItemsPerPage: true,
	}
	htv := modelListVbs{commonVar, model.TorrentsToJSON(torrents), messages.GetAllErrors(), messages.GetAllInfos()}

	err = searchTemplate.ExecuteTemplate(w, "index.html", htv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
