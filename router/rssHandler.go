package router

import(
	"time"
	"net/http"
	"github.com/gorilla/feeds"
	"github.com/ewhal/nyaa/util/search"
	"strconv"
)

func RssHandler(w http.ResponseWriter, r *http.Request) {

	_, torrents, _ := search.SearchByQuery( r, 1 )
	created_as_time := time.Now()

	if len(torrents) > 0 {
		created_as_time = time.Unix(torrents[0].Date, 0)
	}
	feed := &feeds.Feed{
		Title:   "Nyaa Pantsu",
		Link:    &feeds.Link{Href: "https://nyaa.pantsu.cat/"},
		Created: created_as_time,
	}
	feed.Items = []*feeds.Item{}
	feed.Items = make([]*feeds.Item, len(torrents))

	for i, _ := range torrents {
		timestamp_as_time := time.Unix(torrents[0].Date, 0)
		torrent_json := torrents[i].ToJson()
		feed.Items[i] = &feeds.Item{
			// need a torrent view first
			Id:          "https://nyaa.pantsu.cat/view/" + strconv.Itoa(torrents[i].Id),
			Title:       torrents[i].Name,
			Link:        &feeds.Link{Href: string(torrent_json.Magnet)},
			Description: "",
			Created:     timestamp_as_time,
			Updated:     timestamp_as_time,
		}
	}

	rss, err := feed.ToRss()
	if err == nil {
		w.Write([]byte(rss))
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}