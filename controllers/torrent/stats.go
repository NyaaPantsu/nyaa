package torrentController

import (
	"text/template"
	"strconv"
	"strings"
	"net/url"
	"time"
	"fmt"

	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/Stephen304/goscrape"
	"github.com/gin-gonic/gin"
)

// ViewHeadHandler : Controller for getting torrent stats
func GetStatsHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		return
	}

	torrent, err := torrents.FindRawByID(uint(id))

	if err != nil {
		return
	}
	
	var Trackers []string
	for _, line := range strings.Split(torrent.Trackers[3:], "&tr=") {
		tracker, err := url.QueryUnescape(line)
		if tracker[:6] == "udp://" {
			Trackers = append(Trackers, tracker)
		}
	}	

	scraper := goscrape.NewBulk(Trackers)
	
	stats := scraper.ScrapeBulk([]string{
	  torrent.Hash,
	})[0]
	
	emptyStats := goscrape.Result{stats.Btih, 0, 0, 0} 
	
	if stats == emptyStats {
		stats.Seeders = -1
		//If we put seeders on -1, the script instantly knows the fetching did not give any result, avoiding having to check all three stats below and in view.jet.html's javascript
	}
	
	t, err := template.New("foo").Parse(fmt.Sprintf(`{{define "stats"}}{ "seeders": [%d], "leechers": [%d], "downloads": [%d] }{{end}}`, stats.Seeders, stats.Leechers, stats.Completed))
	t.ExecuteTemplate(c.Writer, "stats", "")
	
	if stats.Seeders != -1 {
		var tmp models.Scrape
		if models.ORM.Where("torrent_id = ?", id).Find(&tmp).RecordNotFound() {
			torrent.Scrape = torrent.Scrape.Create(uint(id), uint32(stats.Seeders), uint32(stats.Leechers), uint32(stats.Completed), time.Now())
		} else {
			torrent.Scrape = &models.Scrape{uint(id), uint32(stats.Seeders), uint32(stats.Leechers), uint32(stats.Completed), time.Now()}
			torrent.Scrape.Update(false)
		}
	}
	
	return
}
