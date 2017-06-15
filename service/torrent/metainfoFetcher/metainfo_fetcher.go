package metainfoFetcher

import (
	"sync"
	"time"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/db"
	"github.com/NyaaPantsu/nyaa/model"
	serviceBase "github.com/NyaaPantsu/nyaa/service"
	torrentService "github.com/NyaaPantsu/nyaa/service/torrent"
	"github.com/NyaaPantsu/nyaa/util/log"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"golang.org/x/time/rate"
)

// MetainfoFetcher Struct
type MetainfoFetcher struct {
	torrentClient    *torrent.Client
	results          chan Result
	queueSize        int
	timeout          int
	maxDays          int
	baseFailCooldown int
	maxFailCooldown  int
	newTorrentsOnly  bool
	done             chan struct{}
	queue            []*FetchOperation
	queueMutex       sync.Mutex
	failedOperations map[uint]time.Time
	numFails         map[uint]int
	failsMutex       sync.Mutex
	wakeUp           *time.Ticker
	wg               sync.WaitGroup
}

// New : Creates a MetainfoFetcher struct
func New(fetcherConfig *config.MetainfoFetcherConfig) (*MetainfoFetcher, error) {
	clientConfig := torrent.Config{}

	// Well, it seems this is the right way to convert speed -> rate.Limiter
	// https://github.com/anacrolix/torrent/blob/master/cmd/torrent/main.go
	const uploadBurst = 0x40000    // 256K
	const downloadBurst = 0x100000 // 1M
	uploadLimit := fetcherConfig.UploadRateLimitKiB * 1024
	downloadLimit := fetcherConfig.DownloadRateLimitKiB * 1024
	if uploadLimit > 0 {
		limit := rate.Limit(uploadLimit)
		limiter := rate.NewLimiter(limit, uploadBurst)
		clientConfig.UploadRateLimiter = limiter
	}
	if downloadLimit > 0 {
		limit := rate.Limit(downloadLimit)
		limiter := rate.NewLimiter(limit, downloadBurst)
		clientConfig.DownloadRateLimiter = limiter
	}

	client, err := torrent.NewClient(&clientConfig)
	if err != nil {
		return nil, err
	}

	fetcher := &MetainfoFetcher{
		torrentClient:    client,
		results:          make(chan Result, fetcherConfig.QueueSize),
		queueSize:        fetcherConfig.QueueSize,
		timeout:          fetcherConfig.Timeout,
		maxDays:          fetcherConfig.MaxDays,
		newTorrentsOnly:  fetcherConfig.FetchNewTorrentsOnly,
		baseFailCooldown: fetcherConfig.BaseFailCooldown,
		maxFailCooldown:  fetcherConfig.MaxFailCooldown,
		done:             make(chan struct{}, 1),
		failedOperations: make(map[uint]time.Time),
		numFails:         make(map[uint]int),
		wakeUp:           time.NewTicker(time.Second * time.Duration(fetcherConfig.WakeUpInterval)),
	}

	return fetcher, nil
}

