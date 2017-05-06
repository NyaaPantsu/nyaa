package metainfo

import (
	"github.com/zeebo/bencode"
	"os"
	"strings"
	"testing"
)

func TestLoadTorrent(t *testing.T) {
	f, err := os.Open("test.torrent")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	tf := new(TorrentFile)
	dec := bencode.NewDecoder(f)
	err = dec.Decode(tf)
	if err != nil {
		t.Error(err)
	}

	if strings.ToUpper(tf.Infohash().Hex()) != "6BCDC07177EC43658C1B4D5450640059663A5214" {
		t.Error(tf.Infohash().Hex())
	}
	// TODO: check members
}
