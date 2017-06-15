package router

import (
	"errors"
	"html"
	"net/http"
	"strconv"
	"time"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/model"
	userService "github.com/NyaaPantsu/nyaa/service/user"
	"github.com/NyaaPantsu/nyaa/util/feeds"
	"github.com/NyaaPantsu/nyaa/util/search"
	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
)

// RSSHandler : Controller for displaying rss feed, accepting common search arguments
func RSSHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// We only get the basic variable for rss based on search param
	torrents, createdAsTime, title, err := getTorrentList(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	feed := &nyaafeeds.RssFeed{
		Title:   title,
		Link:    config.WebAddress() + "/",
		PubDate: createdAsTime.String(),
	}
	feed.Items = make([]*nyaafeeds.RssItem, len(torrents))

	for i, torrent := range torrents {
		torrentJSON := torrent.ToJSON()
		feed.Items[i] = &nyaafeeds.RssItem{
			Title:       torrentJSON.Name,
			Link:        config.WebAddress() + "/view/" + strconv.FormatUint(uint64(torrentJSON.ID), 10),
			Description: string(torrentJSON.Description),
			Author:      config.WebAddress() + "/view/" + strconv.FormatUint(uint64(torrentJSON.ID), 10),
			PubDate:     torrent.Date.String(),
			GUID:        config.WebAddress() + "/download/" + torrentJSON.Hash,
			Enclosure: &nyaafeeds.RssEnclosure{
				URL:    config.WebAddress() + "/download/" + torrentJSON.Hash,
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

// RSSEztvHandler : Controller for displaying rss feed, accepting common search arguments
func RSSEztvHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// We only get the basic variable for rss based on search param
	torrents, createdAsTime, title, err := getTorrentList(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	feed := &nyaafeeds.RssFeed{
		Title:   title,
		Link:    config.WebAddress() + "/",
		PubDate: createdAsTime.String(),
	}
	feed.Items = make([]*nyaafeeds.RssItem, len(torrents))

	for i, torrent := range torrents {
		torrentJSON := torrent.ToJSON()
		feed.Items[i] = &nyaafeeds.RssItem{
			Title: torrentJSON.Name,
			Link:  config.WebAddress() + "/download/" + torrentJSON.Hash,
			Category: &nyaafeeds.RssCategory{
				Domain: config.WebAddress() + "/search?c=" + torrentJSON.Category + "_" + torrentJSON.SubCategory,
			},
			Description: string(torrentJSON.Description),
			Comments:    config.WebAddress() + "/view/" + strconv.FormatUint(uint64(torrentJSON.ID), 10),
			PubDate:     torrent.Date.String(),
			GUID:        config.WebAddress() + "/view/" + strconv.FormatUint(uint64(torrentJSON.ID), 10),
			Enclosure: &nyaafeeds.RssEnclosure{
				URL:    config.WebAddress() + "/download/" + torrentJSON.Hash,
				Length: strconv.FormatUint(uint64(torrentJSON.Filesize), 10),
				Type:   "application/x-bittorrent",
			},
			Torrent: &nyaafeeds.RssTorrent{
				Xmlns:         "http://xmlns.ezrss.it/0.1/",
				FileName:      torrentJSON.Name + ".torrent",
				ContentLength: strconv.FormatUint(uint64(torrentJSON.Filesize), 10),
				InfoHash:      torrentJSON.Hash,
				MagnetURI:     string(torrentJSON.Magnet),
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

// RSSEztvHandler : Controller for displaying rss feed, accepting common search arguments
func RSSTorznabHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// We only get the basic variable for rss based on search param
	torrents, createdAsTime, title, err := getTorrentList(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	feed := &nyaafeeds.RssFeed{
		Title:   title,
		Link:    config.WebAddress() + "/",
		PubDate: createdAsTime.String(),
	}
	feed.Items = make([]*nyaafeeds.RssItem, len(torrents))

	for i, torrent := range torrents {
		torrentJSON := torrent.ToJSON()
		feed.Items[i] = &nyaafeeds.RssItem{
			Title: torrentJSON.Name,
			Link:  config.WebAddress() + "/download/" + torrentJSON.Hash,
			Category: &nyaafeeds.RssCategory{
				Domain: config.WebAddress() + "/search?c=" + torrentJSON.Category + "_" + torrentJSON.SubCategory,
			},
			Description: string(torrentJSON.Description),
			Comments:    config.WebAddress() + "/view/" + strconv.FormatUint(uint64(torrentJSON.ID), 10),
			PubDate:     torrent.Date.String(),
			GUID:        config.WebAddress() + "/view/" + strconv.FormatUint(uint64(torrentJSON.ID), 10),
			Enclosure: &nyaafeeds.RssEnclosure{
				URL:    config.WebAddress() + "/download/" + torrentJSON.Hash,
				Length: strconv.FormatUint(uint64(torrentJSON.Filesize), 10),
				Type:   "application/x-bittorrent",
			},
			Torznab: &nyaafeeds.RssTorznab{
				Xmlns:     "http://torznab.com/schemas/2015/feed",
				Size:      strconv.FormatUint(uint64(torrentJSON.Filesize), 10),
				Files:     strconv.Itoa(len(torrentJSON.FileList)),
				Grabs:     strconv.Itoa(torrentJSON.Downloads),
				Seeders:   strconv.Itoa(int(torrentJSON.Seeders)),
				Leechers:  strconv.Itoa(int(torrentJSON.Leechers)),
				Infohash:  torrentJSON.Hash,
				MagnetURL: string(torrentJSON.Magnet),
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

func getTorrentList(r *http.Request) (torrents []model.Torrent, createdAsTime time.Time, title string, err error) {
	vars := mux.Vars(r)
	page := vars["page"]
	userID := vars["id"]

	offset := r.URL.Query().Get("offset")
	pagenum := 1
	if page == "" && offset != "" {
		page = offset
	}
	if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if err != nil {
			return
		}
		if pagenum <= 0 {
			err = errors.New("Page number is invalid")
			return
		}
	}

	if userID != "" {
		userIDnum := 0
		userIDnum, err = strconv.Atoi(html.EscapeString(userID))
		// Should we have a feed for anonymous uploads?
		if err != nil || userIDnum == 0 {
			return
		}

		_, _, err = userService.RetrieveUserForAdmin(userID)
		if err != nil {
			return
		}
		// Set the user ID on the request, so that SearchByQuery finds it.
		query := r.URL.Query()
		query.Set("userID", userID)
		r.URL.RawQuery = query.Encode()
	}

	_, torrents, err = search.SearchByQueryNoCount(r, pagenum)

	createdAsTime = time.Now()

	if len(torrents) > 0 {
		createdAsTime = torrents[0].Date
	}

	title = "Nyaa Pantsu"
	if config.IsSukebei() {
		title = "Sukebei Pantsu"
	}

	return
}
