package metainfoFetcher;

import (
	"github.com/anacrolix/torrent/metainfo"
	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/util"
	"errors"
	"time"
	"strings"
)

type FetchOperation struct {
	fetcher  *MetainfoFetcher
	torrent  model.Torrent
	done     chan int
}

type Result struct {
	operation *FetchOperation
	err       error
	info      *metainfo.Info
}

func NewFetchOperation(fetcher *MetainfoFetcher, dbEntry model.Torrent) (op *FetchOperation) {
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
	defer downloadingTorrent.Drop()

	timeoutTimer := time.NewTicker(time.Second * time.Duration(op.fetcher.timeout))
	defer timeoutTimer.Stop()
	select {
	case <-downloadingTorrent.GotInfo():
		out <- Result{op, nil, downloadingTorrent.Info()}
		break
	case <-timeoutTimer.C:
		out <- Result{op, errors.New("Timeout"), nil}
		break
	case <-op.done:
		break
	}
}


