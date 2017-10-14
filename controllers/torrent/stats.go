package torrentController

import (
	"strconv"
	"strings"
	"net/url"
	"time"

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
	
	c.JSON(200, gin.H{
 		"seeders": stats.Seeders,
 		"leechers": stats.Leechers,
 		"downloads": stats.Completed,
 	})
	
	if stats.Seeders == -1 {
		stats.Seeders = 0
	}
	
	var CurrentData models.Scrape
	if models.ORM.Where("torrent_id = ?", id).Find(&CurrentData).RecordNotFound() {
		torrent.Scrape = torrent.Scrape.Create(uint(id), uint32(stats.Seeders), uint32(stats.Leechers), uint32(stats.Completed), time.Now())
		//Create entry in the DB because none exist
	} else {
		//Entry in the DB already exists, simply update it
		if (CurrentData.Seeders == 0 && CurrentData.Leechers == 0 && CurrentData.Completed == 0) || (stats.Seeders != 0 && stats.Leechers != 0 && stats.Completed != 0 ) {
			torrent.Scrape = &models.Scrape{uint(id), uint32(stats.Seeders), uint32(stats.Leechers), uint32(stats.Completed), time.Now()}
		} else {
			torrent.Scrape = &models.Scrape{uint(id), uint32(CurrentData.Seeders), uint32(CurrentData.Leechers), uint32(CurrentData.Completed), time.Now()}
		}
		//Only overwrite stats if the old one are Unknown OR if the current ones are not unknown, preventing good stats from being turned into unknown own but allowing good stats to be updated to more reliable ones
		torrent.Scrape.Update(false)
		
	}
	
	return
}
