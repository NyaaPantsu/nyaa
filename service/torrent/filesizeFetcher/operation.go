package filesizeFetcher;

import (
	"github.com/anacrolix/torrent/metainfo"
	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/util"
	"errors"
	"time"
	"strings"
)

type FetchOperation struct {
	fetcher  *FilesizeFetcher
	torrent  model.Torrent
	done     chan int
}

type Result struct {
	operation *FetchOperation
	err       error
	info      *metainfo.Info
}

func NewFetchOperation(fetcher *FilesizeFetcher, dbEntry model.Torrent) (op *FetchOperation) {
	op = &FetchOperation{
		fetcher:  fetcher,
		torrent:  dbEntry,
		done:     make(chan int, 1),
	}
	return
}

// Should be started from a goroutine somewhere
func (op *FetchOperation) Start(out chan Result) {
	defer op.fetcher.wg.Done()

	magnet := util.InfoHashToMagnet(strings.TrimSpace(op.torrent.Hash), op.torrent.Name, config.Trackers...)
	downloadingTorrent, err := op.fetcher.torrentClient.AddMagnet(magnet)
	if err != nil {
		out <- Result{op, err, nil}
		return
	}

	timeoutTicker := time.NewTicker(time.Second * time.Duration(op.fetcher.timeout))
	select {
	case <-downloadingTorrent.GotInfo():
		downloadingTorrent.Drop()
		out <- Result{op, nil, downloadingTorrent.Info()}
		break
	case <-timeoutTicker.C:
		downloadingTorrent.Drop()
		out <- Result{op, errors.New("Timeout"), nil}
		break
	case <-op.done:
		downloadingTorrent.Drop()
		break
	}
}

