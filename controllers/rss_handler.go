package controllers

import (
	"errors"
	"html"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/users"
	"github.com/NyaaPantsu/nyaa/utils/categories"
	"github.com/NyaaPantsu/nyaa/utils/feeds"
	"github.com/NyaaPantsu/nyaa/utils/publicSettings"
	"github.com/NyaaPantsu/nyaa/utils/search"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
)

// RSSHandler : Controller for displaying rss feed, accepting common search arguments
func RSSHandler(c *gin.Context) {
	// We only get the basic variable for rss based on search param
	torrents, createdAsTime, title, err := getTorrentList(c)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
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
			Link:        config.WebAddress() + "/download/" + torrentJSON.Hash,
			Description: string(torrentJSON.Description),
			PubDate:     torrent.Date.Format(time.RFC822),
			GUID:        config.WebAddress() + "/view/" + strconv.FormatUint(uint64(torrentJSON.ID), 10),
			Enclosure: &nyaafeeds.RssEnclosure{
				URL:    config.WebAddress() + "/download/" + strings.TrimSpace(torrentJSON.Hash),
				Length: strconv.FormatUint(uint64(torrentJSON.Filesize), 10),
				Type:   "application/x-bittorrent",
			},
		}
	}
	// allow cross domain AJAX requests
	c.Header("Access-Control-Allow-Origin", "*")
	rss, rssErr := feeds.ToXML(feed)
	if rssErr != nil {
		c.AbortWithError(http.StatusInternalServerError, rssErr)
	}

	_, writeErr := c.Writer.Write([]byte(rss))
	if writeErr != nil {
		c.AbortWithError(http.StatusInternalServerError, writeErr)
	}
}

// RSSMagnetHandler : Controller for displaying rss feeds with magnet URL, accepting common search arguments
func RSSMagnetHandler(c *gin.Context) {
	// We only get the basic variable for rss based on search param
	torrents, createdAsTime, title, err := getTorrentList(c)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
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
			Link:        &nyaafeeds.RssMagnetLink{Text: string(torrentJSON.Magnet)},
			Description: string(torrentJSON.Description),
			PubDate:     torrent.Date.Format(time.RFC822),
			GUID:        config.WebAddress() + "/view/" + strconv.FormatUint(uint64(torrentJSON.ID), 10),
			Enclosure: &nyaafeeds.RssEnclosure{
				URL:    config.WebAddress() + "/download/" + strings.TrimSpace(torrentJSON.Hash),
				Length: strconv.FormatUint(uint64(torrentJSON.Filesize), 10),
				Type:   "application/x-bittorrent",
			},
		}
	}
	// allow cross domain AJAX requests
	c.Header("Access-Control-Allow-Origin", "*")
	rss, rssErr := feeds.ToXML(feed)
	if rssErr != nil {
		c.AbortWithError(http.StatusInternalServerError, rssErr)
	}

	_, writeErr := c.Writer.Write([]byte(rss))
	if writeErr != nil {
		c.AbortWithError(http.StatusInternalServerError, writeErr)
	}
}

