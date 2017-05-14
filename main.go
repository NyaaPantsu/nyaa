package main

import (
	"bufio"
	"flag"

	"net/http"
	"os"
	"time"

	"github.com/ewhal/nyaa/cache"
	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/network"
	"github.com/ewhal/nyaa/router"
	"github.com/ewhal/nyaa/service/scraper"
	"github.com/ewhal/nyaa/service/torrent/filesizeFetcher"
	"github.com/ewhal/nyaa/util/languages"
	"github.com/ewhal/nyaa/util/log"
	"github.com/ewhal/nyaa/util/search"
	"github.com/ewhal/nyaa/util/signals"
)

// RunServer runs webapp mainloop
func RunServer(conf *config.Config) {
	http.Handle("/", router.Router)

	// Set up server,
	srv := &http.Server{
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
	l, err := network.CreateHTTPListener(conf)
	log.CheckError(err)
	if err == nil {
		// add http server to be closed gracefully
		signals.RegisterCloser(&network.GracefulHttpCloser{
			Server:   srv,
			Listener: l,
		})
		log.Infof("listening on %s", l.Addr())
		err := srv.Serve(l)
		if err != nil && err != network.ErrListenerStopped {
			log.CheckError(err)
		}

	}
}

// RunScraper runs tracker scraper mainloop
func RunScraper(conf *config.Config) {

	// bind to network
	pc, err := network.CreateScraperSocket(conf)
	if err != nil {
		log.Fatalf("failed to bind udp socket for scraper: %s", err)
	}
	// configure tracker scraperv
	var scraper *scraperService.Scraper
	scraper, err = scraperService.New(&conf.Scrape)
	if err != nil {
		pc.Close()
		log.Fatalf("failed to configure scraper: %s", err)
	}

	workers := conf.Scrape.NumWorkers
	if workers < 1 {
		workers = 1
	}

	// register udp socket with signals
	signals.RegisterCloser(pc)
	// register scraper with signals
	signals.RegisterCloser(scraper)
	// run udp scraper worker
	for workers > 0 {
		log.Infof("starting up worker %d", workers)
		go scraper.RunWorker(pc)
		workers--
	}
	// run scraper
	go scraper.Run()
	scraper.Wait()
}

// RunFilesizeFetcher runs the database filesize fetcher main loop
func RunFilesizeFetcher(conf *config.Config) {
	fetcher, err := filesizeFetcher.New(&conf.FilesizeFetcher)
	if err != nil {
		log.Fatalf("failed to start fetcher, %s", err)
		return
	}

	signals.RegisterCloser(fetcher)
	fetcher.RunAsync()
	fetcher.Wait()
}

func main() {
	conf := config.New()
	processFlags := conf.BindFlags()
	defaults := flag.Bool("print-defaults", false, "print the default configuration file on stdout")
	mode := flag.String("mode", "webapp", "which mode to run daemon in, either webapp, scraper or filesize_fetcher")
	flag.Float64Var(&conf.Cache.Size, "c", config.DefaultCacheSize, "size of the search cache in MB")

	flag.Parse()
	if *defaults {
		stdout := bufio.NewWriter(os.Stdout)
		err := conf.Pretty(stdout)
		if err != nil {
			log.Fatal(err.Error())
		}
		err = stdout.Flush()
		if err != nil {
			log.Fatal(err.Error())
		}
		os.Exit(0)
	} else {
		err := processFlags()
		if err != nil {
			log.CheckError(err)
		}
		db.ORM, err = db.GormInit(conf, db.DefaultLogger)
		if err != nil {
			log.Fatal(err.Error())
		}
		err = languages.InitI18n(conf.I18n)
		if err != nil {
			log.Fatal(err.Error())
		}
		err = cache.Configure(&conf.Cache)
		if err != nil {
			log.Fatal(err.Error())
		}
		err = search.Configure(&conf.Search)
		if err != nil {
			log.Fatal(err.Error())
		}
		go signals.Handle()
		if len(config.TorrentFileStorage) > 0 {
			err := os.MkdirAll(config.TorrentFileStorage, 0700)
			if err != nil {
				log.Fatal(err.Error())
			}
		}
		if *mode == "scraper" {
			RunScraper(conf)
		} else if *mode == "webapp" {
			RunServer(conf)
		} else if *mode == "filesize_fetcher" {
			RunFilesizeFetcher(conf)
		} else {
			log.Fatalf("invalid runtime mode: %s", *mode)
		}
	}
}
