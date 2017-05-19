package router

import (
	"github.com/NyaaPantsu/nyaa/config"
	userService "github.com/NyaaPantsu/nyaa/service/user"
	"github.com/NyaaPantsu/nyaa/util"
	"github.com/NyaaPantsu/nyaa/util/search"
	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
	"html"
	"net/http"
	"strconv"
	"time"
)

func RSSHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]
	userID := vars["id"]

	var err error
	pagenum := 1
	if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if err != nil {
			util.SendError(w, err, 400)
			return
		}
	}

	if userID != "" {
		userIDnum, err := strconv.Atoi(html.EscapeString(userID))
		// Should we have a feed for anonymous uploads?
		if err != nil || userIDnum == 0 {
			util.SendError(w, err, 400)
			return
		}

		_, _, err = userService.RetrieveUserForAdmin(userID)
		if err != nil {
			util.SendError(w, err, 404)
			return
		}

		// Set the user ID on the request, so that SearchByQuery finds it.
		query := r.URL.Query()
		query.Set("userID", userID)
		r.URL.RawQuery = query.Encode()
	}

	_, torrents, err := search.SearchByQueryNoCount(r, pagenum)
	if err != nil {
		util.SendError(w, err, 400)
		return
	}
	createdAsTime := time.Now()

	if len(torrents) > 0 {
		createdAsTime = torrents[0].Date
	}
	feed := &feeds.Feed{
		Title:   "Nyaa Pantsu",
		Link:    &feeds.Link{Href: "https://" + config.WebAddress + "/"},
		Created: createdAsTime,
	}
	feed.Items = make([]*feeds.Item, len(torrents))

	for i, torrent := range torrents {
		torrentJSON := torrent.ToJSON()
		feed.Items[i] = &feeds.Item{
			Id:          "https://" + config.WebAddress + "/view/" + torrentJSON.ID,
			Title:       torrent.Name,
			Link:        &feeds.Link{Href: string(torrentJSON.Magnet)},
			Description: string(torrentJSON.Description),
			Created:     torrent.Date,
			Updated:     torrent.Date,
		}
	}
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "*")
	rss, rssErr := feed.ToRss()
	if rssErr != nil {
		http.Error(w, rssErr.Error(), http.StatusInternalServerError)
	}

	_, writeErr := w.Write([]byte(rss))
	if writeErr != nil {
		http.Error(w, writeErr.Error(), http.StatusInternalServerError)
	}
}
