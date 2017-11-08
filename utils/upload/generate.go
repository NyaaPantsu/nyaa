package upload

import (
	"strconv"
	"errors"
	"fmt"
	"os"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/utils/log"
	"github.com/anacrolix/dht"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/bencode"
)

var queue []string
var client *torrent.Client

func initClient() error {
	clientConfig := torrent.Config{
		DHTConfig: dht.ServerConfig{
			StartingNodes: dht.GlobalBootstrapAddrs,
		},
		ListenAddr: ":" + strconv.Itoa(config.Get().Torrents.GenerationClientPort),
	}
	cl, err := torrent.NewClient(&clientConfig)
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
		return errors.New("Magnet Empty or FileStorage not configured")
	}
	if len(queue) > 0 {
		for _, m := range queue {
			if m == magnet {
				return errors.New("Magnet being generated already")
			}
		}
	}
	queue = append(queue, magnet)

	t, err := client.AddMagnet(magnet)
	if err != nil {
		log.Errorf("error adding magnet to client: %s", err)
		return err
	}
	go func() {
		<-t.GotInfo()
		mi := t.Metainfo()
		t.Drop()
		file := fmt.Sprintf("%s%c%s.torrent", config.Get().Torrents.FileStorage, os.PathSeparator, t.InfoHash().String())
		f, err := os.Create(file)
		if err != nil {
			log.Errorf("error creating torrent metainfo file: %s", err)
			return
		}
		defer f.Close()
		err = bencode.NewEncoder(f).Encode(mi)
		if err != nil {
			log.Errorf("error writing torrent metainfo file: %s", err)
			return
		}
		for k, m := range queue {
			if m == magnet {
				queue = append(queue[:k], queue[k+1:]...)
			}
		}
		log.Infof("New torrent file generated in: %s", file)
	}()
	return nil
}
