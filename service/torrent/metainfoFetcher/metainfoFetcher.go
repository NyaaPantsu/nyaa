package metainfoFetcher;

import (
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/util/log"
	serviceBase "github.com/ewhal/nyaa/service"
	torrentService "github.com/ewhal/nyaa/service/torrent"
	"golang.org/x/time/rate"
	"math"
	"sync"
	"time"
)

type MetainfoFetcher struct {
	torrentClient    *torrent.Client
	results          chan Result
	queueSize        int
	timeout          int
	maxDays          int
	baseFailCooldown int
	maxFailCooldown  int
	done             chan int
	queue            []*FetchOperation
	queueMutex       sync.Mutex
	failedOperations map[uint]time.Time
	numFails         map[uint]int
	failsMutex       sync.Mutex
	wakeUp           *time.Ticker
	wg               sync.WaitGroup
}

func New(fetcherConfig *config.MetainfoFetcherConfig) (fetcher *MetainfoFetcher, err error) {
	clientConfig := torrent.Config{}
	// Well, it seems this is the right way to convert speed -> rate.Limiter
	// https://github.com/anacrolix/torrent/blob/master/cmd/torrent/main.go
	if fetcherConfig.UploadRateLimiter != -1 {
		clientConfig.UploadRateLimiter = rate.NewLimiter(rate.Limit(fetcherConfig.UploadRateLimiter * 1024), 256<<10)
	}
	if fetcherConfig.DownloadRateLimiter != -1 {
		clientConfig.DownloadRateLimiter = rate.NewLimiter(rate.Limit(fetcherConfig.DownloadRateLimiter * 1024), 1<<20)
	}

	client, err := torrent.NewClient(&clientConfig)

	fetcher = &MetainfoFetcher{
		torrentClient:    client,
		results:          make(chan Result, fetcherConfig.QueueSize),
		queueSize:        fetcherConfig.QueueSize,
		timeout:          fetcherConfig.Timeout,
		maxDays:          fetcherConfig.MaxDays,
		baseFailCooldown: fetcherConfig.BaseFailCooldown,
		maxFailCooldown:  fetcherConfig.MaxFailCooldown,
		done:             make(chan int, 1),
		failedOperations: make(map[uint]time.Time),
		numFails:         make(map[uint]int),
		wakeUp:           time.NewTicker(time.Second * time.Duration(fetcherConfig.WakeUpInterval)),
	}

	return
}

func (fetcher *MetainfoFetcher) isFetchingOrFailed(t model.Torrent) bool {
	for _, op := range fetcher.queue {
		if op.torrent.ID == t.ID {
			return true
		}
	}

	_, ok := fetcher.failedOperations[t.ID]
	return ok
}

func (fetcher *MetainfoFetcher) addToQueue(op *FetchOperation) bool {
	fetcher.queueMutex.Lock()
	defer fetcher.queueMutex.Unlock()

	if len(fetcher.queue) + 1 > fetcher.queueSize {
		return false
	}

	fetcher.queue = append(fetcher.queue, op)
	return true
}


