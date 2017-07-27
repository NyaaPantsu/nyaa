package feedController

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/utils/categories"
	"github.com/NyaaPantsu/nyaa/utils/feeds"
	"github.com/NyaaPantsu/nyaa/utils/publicSettings"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
)

// RSSTorznabHandler : Controller for displaying rss feed, accepting common search arguments
func RSSTorznabHandler(c *gin.Context) {
	t := c.Query("t")
	var rss string
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
