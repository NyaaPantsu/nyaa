package torrentController

import (
	"text/template"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"fmt"

	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/Stephen304/goscrape"
	"github.com/anacrolix/torrent"
	"github.com/gin-gonic/gin"
)

// ViewHeadHandler : Controller for getting torrent stats
func GetStatsHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		return
	}

	_, err = torrents.FindRawByID(uint(id))

	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	
	seeders := -1
	leechers := -1
	downloads := -1
  	//TODO: fetch torrent stats and store it in the above variables 
	//if unknown let all three on -1
	
	scraper := NewBulk([]string{
	  "udp://tracker.doko.moe:6969",
	  "udp://tracker.leechers-paradise.org:6969"})
	
	result := scraper.ScrapeBulk([]string{
	  "d791154b92068ee15dd0fd1f71686a0fc79b0aef",
	  "d7cef3aa4a5f71963eee4c481f8d48490febee18",
	})
	
	
	
	t, err := template.New("foo").Parse(fmt.Sprintf(`{{define "stats"}}{ "seeders":[%d], "leechers": [%d], "downloads": [%d] }{{end}}`, seeders, leechers, downloads))
	t, err := template.New("foo").Parse(fmt.Sprintf(`{{define "stats"}}{ "seeders":[%s], "leechers": [%d], "downloads": [%d] }{{end}}`, result[0] , leechers, downloads))
	err = t.ExecuteTemplate(c.Writer, "stats", "")
	
	return
}
