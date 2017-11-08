package torrentController

import (
	"html/template"
	"encoding/hex"
	"net/http"
	"strings"
	"strconv"

	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/templates"
	"github.com/NyaaPantsu/nyaa/utils/format"
	"github.com/NyaaPantsu/nyaa/utils/filelist"
	"github.com/Stephen304/goscrape"
	"github.com/gin-gonic/gin"
)

func GetFilesHandler(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)
	torrent, err := torrents.FindByID(uint(id))	

	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	
	
	if len(torrent.FileList) == 0 {
		var blankScrape models.Scrape
		ScrapeFiles(format.InfoHashToMagnet(strings.TrimSpace(torrent.Hash), torrent.Name, GetTorrentTrackers(torrent)...), torrent, blankScrape, true)
	}
	
	folder := filelist.FileListToFolder(torrent.FileList, "root")
	templates.TorrentFileList(c, torrent.ToJSON(), folder)
}

// ScrapeFiles : Scrape torrent files
func ScrapeFiles(magnet string, torrent *models.Torrent, currentStats models.Scrape, statsExists bool) (error, []FileJSON) {
	if client == nil {
		err := initClient()
		if err != nil {
			return err, []FileJSON{}
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

// FileJSON for file model in json, 
type FileJSON struct {
	Path       string         `json:"path"`
	Filesize   template.HTML  `json:"filesize"`
}

func fileSize(filesize int64) template.HTML {
	return template.HTML(format.FileSize(filesize))
}