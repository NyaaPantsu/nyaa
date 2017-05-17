package metainfoFetcher;

import (
	"testing"

	"github.com/anacrolix/torrent"
	"github.com/NyaaPantsu/nyaa/model"
)

func TestInvalidHash(t *testing.T) {
	client, err := torrent.NewClient(nil)
	if err != nil {
		t.Skipf("Failed to create client, with err %v. Skipping.", err)
	}

	fetcher := &MetainfoFetcher{
		timeout: 5,
		torrentClient: client,
		results: make(chan Result, 1),
	}

	dbEntry := model.Torrent{
		Hash: "INVALID",
		Name: "Invalid",
	}

	op := NewFetchOperation(fetcher, dbEntry)
	fetcher.wg.Add(1)
	op.Start(fetcher.results)

	var res Result
	select {
	case res = <-fetcher.results:
		break
	default:
		t.Fatal("No result in channel, should have one")
	}

	if res.err == nil {
		t.Fatal("Got no error, should have got invalid magnet")
	}

	t.Logf("Got error %s, shouldn't be timeout", res.err)
}

