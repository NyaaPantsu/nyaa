package torrentController

import (
	"encoding/hex"
	"strconv"
	"strings"
	"net/url"
	"time"

	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/utils/log"
	"github.com/NyaaPantsu/nyaa/utils/format"
	"github.com/Stephen304/goscrape"
	"github.com/gin-gonic/gin"
	
	"github.com/anacrolix/dht"
	"github.com/anacrolix/torrent"
)

var client *torrent.Client

func initClient() error {
	clientConfig := torrent.Config{
		DHTConfig: dht.ServerConfig{
			StartingNodes: dht.GlobalBootstrapAddrs,
		},
		ListenAddr: ":5977",
	}
	cl, err := torrent.NewClient(&clientConfig)
	if err != nil {
		log.Errorf("error creating client: %s", err)
		return err
	}
	client = cl
	return nil
}

// ViewHeadHandler : Controller for getting torrent stats
func GetStatsHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		return
	}
	
	updateTorrent, err := torrents.FindRawByID(uint(id))

	if err != nil {
		return
	}
	
	var CurrentData models.Scrape
	statsExists := !(models.ORM.Where("torrent_id = ?", id).Find(&CurrentData).RecordNotFound())

	if statsExists {
		//Stats already exist, we check if the torrent stats have been scraped already very recently and if so, we stop there to avoid abuse of the /stats/:id route
		if isEmptyScrape(CurrentData)  && time.Since(CurrentData.LastScrape).Minutes() <= config.Get().Scrape.MaxStatScrapingFrequencyUnknown {
			//Unknown stats but has been scraped less than X minutes ago (X being the limit set in the config file)
			return
		}
		if !isEmptyScrape(CurrentData) && time.Since(CurrentData.LastScrape).Minutes() <= config.Get().Scrape.MaxStatScrapingFrequency  {
			//Known stats but has been scraped less than X minutes ago (X being the limit set in the config file)
			return
		}
	}
	
	var Trackers []string
	if len(updateTorrent.Trackers) > 3 {
		for _, line := range strings.Split(updateTorrent.Trackers[3:], "&tr=") {
			tracker, error := url.QueryUnescape(line)
			if error == nil && strings.HasPrefix(tracker, "udp") {
				Trackers = append(Trackers, tracker)
			}
			//Cannot scrape from http trackers only keep UDP ones
		}
	}
	
	for _, line := range config.Get().Torrents.Trackers.Default {
		if !contains(Trackers, line) {
			Trackers = append(Trackers, line)
		}
	}
	
	var stats goscrape.Result
	var torrentFiles []models.FileJSON
	
	if c.Request.URL.Query()["files"] != nil {
		err, torrentFiles = ScrapeFiles(format.InfoHashToMagnet(strings.TrimSpace(updateTorrent.Hash), updateTorrent.Name, Trackers...), updateTorrent, CurrentData, statsExists)
		if err != nil {
			return
		}
	} else {
		//Single() returns an array which contain results for each torrent Hash it is fed, since we only feed him one we want to directly access the results
		stats = goscrape.Single(Trackers, []string{
		  updateTorrent.Hash,
		})[0]
		UpdateTorrentStats(updateTorrent, stats, CurrentData, []torrent.File{}, statsExists)
	}
	
	
	//If we put seeders on -1, the script instantly knows the fetching did not give any result, avoiding having to check all three stats below and in view.jet.html's javascript
	if isEmptyResult(stats) {
		stats.Seeders = -1
	}
	
	c.JSON(200, gin.H{
 		"seeders": stats.Seeders,
 		"leechers": stats.Leechers,
 		"downloads": stats.Completed,
		"filelist": torrentFiles,
 	})
	
	return
}

func ScrapeFiles(magnet string, torrent models.Torrent, currentStats models.Scrape, statsExists bool) (error, []models.FileJSON) {
	if client == nil {
		err := initClient()
		if err != nil {
			return err, []models.FileJSON{}
		}
	}
	
	t, _ := client.AddMagnet(magnet)
	<-t.GotInfo()
	
	infoHash := t.InfoHash()
	dst := make([]byte, hex.EncodedLen(len(t.InfoHash())))
	hex.Encode(dst, infoHash[:])
	
	var UDP []string
	
	for _, tracker := range t.Metainfo().AnnounceList[0] {
		if strings.HasPrefix(tracker, "udp") {
			UDP = append(UDP, tracker)
		}
	}
	var results goscrape.Result
	if len(UDP) != 0 {
		udpscrape := goscrape.NewBulk(UDP)
		results = udpscrape.ScrapeBulk([]string{torrent.Hash})[0]
	}
	t.Drop()
	return nil, UpdateTorrentStats(torrent, results, currentStats, t.Files(), statsExists)
}

// UpdateTorrentStats : Update stats & filelist if files are specified, otherwise just stats
func UpdateTorrentStats(torrent models.Torrent, stats goscrape.Result, currentStats models.Scrape, Files []torrent.File, statsExists bool) (JSONFilelist []models.FileJSON) {
	if stats.Seeders == -1 {
		stats.Seeders = 0
	}
	
	if !statsExists {
		torrent.Scrape = torrent.Scrape.Create(torrent.ID, uint32(stats.Seeders), uint32(stats.Leechers), uint32(stats.Completed), time.Now())
		//Create a stat entry in the DB because none exist
	} else {
		//Entry in the DB already exists, simply update it
		if isEmptyScrape(currentStats) || !isEmptyResult(stats) {
			torrent.Scrape = &models.Scrape{torrent.ID, uint32(stats.Seeders), uint32(stats.Leechers), uint32(stats.Completed), time.Now()}
		} else {
			torrent.Scrape = &models.Scrape{torrent.ID, uint32(currentStats.Seeders), uint32(currentStats.Leechers), uint32(currentStats.Completed), time.Now()}
		}
		//Only overwrite stats if the old one are Unknown OR if the new ones are not unknown, preventing good stats from being turned into unknown but allowing good stats to be updated to more reliable ones
		torrent.Scrape.Update(false)
	}
	
	if len(Files) > 0 {
		torrent.FileList = []models.File{}
		for i, file := range Files {
			torrent.FileList = append(torrent.FileList, models.File{uint(i), torrent.ID, file.DisplayPath(), file.Length()})
			JSONFilelist = append(JSONFilelist, models.FileJSON{file.DisplayPath(), file.Length()})
		}
		torrent.Update(true)
	}
	
	return
}

func isEmptyResult(stats goscrape.Result) bool {
	return stats.Seeders == 0 && stats.Leechers == 0 && stats.Completed == 0 
}

func isEmptyScrape(stats models.Scrape) bool {
	return stats.Seeders == 0 && stats.Leechers == 0 && stats.Completed == 0 
}

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}
