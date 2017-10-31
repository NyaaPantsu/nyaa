package upload

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/utils/format"
	"github.com/NyaaPantsu/nyaa/utils/log"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/bencode"
)

var queue []string
var client *torrent.Client

func initClient() error {
	cl, err := torrent.NewClient(nil)
	if err != nil {
		log.Errorf("error creating client: %s", err)
		return err
	}
	client = cl
	return nil
}

// GenerateTorrent generates a torrent file in the specified directory in config.yml from a magnet URI
func GenerateTorrent(magnet string) error {
	if client == nil {
		err := initClient()
		if err != nil {
			return err
		}
	}
	if magnet == "" || len(config.Get().Torrents.FileStorage) == 0 {
		return errConfig
	}
	if len(queue) > 0 {
		for _, m := range queue {
			if m == magnet {
				return errPending
			}
		}
	}
	queue = append(queue, magnet)

	t, err := client.AddMagnet(magnet)
	if err != nil {
		log.Errorf("error adding magnet to client: %s", err)
		return err
	}
	<-t.GotInfo()
	mi := t.Metainfo()
	t.Drop()
	file := fmt.Sprintf("%s%c%s.torrent", config.Get().Torrents.FileStorage, os.PathSeparator, t.InfoHash().String())
	f, err := os.Create(file)
	if err != nil {
		log.Errorf("error creating torrent metainfo file: %s", err)
		return err
	}
	defer f.Close()
	err = bencode.NewEncoder(f).Encode(mi)
	if err != nil {
		log.Errorf("error writing torrent metainfo file: %s", err)
		return err
	}
	for k, m := range queue {
		if m == magnet {
			queue = append(queue[:k], queue[k+1:]...)
		}
	}
	log.Infof("New torrent file generated in: %s", file)

	return nil
}

// GotFile will check if a torrent file exists and if not, try to generate it
func GotFile(torrent *models.Torrent) error {
	var trackers []string
	if torrent.Trackers == "" {
		trackers = config.Get().Torrents.Trackers.Default
	} else {
		trackers = torrent.GetTrackersArray()
	}
	// We generate a new magnet link with all trackers (ours + ones from uploader)
	magnet := format.InfoHashToMagnet(strings.TrimSpace(torrent.Hash), torrent.Name, trackers...)
	//Check if file exists and open
	_, err := os.Open(torrent.GetPath())
	if err != nil {
		err := GenerateTorrent(magnet)
		if err != errPending && err != nil {
			return err
		}
		if err == errPending {
			for {
				//Check again if file exists and open
				_, err := os.Open(torrent.GetPath())
				if err == nil {
					break
				}
				// We check every 10 seconds
				time.Sleep(10000)
			}
		}
	}
	return nil
}
