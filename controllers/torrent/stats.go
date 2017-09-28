package torrentController

import (
	"net/http"
	"strconv"
	"fmt"

	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/gin-gonic/gin"
	"text/template"
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
	
	seeders := 3
	leechers := 3
	downloads := 3
  	//TODO: fetch torrent stats and store it in the above variables 
	//if unknown put all three on -1
	
	t, err := template.New("foo").Parse(fmt.Sprintf(`{{define "stats"}}{ "seeders":[%d], "leechers": [%d], "downloads": [%d] }{{end}}`, seeders, leechers, downloads))
	err = t.ExecuteTemplate(c.Writer, "stats", "")
	
	return
}