func (fetcher *MetainfoFetcher) removeFromQueue(op *FetchOperation) bool {
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

func (fetcher *MetainfoFetcher) markAsFailed(tID uint) {
	fetcher.failsMutex.Lock()
	defer fetcher.failsMutex.Unlock()

	if n, ok := fetcher.numFails[tID]; ok {
		fetcher.numFails[tID] = n + 1
	} else {
		fetcher.numFails[tID] = 1
	}

	fetcher.failedOperations[tID] = time.Now()
}

func (fetcher *MetainfoFetcher) removeFromFailed(tID uint) {
	fetcher.failsMutex.Lock()
	defer fetcher.failsMutex.Unlock()

	delete(fetcher.failedOperations, tID)
}

func updateFileList(dbEntry model.Torrent, info *metainfo.Info) error {
	torrentFiles := info.UpvertedFiles()
	log.Infof("TID %d has %d files.", dbEntry.ID, len(torrentFiles))
	for _, file := range torrentFiles {
		path := file.DisplayPath(info)

		// Can't read FileList from the GetTorrents output, rely on the unique_index
		// to ensure no files are duplicated.
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

	return nil
}

func (fetcher *MetainfoFetcher) gotResult(r Result) {
	updatedSuccessfully := false
	if r.err != nil {
		log.Infof("Failed to get torrent metainfo (TID: %d), err %v", r.operation.torrent.ID, r.err)
	} else if r.info.TotalLength() == 0 {
		log.Infof("Got length 0 for torrent TID: %d. Possible bug?", r.operation.torrent.ID)
	} else {
		lengthOK := true

		if r.operation.torrent.Filesize != r.info.TotalLength() {
			log.Infof("Got length %d for torrent TID: %d. Updating.", r.info.TotalLength(), r.operation.torrent.ID)
			r.operation.torrent.Filesize = r.info.TotalLength()
			_, err := torrentService.UpdateTorrent(r.operation.torrent)
			if err != nil {
				log.Infof("Failed to update torrent TID: %d with new filesize", r.operation.torrent.ID)
				lengthOK = false
			}
		}

		filelistOK := true
		// Create the file list, if it's missing.
		if len(r.operation.torrent.FileList) == 0 {
			err := updateFileList(r.operation.torrent, r.info)
			if err != nil {
				log.Infof("Failed to update file list of TID %d", r.operation.torrent.ID)
				filelistOK = false
			}
		}

		updatedSuccessfully = lengthOK && filelistOK
	}

	if !updatedSuccessfully {
		fetcher.markAsFailed(r.operation.torrent.ID)
	}

	fetcher.removeFromQueue(r.operation)
}

func (fetcher *MetainfoFetcher) removeOldFailures() {
	// Cooldown is disabled
	if fetcher.baseFailCooldown < 0 {
		return
	}

	maxCd := time.Duration(fetcher.maxFailCooldown) * time.Second
	now := time.Now()
	for id, failTime := range fetcher.failedOperations {
		cdMult := int(math.Pow(2, float64(fetcher.numFails[id] - 1)))
		cd := time.Duration(cdMult * fetcher.baseFailCooldown) * time.Second
		if cd > maxCd {
			cd = maxCd
		}

		if failTime.Add(cd).Before(now) {
			log.Infof("Torrent TID %d gone through cooldown, removing from failures", id)
			// Deleting keys inside a loop seems to be safe.
			fetcher.removeFromFailed(id)
		}
	}
}

func (fetcher *MetainfoFetcher) fillQueue() {
	toFill := fetcher.queueSize - len(fetcher.queue)

	if toFill <= 0 {
		return
	}

	oldest := time.Now().Add(0 - (time.Hour * time.Duration(24 * fetcher.maxDays)))
	excludedIDS := make([]uint, 0, len(fetcher.failedOperations))
	for id, _ := range fetcher.failedOperations {
		excludedIDS = append(excludedIDS, id)
	}

	// Nice query lol
	// Select the torrents with no filesize, or without any rows with torrent_id in the files table, that are younger than fetcher.MaxDays
	var params serviceBase.WhereParams
	
	if len(excludedIDS) > 0 {
		params = serviceBase.CreateWhereParams("((filesize IS NULL OR filesize = 0) OR (torrents.torrent_id NOT IN (SELECT files.torrent_id FROM files WHERE files.torrent_id = torrents.torrent_id))) AND date > ? AND torrent_id NOT IN (?)", oldest, excludedIDS)
	} else {
		params = serviceBase.CreateWhereParams("((filesize IS NULL OR filesize = 0) OR (torrents.torrent_id NOT IN (SELECT files.torrent_id FROM files WHERE files.torrent_id = torrents.torrent_id))) AND date > ?", oldest)
	}
	dbTorrents, err := torrentService.GetTorrentsOrderByNoCount(&params, "", fetcher.queueSize, 0)

	if len(dbTorrents) == 0 {
		log.Infof("No torrents for filesize update")
		return
	}

	if err != nil {
		log.Infof("Failed to get torrents for metainfo updating")
		return
	}


	for _, T := range dbTorrents {
		if fetcher.isFetchingOrFailed(T) {
			continue
		}

		operation := NewFetchOperation(fetcher, T)
		if fetcher.addToQueue(operation) {
			log.Infof("Added TID %d for filesize fetching", T.ID)
			fetcher.wg.Add(1)
			go operation.Start(fetcher.results)
		} else {
			break
		}
	}
}

func (fetcher *MetainfoFetcher) run() {
	var result Result

	defer fetcher.wg.Done()

	done := 0
	fetcher.fillQueue()
	for done == 0 {
		fetcher.removeOldFailures()
		fetcher.fillQueue()
		select {
		case done = <-fetcher.done:
			break
		case result = <-fetcher.results:
			fetcher.gotResult(result)
			break
		case <-fetcher.wakeUp.C:
			break
		}
	}
}

func (fetcher *MetainfoFetcher) RunAsync() {
	fetcher.wg.Add(1)

	go fetcher.run()
}

func (fetcher *MetainfoFetcher) Close() error {
	fetcher.queueMutex.Lock()
	defer fetcher.queueMutex.Unlock()

	// Send the done event to every Operation
	for _, op := range fetcher.queue {
		op.done <- 1
	}

	fetcher.done <- 1
	fetcher.torrentClient.Close()
	log.Infof("Send done signal to everyone, waiting...")
	fetcher.wg.Wait()
	return nil
}

func (fetcher *MetainfoFetcher) Wait() {
	fetcher.wg.Wait()
}

