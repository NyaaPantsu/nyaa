package main

import (
	"bufio"
	"flag"

	"net/http"
	"os"
	"time"

	"github.com/NyaaPantsu/nyaa/cache"
	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/db"
	"github.com/NyaaPantsu/nyaa/network"
	"github.com/NyaaPantsu/nyaa/router"
	"github.com/NyaaPantsu/nyaa/service/scraper"
	"github.com/NyaaPantsu/nyaa/service/torrent/metainfoFetcher"
	"github.com/NyaaPantsu/nyaa/service/user"
	"github.com/NyaaPantsu/nyaa/util/languages"
	"github.com/NyaaPantsu/nyaa/util/log"
	"github.com/NyaaPantsu/nyaa/util/search"
	"github.com/NyaaPantsu/nyaa/util/signals"
)

// RunServer runs webapp mainloop
func RunServer(conf *config.Config) {
	// TODO Use config from cli
	os.Mkdir(router.DatabaseDumpPath, 700)
	// TODO Use config from cli
	os.Mkdir(router.GPGPublicKeyPath, 700)
	http.Handle("/", router.Router)

	// Set up server,
	srv := &http.Server{
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
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

// RunMetainfoFetcher runs the database filesize fetcher main loop
func RunMetainfoFetcher(conf *config.Config) {
	fetcher, err := metainfoFetcher.New(&conf.MetainfoFetcher)
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
	mode := flag.String("mode", "webapp", "which mode to run daemon in, either webapp, scraper or metainfo_fetcher")
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
		err = languages.InitI18n(conf.I18n, userService.NewCurrentUserRetriever())
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
		signals.Handle()
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
		} else if *mode == "metainfo_fetcher" {
			RunMetainfoFetcher(conf)
		} else {
			log.Fatalf("invalid runtime mode: %s", *mode)
		}
	}
}
