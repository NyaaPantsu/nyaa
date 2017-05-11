package router

import (
	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/util"
	"github.com/ewhal/nyaa/util/search"
	"github.com/gorilla/feeds"
	"net/http"
	"strconv"
	"time"
)

func RSSHandler(w http.ResponseWriter, r *http.Request) {
	_, torrents, err := search.SearchByQueryNoCount(r, 1)
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
			Id:          "https://" + config.WebAddress + "/view/" + strconv.FormatUint(uint64(torrents[i].ID), 10),
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, writeErr := w.Write([]byte(rss))
	if writeErr != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
