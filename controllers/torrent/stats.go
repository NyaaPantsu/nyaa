package torrentController

import (
	"strconv"
	"strings"
	"net/url"
	"time"

	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/config"
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
	
	var CurrentData models.Scrape
	statsExists := !(models.ORM.Where("torrent_id = ?", id).Find(&CurrentData).RecordNotFound())

	if statsExists {
		//Stats already exist, we check if the torrent stats have been scraped already very recently and if so, we stop there to avoid abuse of the /stats/:id route
		if (CurrentData.Seeders == 0 && CurrentData.Leechers == 0 && CurrentData.Completed == 0)  && time.Since(CurrentData.LastScrape).Minutes() <= config.Get().Scrape.MaxStatScrapingFrequencyUnknown {
			//Unknown stats but has been scraped less than X minutes ago (X being the limit set in the config file)
			return
		}
		if (CurrentData.Seeders != 0 || CurrentData.Leechers != 0 || CurrentData.Completed != 0) && time.Since(CurrentData.LastScrape).Minutes() <= config.Get().Scrape.MaxStatScrapingFrequency  {
			//Known stats but has been scraped less than X minutes ago (X being the limit set in the config file)
			return
		}
	}
	
	var Trackers []string
	if len(torrent.Trackers) > 3 {
		for _, line := range strings.Split(torrent.Trackers[3:], "&tr=") {
			tracker, error := url.QueryUnescape(line)
			if error == nil && strings.HasPrefix(tracker, "udp") {
				Trackers = append(Trackers, tracker)
			}
			//Cannot scrape from http trackers so don't put them in the array
		}
	}
	
	for _, tracker := range config.Get().Torrents.Trackers.Default {
		if !contains(Trackers, tracker) && strings.HasPrefix(tracker, "udp")  {
			Trackers = append(Trackers, tracker)
		}
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
	
	if !statsExists {
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

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}
