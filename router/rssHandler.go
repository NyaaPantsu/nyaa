package router

import (
	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/util/search"
	"github.com/gorilla/feeds"
	"net/http"
	"strconv"
	"time"
)

func RssHandler(w http.ResponseWriter, r *http.Request) {

	_, torrents, _ := search.SearchByQuery(r, 1)
	created_as_time := time.Now()

	if len(torrents) > 0 {
		created_as_time = torrents[0].Date
	}
	feed := &feeds.Feed{
		Title:   "Nyaa Pantsu",
		Link:    &feeds.Link{Href: "https://" + config.WebAddress + "/"},
		Created: created_as_time,
	}
	feed.Items = []*feeds.Item{}
	feed.Items = make([]*feeds.Item, len(torrents))

	for i, _ := range torrents {
		torrent_json := torrents[i].ToJSON()
		feed.Items[i] = &feeds.Item{
			// need a torrent view first
			Id:          "https://" + config.WebAddress + "/view/" + strconv.FormatUint(uint64(torrents[i].ID), 10),
			Title:       torrents[i].Name,
			Link:        &feeds.Link{Href: string(torrent_json.Magnet)},
			Description: "",
			Created:     torrents[0].Date,
			Updated:     torrents[0].Date,
		}
	}

	rss, err := feed.ToRss()
	if err == nil {
		w.Write([]byte(rss))
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
