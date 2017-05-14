package filesizeFetcher;

import (
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/util/log"
	serviceBase "github.com/ewhal/nyaa/service"
	torrentService "github.com/ewhal/nyaa/service/torrent"
	"sync"
	"time"
)

type FilesizeFetcher struct {
	torrentClient    *torrent.Client
	results          chan Result
	queueSize        int
	timeout          int
	maxDays          int
	done             chan int
	queue            []*FetchOperation
	queueMutex       sync.Mutex
	failedOperations map[uint]struct{}
	wakeUp           *time.Ticker
	wg               sync.WaitGroup
}

func New(fetcherConfig *config.FilesizeFetcherConfig) (fetcher *FilesizeFetcher, err error) {
	client, err := torrent.NewClient(nil)
	fetcher = &FilesizeFetcher{
		torrentClient:    client,
		results:          make(chan Result, fetcherConfig.QueueSize),
		queueSize:        fetcherConfig.QueueSize,
		timeout:          fetcherConfig.Timeout,
		maxDays:          fetcherConfig.MaxDays,
		done:             make(chan int, 1),
		failedOperations: make(map[uint]struct{}),
		wakeUp:           time.NewTicker(time.Second * time.Duration(fetcherConfig.WakeUpInterval)),
	}

	return
}

func (fetcher *FilesizeFetcher) isFetchingOrFailed(t model.Torrent) bool {
	for _, op := range fetcher.queue {
		if op.torrent.ID == t.ID {
			return true
		}
	}

	_, ok := fetcher.failedOperations[t.ID]
	return ok
}

func (fetcher *FilesizeFetcher) addToQueue(op *FetchOperation) bool {
	fetcher.queueMutex.Lock()
	defer fetcher.queueMutex.Unlock()

	if len(fetcher.queue) + 1 > fetcher.queueSize {
		return false
	}

	fetcher.queue = append(fetcher.queue, op)
	return true
}


func (fetcher *FilesizeFetcher) removeFromQueue(op *FetchOperation) bool {
	fetcher.queueMutex.Lock()
	defer fetcher.queueMutex.Unlock()

	for i, queueOP := range fetcher.queue {
		if queueOP == op {
			fetcher.queue = append(fetcher.queue[:i], fetcher.queue[i+1:]...)
			return true
		}
	}

	return false
}

func updateFileList(dbEntry model.Torrent, info *metainfo.Info) error {
	log.Infof("TID %d has %d files.", dbEntry.ID, len(info.Files))
	for _, file := range info.Files {
		path := file.DisplayPath(info)
		fileExists := false
		for _, existingFile := range dbEntry.FileList {
			if existingFile.Path == path {
				fileExists = true
				break
			}
		}

		if !fileExists {
			log.Infof("Adding file %s to filelist of TID %d", path, dbEntry.ID)
			dbFile := model.File{
				TorrentID: dbEntry.ID,
				Path: path,
				Filesize: file.Length,
			}

			err := db.ORM.Create(&dbFile).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (fetcher *FilesizeFetcher) gotResult(r Result) {
	updatedSuccessfully := false
	if r.err != nil {
		log.Infof("Failed to get torrent filesize (TID: %d), err %v", r.operation.torrent.ID, r.err)
	} else if r.info.TotalLength() == 0 {
		log.Infof("Got length 0 for torrent TID: %d. Possible bug?", r.operation.torrent.ID)
	} else {
		log.Infof("Got length %d for torrent TID: %d. Updating.", r.info.TotalLength(), r.operation.torrent.ID)
		r.operation.torrent.Filesize = r.info.TotalLength()
		_, err := torrentService.UpdateTorrent(r.operation.torrent)
		if err != nil {
			log.Infof("Failed to update torrent TID: %d with new filesize", r.operation.torrent.ID)
		} else {
			updatedSuccessfully = true
		}

		// Also update the File list with FileInfo, I guess.
		err = updateFileList(r.operation.torrent, r.info)
		if err != nil {
			log.Infof("Failed to update file list of TID %d", r.operation.torrent.ID)
		}
	}

	if !updatedSuccessfully {
		fetcher.failedOperations[r.operation.torrent.ID] = struct{}{}
	}

	fetcher.removeFromQueue(r.operation)
}

func (fetcher *FilesizeFetcher) fillQueue() {
	toFill := fetcher.queueSize - len(fetcher.queue)

	if toFill <= 0 {
		return
	}

	oldest := time.Now().Add(0 - (time.Hour * time.Duration(24 * fetcher.maxDays)))
	params := serviceBase.CreateWhereParams("(filesize IS NULL OR filesize = 0) AND date > ?", oldest)
	// Get up to queueSize + len(failed) torrents, so we get at least some fresh new ones.
	dbTorrents, count, err := torrentService.GetTorrents(params, fetcher.queueSize + len(fetcher.failedOperations), 0)

	if err != nil {
		log.Infof("Failed to get torrents for filesize updating")
		return
	}
	
	if count == 0 {
		log.Infof("No torrents for filesize update")
		return
	}

	for _, T := range dbTorrents {
		if fetcher.isFetchingOrFailed(T) {
			continue
		}

		log.Infof("Added TID %d for filesize fetching", T.ID)
		operation := NewFetchOperation(fetcher, T)

		if fetcher.addToQueue(operation) {
			fetcher.wg.Add(1)
			go operation.Start(fetcher.results)
		} else {
			break
		}
	}
}

func (fetcher *FilesizeFetcher) run() {
	var result Result

	defer fetcher.wg.Done()

	done := 0
	fetcher.fillQueue()
	for done == 0 {
		select {
		case done = <-fetcher.done:
			break
		case result = <-fetcher.results:
			fetcher.gotResult(result)
			fetcher.fillQueue()
			break
		case <-fetcher.wakeUp.C:
			fetcher.fillQueue()
			break
		}
	}
}

func (fetcher *FilesizeFetcher) RunAsync() {
	fetcher.wg.Add(1)

	go fetcher.run()
}

func (fetcher *FilesizeFetcher) Close() error {
	fetcher.queueMutex.Lock()
	defer fetcher.queueMutex.Unlock()

	// Send the done event to every Operation
	for _, op := range fetcher.queue {
		op.done <- 1
	}

	fetcher.done <- 1
	fetcher.wg.Wait()
	return nil
}

func (fetcher *FilesizeFetcher) Wait() {
	fetcher.wg.Wait()
}

