package main

import (
	"bufio"
	"flag"

	"context"
	"net/http"
	"os"
	"time"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/controllers"
	"github.com/NyaaPantsu/nyaa/controllers/databasedumps"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/utils/cookies"
	"github.com/NyaaPantsu/nyaa/utils/log"
	"github.com/NyaaPantsu/nyaa/utils/publicSettings"
	"github.com/NyaaPantsu/nyaa/utils/search"
	"github.com/NyaaPantsu/nyaa/utils/signals"
)

var buildversion string

// RunServer runs webapp mainloop
func RunServer(conf *config.Config) {
	// TODO Use config from cli
	os.Mkdir(databasedumpsController.DatabaseDumpPath, 0700)
	// TODO Use config from cli
	os.Mkdir(databasedumpsController.GPGPublicKeyPath, 0700)

	http.Handle("/", controllers.CSRFRouter)

	// Set up server,
	srv := &http.Server{
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	l, err := CreateHTTPListener(conf)
	log.CheckError(err)
	if err != nil {
		return
	}
	log.Infof("listening on %s", l.Addr())
	// http.Server.Shutdown closes associated listeners/clients.
	// context.Background causes srv to indefinitely try to
	// gracefully shutdown. add a timeout if this becomes a problem.
	signals.OnInterrupt(func() {
		srv.Shutdown(context.Background())
	})
	err = srv.Serve(l)
	if err == nil {
		log.Panic("http.Server.Serve never returns nil")
	}
	if err == http.ErrServerClosed {
		return
	}
	log.CheckError(err)
}

func main() {
	conf := config.Get()
	if buildversion != "" {
		conf.Build = buildversion
	} else {
		conf.Build = "unknown"
	}
	defaults := flag.Bool("print-defaults", false, "print the default configuration file on stdout")
	callback := config.BindFlags()
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
		callback()
		var err error
		models.ORM, err = models.GormInit(conf, models.DefaultLogger)
		if err != nil {
			log.Fatal(err.Error())
		}
		if config.Get().Search.EnableElasticSearch {
			models.ElasticSearchClient, _ = models.ElasticSearchInit()
		}
		err = publicSettings.InitI18n(conf.I18n, cookies.NewCurrentUserRetriever())
		if err != nil {
			log.Fatal(err.Error())
		}
		err = search.Configure(&conf.Search)
		if err != nil {
			log.Fatal(err.Error())
		}
		signals.Handle()
		if len(config.Get().Torrents.FileStorage) > 0 {
			err := os.MkdirAll(config.Get().Torrents.FileStorage, 0700)
			if err != nil {
				log.Fatal(err.Error())
			}
		}
		RunServer(conf)
	}
}
