package feedController

import (
	"net/http"
	"strconv"
	"strings"
	"time"
	"fmt"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/utils/feeds"

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
		PubDate: formatRSSDate(createdAsTime),
	}
	feed.Items = make([]*nyaafeeds.RssItem, len(torrents))

	for i, torrent := range torrents {
		torrentJSON := torrent.ToJSON()
		feed.Items[i] = &nyaafeeds.RssItem{
			Title:       torrentJSON.Name,
			Link:        config.WebAddress() + "/download/" + torrentJSON.Hash,
			Description: string(torrentJSON.Description),
			PubDate:     formatRSSDate(torrent.Date),
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

//Return date in an RFC 2822 format, the official one for RSS2
func formatRSSDate(Date time.Time) string {
	Date = Date.UTC()
	return fmt.Sprintf("%.3s, %.2d %.3s %d %.2d:%.2d:%.2d +0000", Date.Weekday(), Date.Day(), Date.Month(), Date.Year(), Date.Hour(), Date.Minute(), Date.Second())
}
