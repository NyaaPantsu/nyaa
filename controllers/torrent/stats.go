package torrentController

import (
	"path/filepath"
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
	"github.com/bradfitz/slice"
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
	
	updateTorrent, err := torrents.FindByID(uint(id))
	
	if err != nil {
		return
	}
	
	var CurrentData models.Scrape
	statsExists := !(models.ORM.Where("torrent_id = ?", id).Find(&CurrentData).RecordNotFound())

	if statsExists && c.Request.URL.Query()["files"] == nil {
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
	
	Trackers := GetTorrentTrackers(updateTorrent)
	
	var stats goscrape.Result
	var torrentFiles []FileJSON
	
	if c.Request.URL.Query()["files"] != nil {
		if len(updateTorrent.FileList) > 0 {
			return
		}
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
		"totalsize": fileSize(updateTorrent.Filesize),
 	})
	
	return
}

// UpdateTorrentStats : Update stats & filelist if files are specified, otherwise just stats
func UpdateTorrentStats(torrent *models.Torrent, stats goscrape.Result, currentStats models.Scrape, Files []torrent.File, statsExists bool) (JSONFilelist []FileJSON) {
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
		files, err := torrent.CreateFileList(Files)
		
		if err != nil {
			return
		}
		
		JSONFilelist = make([]FileJSON, 0, len(files))
		for _, f := range files {
			JSONFilelist = append(JSONFilelist, FileJSON{
				Path:     filepath.Join(f.Path()...),
				Filesize: fileSize(f.Filesize),
			})
		}

		// Sort file list by lowercase filename
		slice.Sort(JSONFilelist, func(i, j int) bool {
			return strings.ToLower(JSONFilelist[i].Path) < strings.ToLower(JSONFilelist[j].Path)
		})
	}
	
	return
}

// GetTorrentTrackers : Get the torrent trackers and add the default ones if they are missing
func GetTorrentTrackers(torrent *models.Torrent) []string {
	var Trackers []string
	if len(torrent.Trackers) > 3 {
		for _, line := range strings.Split(torrent.Trackers[3:], "&tr=") {
			tracker, error := url.QueryUnescape(line)
			if error == nil && strings.HasPrefix(tracker, "udp") {
				Trackers = append(Trackers, tracker)
			}
			//Cannot scrape from http trackers only keep UDP ones
		}
	}
	
	for _, tracker := range config.Get().Torrents.Trackers.Default {
		if !contains(Trackers, tracker) && strings.HasPrefix(tracker, "udp") {
			Trackers = append(Trackers, line)
		}
	}
	return Trackers
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
