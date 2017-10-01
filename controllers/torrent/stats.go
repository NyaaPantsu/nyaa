package torrentController

import (
	"text/template"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"regexp"
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
	
	t, err := template.New("foo").Parse(fmt.Sprintf(`{{define "stats"}}{ "seeders":[%d], "leechers": [%d], "downloads": [%d] }{{end}}`, seeders, leechers, downloads))
	err = t.ExecuteTemplate(c.Writer, "stats", "")
	
	return
}

//Copypasted from Nyaapantsu/scrapers/stats.go

type Stats struct {
	Btih      string
	Seeders   int
	Leechers  int
	Completed int
}

type TStruct struct {
	Peers    Stats
	Trackers []string
	//Files    []metainfo.FileInfo
	Files  []torrent.File
	Magnet string
}

var validChar = regexp.MustCompile(`[a-zA-Z0-9_\-\~]`)

//byteToStr : Function to encode any non-tolerated characters in a tracker request to hex
func byteToStr(arr []byte) (str string) {
	for _, b := range arr {
		c := string(b)
		if !validChar.MatchString(c) {
			dst := make([]byte, hex.EncodedLen(len(c)))
			hex.Encode(dst, []byte(c))
			c = string(dst)
		}
		str += c
	}
	return
}

func udpScrape(trackers []string, hash string, chFin chan<- bool, torr *TStruct) {
	udpscrape := goscrape.NewBulk(trackers)
	results := udpscrape.ScrapeBulk([]string{hash})
	if results[0].Btih != "0" {
		torr.Peers = Stats(results[0])
	} else {
		fmt.Println("Bad results: ", results[0])
		udpScrape(trackers, hash, chFin, torr)
	}
	chFin <- true
}

func fileScrape(client *torrent.Client, torr *TStruct, chFin chan<- bool) {
	t, _ := client.AddMagnet(torr.Magnet)
	<-t.GotInfo()
	infoHash := t.InfoHash()
	dst := make([]byte, hex.EncodedLen(len(t.InfoHash())))
	hex.Encode(dst, infoHash[:])
	var UDP []string
	var HTTP []string
	torr.Trackers = t.Metainfo().AnnounceList[0]
	for _, tracker := range torr.Trackers {
		if strings.HasPrefix(tracker, "http") {
			HTTP = append(HTTP, tracker)
		} else if strings.HasPrefix(tracker, "udp") {
			UDP = append(UDP, tracker)
		}
	}
	if len(UDP) != 0 {
		go udpScrape(UDP, string(dst), chFin, torr)
	}
	//	metaInfo := t.Info()
	//	torr.Files = t.UpvertedFiles()
	torr.Files = t.Files()
	t.Drop()
	chFin <- true
}

//I'm sure there's a less sloppy way to do this, but let's call this an "alpha" version
func injectStats(t *Torrent, torr *TStruct) {
	fmt.Println("Injecting stats!")
	t.Seeders = torr.Peers.Seeders
	t.Leechers = torr.Peers.Leechers
	t.Completed = torr.Peers.Completed
	fmt.Println("Printing files for", t.Magnet)
	fmt.Println(torr.Files)
	//Current filelist struct uses a string array, not sure how to convert that
	//t.FileList = torr.Files
}

func grabEverything(client *torrent.Client, torr TStruct, t Torrent, chOut chan<- Torrent) {
	chFin := make(chan bool)
	go fileScrape(client, &torr, chFin)
	for i := 0; i < 2; {
		select {
		case <-chFin:
			i++
		}
	}
	injectStats(&t, &torr)
	chOut <- t
}

func statWorker(chIn <-chan Torrent, chOut chan<- Torrent) {
	client, _ := torrent.NewClient(nil)
	for t := range chIn {
		torr := TStruct{}
		torr.Magnet = t.Magnet
		go grabEverything(client, torr, t, chOut)
	}
}
