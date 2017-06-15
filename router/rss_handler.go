package router

import (
	"html"
	"net/http"
	"strconv"
	"time"

	"github.com/NyaaPantsu/nyaa/config"
	userService "github.com/NyaaPantsu/nyaa/service/user"
	"github.com/NyaaPantsu/nyaa/util/search"
	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
)

// RSSHandler : Controller for displaying rss feed, accepting common search arguments
func RSSHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	page := vars["page"]
	userID := vars["id"]

	offset := r.URL.Query().Get("offset")
	var err error
	pagenum := 1
	if page == "" && offset != "" {
		page = offset
	}
	if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if pagenum <= 0 {
			NotFoundHandler(w, r)
			return
		}
	}

	if userID != "" {
		userIDnum, err := strconv.Atoi(html.EscapeString(userID))
		// Should we have a feed for anonymous uploads?
		if err != nil || userIDnum == 0 {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, _, err = userService.RetrieveUserForAdmin(userID)
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		// Set the user ID on the request, so that SearchByQuery finds it.
		query := r.URL.Query()
		query.Set("userID", userID)
		r.URL.RawQuery = query.Encode()
	}

	_, torrents, err := search.SearchByQueryNoCount(r, pagenum)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	createdAsTime := time.Now()

	if len(torrents) > 0 {
		createdAsTime = torrents[0].Date
	}
	title := "Nyaa Pantsu"
	if config.IsSukebei() {
		title = "Sukebei Pantsu"
	}
	feed := &feeds.RssFeed{
		Title:   title,
		Link:    config.WebAddress() + "/",
		PubDate: createdAsTime.String(),
	}
	feed.Items = make([]*feeds.RssItem, len(torrents))

	for i, torrent := range torrents {
		torrentJSON := torrent.ToJSON()
		feed.Items[i] = &feeds.RssItem{
			Title:       torrentJSON.Name,
			Link:        config.WebAddress() + "/view/" + strconv.FormatUint(uint64(torrentJSON.ID), 10),
			Description: string(torrentJSON.Description),
			Author:      config.WebAddress() + "/view/" + strconv.FormatUint(uint64(torrentJSON.ID), 10),
			PubDate:     torrent.Date.String(),
			Guid:        config.WebAddress() + "/download/" + torrentJSON.Hash,
			Enclosure: &feeds.RssEnclosure{
				Url:    config.WebAddress() + "/download/" + torrentJSON.Hash,
				Length: strconv.FormatUint(uint64(torrentJSON.Filesize), 10),
				Type:   "application/x-bittorrent",
			},
		}
	}
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "*")
	rss, rssErr := feeds.ToXML(feed)
	if rssErr != nil {
		http.Error(w, rssErr.Error(), http.StatusInternalServerError)
	}

	_, writeErr := w.Write([]byte(rss))
	if writeErr != nil {
		http.Error(w, writeErr.Error(), http.StatusInternalServerError)
	}
}
