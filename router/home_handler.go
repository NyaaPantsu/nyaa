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
	"github.com/NyaaPantsu/nyaa/util/log"
	msg "github.com/NyaaPantsu/nyaa/util/messages"
	"github.com/gorilla/mux"
)

// HomeHandler : Controller for Home page, can have some query arguments
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]
	messages := msg.GetMessages(r)
	deleteVar := r.URL.Query()["deleted"]
	defer r.Body.Close()

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
	if deleteVar != nil {
		messages.AddInfoTf("infos", "torrent_deleted", "")
	}
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

	navigationTorrents := navigation{
		TotalItem:      nbTorrents,
		MaxItemPerPage: maxPerPage,
		CurrentPage:    pagenum,
		Route:          "search_page",
	}

	torrentsJSON := model.TorrentsToJSON(torrents)
	common := newCommonVariables(r)
	common.Navigation = navigationTorrents
	htv := modelListVbs{
		commonTemplateVariables: common,
		Models:                  torrentsJSON,
		Infos:                   messages.GetAllInfos(),
	}

	err = homeTemplate.ExecuteTemplate(w, "index.html", htv)
	if err != nil {
		log.Errorf("HomeHandler(): %s", err)
	}
}
