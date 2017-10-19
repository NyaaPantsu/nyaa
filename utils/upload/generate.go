package upload

import (
	"errors"
	"fmt"
	"os"

	"github.com/NyaaPantsu/nyaa/config"
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
		fmt.Println("got info")
		mi := t.Metainfo()
		fmt.Println("meta info")
		t.Drop()
		fmt.Println("drop")
		f, err := os.Create(fmt.Sprintf("%s%c%s.torrent", config.Get().Torrents.FileStorage, os.PathSeparator, t.InfoHash().String()))
		fmt.Println("open file")
		if err != nil {
			log.Errorf("error creating torrent metainfo file: %s", err)
			return
		}
		fmt.Println("defer")
		defer f.Close()
		fmt.Println("bencode")
		err = bencode.NewEncoder(f).Encode(mi)
		if err != nil {
			log.Errorf("error writing torrent metainfo file: %s", err)
			return
		}
		fmt.Println("for loop")
		for k, m := range queue {
			if m == magnet {
				queue = append(queue[:k], queue[k+1:]...)
			}
		}
	}()
	return nil
}