// RSSEztvHandler : Controller for displaying rss feed, accepting common search arguments
func RSSEztvHandler(c *gin.Context) {
	// We only get the basic variable for rss based on search param
	torrents, createdAsTime, title, err := getTorrentList(c)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
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
			PubDate:     torrent.Date.Format(time.RFC822),
			GUID:        config.WebAddress() + "/view/" + strconv.FormatUint(uint64(torrentJSON.ID), 10),
			Enclosure: &nyaafeeds.RssEnclosure{
				URL:    config.WebAddress() + "/download/" + strings.TrimSpace(torrentJSON.Hash),
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
	c.Header("Access-Control-Allow-Origin", "*")
	rss, rssErr := feeds.ToXML(feed)
	if rssErr != nil {
		c.AbortWithError(http.StatusInternalServerError, rssErr)
	}

	_, writeErr := c.Writer.Write([]byte(rss))
	if writeErr != nil {
		c.AbortWithError(http.StatusInternalServerError, writeErr)
	}
}

// RSSTorznabHandler : Controller for displaying rss feed, accepting common search arguments
func RSSTorznabHandler(c *gin.Context) {
	t := c.Query("t")
	rss := ""
	title := "Nyaa Pantsu"
	if config.IsSukebei() {
		title = "Sukebei Pantsu"
	}
	if t == "caps" {
		T := publicSettings.GetTfuncFromRequest(c)
		cats := categories.GetSelect(true, true)
		var categories []*nyaafeeds.RssCategoryTorznab
		categories = append(categories, &nyaafeeds.RssCategoryTorznab{
			ID:          "5070",
			Name:        "Anime",
			Description: "Anime",
		})

		last := 0
		for _, v := range cats {
			if len(v.ID) <= 2 {
				categories = append(categories, &nyaafeeds.RssCategoryTorznab{
					ID:   nyaafeeds.ConvertFromCat(v.ID),
					Name: string(T(v.Name)),
				})
				last++
			} else {
				categories[last].Subcat = append(categories[last].Subcat, &nyaafeeds.RssSubCat{
					ID:   nyaafeeds.ConvertFromCat(v.ID),
					Name: string(T(v.Name)),
				})
			}
		}
		feed := &nyaafeeds.RssCaps{
			Server: &nyaafeeds.RssServer{
				Version:   "1.0",
				Title:     title,
				Strapline: "...",
				Email:     config.Get().Email.From,
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
			c.AbortWithError(http.StatusInternalServerError, rssErr)
		}
	} else {
		// We only get the basic variable for rss based on search param
		torrents, createdAsTime, title, err := getTorrentList(c)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
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
				PubDate:     torrent.Date.Format(time.RFC822),
				GUID:        config.WebAddress() + "/view/" + strconv.FormatUint(uint64(torrentJSON.ID), 10),
				Enclosure: &nyaafeeds.RssEnclosure{
					URL:    config.WebAddress() + "/download/" + strings.TrimSpace(torrentJSON.Hash),
					Length: strconv.FormatUint(uint64(torrentJSON.Filesize), 10),
					Type:   "application/x-bittorrent",
				},
			}
			torznab := []*nyaafeeds.RssTorznab{}
			if torrentJSON.Filesize > 0 {
				torznab = append(torznab, &nyaafeeds.RssTorznab{
					Name:  "size",
					Value: strconv.FormatUint(uint64(torrentJSON.Filesize), 10),
				})
			}
			if filesNumber != "" {
				torznab = append(torznab, &nyaafeeds.RssTorznab{
					Name:  "files",
					Value: filesNumber,
				})
			}
			torznab = append(torznab, &nyaafeeds.RssTorznab{
				Name:  "grabs",
				Value: strconv.Itoa(int(torrentJSON.Completed)),
			})
			if seeders != "" {
				torznab = append(torznab, &nyaafeeds.RssTorznab{
					Name:  "seeders",
					Value: seeders,
				})
			}
			if leechers != "" {
				torznab = append(torznab, &nyaafeeds.RssTorznab{
					Name:  "leechers",
					Value: leechers,
				})
			}
			if torrentJSON.Hash != "" {
				torznab = append(torznab, &nyaafeeds.RssTorznab{
					Name:  "infohash",
					Value: torrentJSON.Hash,
				})
			}
			if torrentJSON.Magnet != "" {
				torznab = append(torznab, &nyaafeeds.RssTorznab{
					Name:  "magneturl",
					Value: string(torrentJSON.Magnet),
				})
			}
			if len(torznab) > 0 {
				feed.Items[i].Torznab = torznab
			}
		}
		var rssErr error
		rss, rssErr = feeds.ToXML(feed)
		if rssErr != nil {
			c.AbortWithError(http.StatusInternalServerError, rssErr)
		}
	}
	// allow cross domain AJAX requests
	c.Header("Access-Control-Allow-Origin", "*")

	_, writeErr := c.Writer.Write([]byte(rss))
	if writeErr != nil {
		c.AbortWithError(http.StatusInternalServerError, writeErr)
	}
}

func getTorrentList(c *gin.Context) (torrents []models.Torrent, createdAsTime time.Time, title string, err error) {
	page := c.Param("page")
	userID := c.Param("id")
	cat := c.Query("cat")
	offset := 0
	if c.Query("offset") != "" {
		offset, err = strconv.Atoi(html.EscapeString(c.Query("offset")))
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

		_, _, err = users.FindForAdmin(uint(userIDnum))
		if err != nil {
			return
		}
		// Set the user ID on the request, so that SearchByQuery finds it.
		query := c.Request.URL.Query()
		query.Set("userID", userID)
		c.Request.URL.RawQuery = query.Encode()
	}

	if cat != "" {
		query := c.Request.URL.Query()
		catConv := nyaafeeds.ConvertToCat(cat)
		if catConv == "" {
			return
		}
		query.Set("c", catConv)
		c.Request.URL.RawQuery = query.Encode()
	}

	_, torrents, err = search.ByQueryNoCount(c, pagenum)

	return
}
