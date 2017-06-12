package metainfoFetcher

import (
	"errors"
	"strings"
	"time"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/util"
	"github.com/anacrolix/torrent/metainfo"
)

// FetchOperation struct
type FetchOperation struct {
	fetcher *MetainfoFetcher
	torrent model.Torrent
	done    chan struct{}
}

// Result struct
type Result struct {
	operation *FetchOperation
	err       error
	info      *metainfo.Info
}

// NewFetchOperation : Creates a new fetchoperation
func NewFetchOperation(fetcher *MetainfoFetcher, dbEntry model.Torrent) (op *FetchOperation) {
	op = &FetchOperation{
		fetcher: fetcher,
		torrent: dbEntry,
		done:    make(chan struct{}, 1),
	}
	return
}

// Start : Should be started from a goroutine somewhere
func (op *FetchOperation) Start(out chan Result) {
	defer op.fetcher.wg.Done()

	magnet := util.InfoHashToMagnet(strings.TrimSpace(op.torrent.Hash), op.torrent.Name, config.Conf.Torrents.Trackers.Default...)
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
	case <-timeoutTimer.C:
		out <- Result{op, errors.New("Timeout"), nil}
	case <-op.done:
	}
}
