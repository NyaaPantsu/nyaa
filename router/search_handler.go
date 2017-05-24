package router

import (
	"net/http"

	"github.com/NyaaPantsu/nyaa/common"
	"github.com/NyaaPantsu/nyaa/util/log"
	msg "github.com/NyaaPantsu/nyaa/util/messages"
	elastic "gopkg.in/olivere/elastic.v5"
)

// SearchHandler : Controller for displaying search result page, accepting common search arguments
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	// TODO Don't create a new client for each request
	client, err := elastic.NewClient()
	if err != nil {
		log.Errorf("Unable to create elasticsearch client: %s\n", err)
	}
	var torrentParam common.TorrentParam
	torrentParam.FromRequest(r)
	totalHits, torrents, err := torrentParam.Find(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	messages := msg.GetMessages(r)
	// TODO Fallback to postgres search if es is down

	commonVar := newCommonVariables(r)
	commonVar.Navigation = navigation{int(totalHits), int(torrentParam.Max), int(torrentParam.Offset), "search_page"}
	// Convert back to strings for now.
	// Convert back to strings for now.
	// TODO Deprecate fully SearchParam and only use TorrentParam
	searchParam := common.SearchParam{
		Order: torrentParam.Order,
		Status: torrentParam.Status,
		Sort: torrentParam.Sort,
		Category: torrentParam.Category,
		Page: int(torrentParam.Offset),
		UserID: uint(torrentParam.UserID),
		Max: uint(torrentParam.Max),
		NotNull: torrentParam.NotNull,
		Query: torrentParam.NameLike,
	}

	commonVar.Search = searchForm{
		SearchParam:      searchParam,
		Category:         searchParam.Category.String(),
		ShowItemsPerPage: true,
	}
	htv := modelListVbs{commonVar, torrents, messages.GetAllErrors(), messages.GetAllInfos()}

	err = searchTemplate.ExecuteTemplate(w, "index.html", htv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
