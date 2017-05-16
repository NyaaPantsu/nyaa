package router

import (
	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/feeds"
	"github.com/ewhal/nyaa/util"
	"github.com/ewhal/nyaa/util/search"
	"net/http"
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
			Id:          "https://" + config.WebAddress + "/view/" + torrentJSON.ID,
			Title:       torrent.Name,
			Link:        &feeds.Link{Href: string(torrentJSON.Magnet)},
			Description: string(torrentJSON.Description),
			Created:     torrent.Date,
			Updated:     torrent.Date,
			Torrent: &feeds.Torrent{
				Seeders:    torrent.Seeders,
				Leechers:   torrent.Leechers,
				Hash:       torrent.Hash,
				Completed:  torrent.Completed,
				LastScrape: torrent.LastScrape,
			},
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
