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
		//Starts at character 3 because the three first characters are always "tr=" so we need to dismiss them
		tracker, error := url.QueryUnescape(line)
		if error == nil && tracker[:6] == "udp://" {
			Trackers = append(Trackers, tracker)
		}
		//Cannot scrape from http trackers so don't put them in the array
	}	

	stats := goscrape.Single(Trackers, []string{
	  torrent.Hash,
	})[0]
	//Single() returns an array which contain results for each torrent Hash it is fed, since we only feed him one we want to directly access the results
	
	//If we put seeders on -1, the script instantly knows the fetching did not give any result, avoiding having to check all three stats below and in view.jet.html's javascript
	if stats.Seeders == 0 && stats.Leechers == 0 && stats.Completed == 0  {
		stats.Seeders = -1
	}
	
	t, err := template.New("foo").Parse(fmt.Sprintf(`{{define "stats"}}{ "seeders": [%d], "leechers": [%d], "downloads": [%d] }{{end}}`, stats.Seeders, stats.Leechers, stats.Completed))
	t.ExecuteTemplate(c.Writer, "stats", "")
	//No idea how to output JSON properly
	
	//We don't want to do useless DB queries if the stats are empty, and we don't want to overwrite good stats with empty ones
	if stats.Seeders != -1 {
		var tmp models.Scrape
		if models.ORM.Where("torrent_id = ?", id).Find(&tmp).RecordNotFound() {
			torrent.Scrape = torrent.Scrape.Create(uint(id), uint32(stats.Seeders), uint32(stats.Leechers), uint32(stats.Completed), time.Now())
			//Create entry in the DB because none exist
		} else {
			torrent.Scrape = &models.Scrape{uint(id), uint32(stats.Seeders), uint32(stats.Leechers), uint32(stats.Completed), time.Now()}
			torrent.Scrape.Update(false)
			//Entry in the DB already exists, simply update it
		}
	}
	
	return
}
