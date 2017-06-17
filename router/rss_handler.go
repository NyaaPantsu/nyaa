package router

import (
	"errors"
	"html"
	"net/http"
	"strconv"
	"time"

	"sort"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/model"
	userService "github.com/NyaaPantsu/nyaa/service/user"
	"github.com/NyaaPantsu/nyaa/util/categories"
	"github.com/NyaaPantsu/nyaa/util/feeds"
	"github.com/NyaaPantsu/nyaa/util/publicSettings"
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

// RSSTorznabHandler : Controller for displaying rss feed, accepting common search arguments
func RSSTorznabHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	t := r.URL.Query().Get("t")
	rss := ""
	title := "Nyaa Pantsu"
	if config.IsSukebei() {
		title = "Sukebei Pantsu"
	}
	if t == "caps" {
		T := publicSettings.GetTfuncFromRequest(r)
		cat := categories.GetCategoriesSelect(true, true)
		var categories []*nyaafeeds.RssCategoryTorznab
		categories = append(categories, &nyaafeeds.RssCategoryTorznab{
			ID:          "5070",
			Name:        "Anime",
			Description: "Anime",
		})
		var keys []string
		for name := range cat {
			keys = append(keys, name)
		}
		sort.Strings(keys)
		last := 0
		for _, key := range keys {
			if len(cat[key]) <= 2 {
				categories = append(categories, &nyaafeeds.RssCategoryTorznab{
					ID:   nyaafeeds.ConvertFromCat(cat[key]),
					Name: string(T(key)),
				})
				last++
			} else {
				categories[last].Subcat = append(categories[last].Subcat, &nyaafeeds.RssSubCat{
					ID:   nyaafeeds.ConvertFromCat(cat[key]),
					Name: string(T(key)),
				})
			}
		}
		feed := &nyaafeeds.RssCaps{
			Server: &nyaafeeds.RssServer{
				Version:   "1.0",
				Title:     title,
				Strapline: "...",
				Email:     config.Conf.Email.From,
				URL:       config.WebAddress(),
				Image:     config.WebAddress() + "/img/logo.png",
			},
			Limits: &nyaafeeds.RssLimits{
				Max:     "300",
				Default: "50",
			},
			Registration: &nyaafeeds.RssRegistration{
				Available: "yes",
				Open:      "yes",
			},
			Searching: &nyaafeeds.RssSearching{
				Search: &nyaafeeds.RssSearch{
					Available:       "yes",
					SupportedParams: "q",
				},
				TvSearch: &nyaafeeds.RssSearch{
					Available: "no",
				},
				MovieSearch: &nyaafeeds.RssSearch{
					Available: "no",
				},
			},
			Categories: &nyaafeeds.RssCategories{
				Category: categories,
			},
		}
		var rssErr error
		rss, rssErr = feeds.ToXML(feed)
		if rssErr != nil {
			http.Error(w, rssErr.Error(), http.StatusInternalServerError)
		}
	} else {
		// We only get the basic variable for rss based on search param
		torrents, createdAsTime, title, err := getTorrentList(r)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		feed := &nyaafeeds.RssFeed{
			Title:   title,
			Xmlns:   "http://torznab.com/schemas/2015/feed",
			Link:    config.WebAddress() + "/",
			PubDate: createdAsTime.String(),
		}
		feed.Items = make([]*nyaafeeds.RssItem, len(torrents))

		for i, torrent := range torrents {

			torrentJSON := torrent.ToJSON()
			filesNumber := ""
			if len(torrentJSON.FileList) > 0 {
				filesNumber = strconv.Itoa(len(torrentJSON.FileList))
			}
			seeders := ""
			if torrentJSON.Seeders > 0 {
				seeders = strconv.Itoa(int(torrentJSON.Seeders))
			}
			leechers := ""
			if torrentJSON.Leechers > 0 {
				leechers = strconv.Itoa(int(torrentJSON.Leechers))
			}
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
				Torznab: []*nyaafeeds.RssTorznab{
					&nyaafeeds.RssTorznab{
						Name:  "size",
						Value: strconv.FormatUint(uint64(torrentJSON.Filesize), 10),
					},
					&nyaafeeds.RssTorznab{
						Name:  "files",
						Value: filesNumber,
					},
					&nyaafeeds.RssTorznab{
						Name:  "grabs",
						Value: strconv.Itoa(torrentJSON.Downloads),
					},
					&nyaafeeds.RssTorznab{
						Name:  "seeders",
						Value: seeders,
					},
					&nyaafeeds.RssTorznab{
						Name:  "leechers",
						Value: leechers,
					},
					&nyaafeeds.RssTorznab{
						Name:  "infohash",
						Value: torrentJSON.Hash,
					},
					&nyaafeeds.RssTorznab{
						Name:  "magneturl",
						Value: string(torrentJSON.Magnet),
					},
				},
			}
		}
		var rssErr error
		rss, rssErr = feeds.ToXML(feed)
		if rssErr != nil {
			http.Error(w, rssErr.Error(), http.StatusInternalServerError)
		}
	}
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "*")

	_, writeErr := w.Write([]byte(rss))
	if writeErr != nil {
		http.Error(w, writeErr.Error(), http.StatusInternalServerError)
	}
}

func getTorrentList(r *http.Request) (torrents []model.Torrent, createdAsTime time.Time, title string, err error) {
	vars := mux.Vars(r)
	page := vars["page"]
	userID := vars["id"]
	cat := r.URL.Query().Get("cat")
	offset := 0
	if r.URL.Query().Get("offset") != "" {
		offset, err = strconv.Atoi(html.EscapeString(r.URL.Query().Get("offset")))
		if err != nil {
			return
		}
	}

	createdAsTime = time.Now()

	if len(torrents) > 0 {
		createdAsTime = torrents[0].Date
	}

	title = "Nyaa Pantsu"
	if config.IsSukebei() {
		title = "Sukebei Pantsu"
	}

	pagenum := 1
	if page == "" && offset > 0 { // first page for offset is 0
		pagenum = offset + 1
	} else if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if err != nil {
			return
		}
	}
	if pagenum <= 0 {
		err = errors.New("Page number is invalid")
		return
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

	if cat != "" {
		query := r.URL.Query()
		if cat == "5070" {
			query.Set("c", "3_")
		} else {
			c := nyaafeeds.ConvertToCat(cat)
			if c == "" {
				return
			}
			query.Set("c", c)
		}
		r.URL.RawQuery = query.Encode()
	}

	_, torrents, err = search.SearchByQueryNoCount(r, pagenum)

	return
}