func (fetcher *MetainfoFetcher) addToQueue(op *FetchOperation) bool {
	fetcher.queueMutex.Lock()
	defer fetcher.queueMutex.Unlock()

	if fetcher.queueSize > 0 && len(fetcher.queue) > fetcher.queueSize-1 {
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
		var path []string
		if file.Path != nil {
			path = file.Path
		} else {
			// If it's nil, use the torrent name (info.Name) as the path (single-file torrent)
			path = append(path, info.Name)
		}

		// Can't read FileList from the GetTorrents output, rely on the unique_index
		// to ensure no files are duplicated.
		log.Infof("Adding file %s to filelist of TID %d", file.DisplayPath(info), dbEntry.ID)
		dbFile := model.File{
			TorrentID: dbEntry.ID,
			Filesize:  file.Length,
		}
		err := dbFile.SetPath(path)
		if err != nil {
			return err
		}

		err = db.ORM.Create(&dbFile).Error
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
			_, err := torrentService.UpdateTorrent(&r.operation.torrent)
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
	if fetcher.baseFailCooldown < 0 {
		// XXX: Cooldown is disabled.
		// this means that if any attempt to fetch metadata fails
		// it will never be retried. it also means that
		// fetcher.failedOperations will keep accumulating torrent IDs
		// that are never freed.
		return
	}

	max := time.Duration(fetcher.maxFailCooldown) * time.Second
	now := time.Now()
	for id, failTime := range fetcher.failedOperations {
		// double the amount of time waited for ever failed attempt.
		// | nfailed | cooldown
		// | 1       | 2 * base
		// | 2       | 4 * base
		// | 3       | 8 * base
		// integers inside fetcher.numFails are never less than or equal to zero
		mul := 1 << uint(fetcher.numFails[id]-1)
		cd := time.Duration(mul*fetcher.baseFailCooldown) * time.Second
		if cd > max {
			cd = max
		}

		if failTime.Add(cd).Before(now) {
			log.Infof("Torrent TID %d gone through cooldown, removing from failures", id)
			// Deleting keys inside a loop seems to be safe.
			fetcher.removeFromFailed(id)
		}
	}
}

func (fetcher *MetainfoFetcher) fillQueue() {
	if left := fetcher.queueSize - len(fetcher.queue); left <= 0 {
		// queue is already full.
		return
	}

	oldest := time.Now().Add(0 - (time.Hour * time.Duration(24*fetcher.maxDays)))
	excludedIDS := make([]uint, 0, len(fetcher.failedOperations))
	for id := range fetcher.failedOperations {
		excludedIDS = append(excludedIDS, id)
	}

	tFiles := config.Conf.Models.FilesTableName
	tTorrents := config.Conf.Models.TorrentsTableName
	// Select the torrents with no filesize, or without any rows with torrent_id in the files table...
	queryString := "((filesize IS NULL OR filesize = 0) OR (" + tTorrents + ".torrent_id NOT " +
		"IN (SELECT " + tFiles + ".torrent_id FROM " + tFiles + " WHERE " + tFiles +
		".torrent_id = " + tTorrents + ".torrent_id)))"
	var whereParamsArgs []interface{}

	// that are newer than maxDays...
	queryString += " AND date > ? "
	whereParamsArgs = append(whereParamsArgs, oldest)

	// that didn't fail recently...
	if len(excludedIDS) > 0 {
		queryString += " AND torrent_id NOT IN (?) "
		whereParamsArgs = append(whereParamsArgs, excludedIDS)
	}

	// and, if true, that aren't from the old Nyaa database
	if fetcher.newTorrentsOnly {
		queryString += " AND torrent_id > ? "
		whereParamsArgs = append(whereParamsArgs, config.Conf.Models.LastOldTorrentID)
	}

	params := serviceBase.CreateWhereParams(queryString, whereParamsArgs...)
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
		for _, v := range fetcher.queue {
			// skip torrents that are already being processed.
			if v.torrent.ID == T.ID {
				continue
			}
		}
		if _, ok := fetcher.failedOperations[T.ID]; ok {
			// do not start new jobs that have recently failed.
			// these are on cooldown and will be removed from
			// fetcher.failedOperations when time is up.
			continue
		}

		operation := NewFetchOperation(fetcher, T)
		if !fetcher.addToQueue(operation) {
			// queue is full, stop and wait for results
			break
		}

		log.Infof("Added TID %d for filesize fetching", T.ID)
		fetcher.wg.Add(1)
		go operation.Start(fetcher.results)
	}
}

func (fetcher *MetainfoFetcher) run() {
	var result Result

	defer fetcher.wg.Done()

	for {
		fetcher.removeOldFailures()
		fetcher.fillQueue()

		select {
		case <-fetcher.done:
			log.Infof("Got done signal on main loop, leaving...")
			return
		case result = <-fetcher.results:
			fetcher.gotResult(result)
		case <-fetcher.wakeUp.C:
			log.Infof("Got wake up signal...")
		}
	}
}

// RunAsync method
func (fetcher *MetainfoFetcher) RunAsync() {
	fetcher.wg.Add(1)

	go fetcher.run()
}

// Close method
func (fetcher *MetainfoFetcher) Close() error {
	fetcher.queueMutex.Lock()
	defer fetcher.queueMutex.Unlock()

	// Send the done event to every Operation
	for _, op := range fetcher.queue {
		op.done <- struct{}{}
	}

	fetcher.done <- struct{}{}
	fetcher.torrentClient.Close()
	log.Infof("Send done signal to everyone, waiting...")
	fetcher.wg.Wait()
	return nil
}

// Wait method
func (fetcher *MetainfoFetcher) Wait() {
	fetcher.wg.Wait()
}
