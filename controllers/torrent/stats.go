package torrentController

import (
	"text/template"
	"net/http"
	"strconv"
	"fmt"

	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/Stephen304/goscrape"
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
	  "udp://tracker.example.com:80",
	  "udp://tracker.example2.com:80"})
	
	t, err := template.New("foo").Parse(fmt.Sprintf(`{{define "stats"}}{ "seeders":[%d], "leechers": [%d], "downloads": [%d] }{{end}}`, seeders, leechers, downloads))
	err = t.ExecuteTemplate(c.Writer, "stats", "")
	
	return
}
