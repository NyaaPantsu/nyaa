package main

import (
	"bufio"
	"flag"

	"context"
	"net/http"
	"os"
	"time"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/db"
	"github.com/NyaaPantsu/nyaa/network"
	"github.com/NyaaPantsu/nyaa/router"
	"github.com/NyaaPantsu/nyaa/service/user"
	"github.com/NyaaPantsu/nyaa/util/log"
	"github.com/NyaaPantsu/nyaa/util/publicSettings"
	"github.com/NyaaPantsu/nyaa/util/search"
	"github.com/NyaaPantsu/nyaa/util/signals"
)

var buildversion string

// RunServer runs webapp mainloop
func RunServer(conf *config.Config) {
	// TODO Use config from cli
	os.Mkdir(router.DatabaseDumpPath, 0700)
	// TODO Use config from cli
	os.Mkdir(router.GPGPublicKeyPath, 0700)

	http.Handle("/", router.CSRFRouter)

	// Set up server,
	srv := &http.Server{
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	l, err := network.CreateHTTPListener(conf)
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
	conf := config.Conf
	if buildversion != "" {
		conf.Build = buildversion
	} else {
		conf.Build = "unknown"
	}
	processFlags := conf.BindFlags()
	defaults := flag.Bool("print-defaults", false, "print the default configuration file on stdout")

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
		db.ElasticSearchClient, _ = db.ElasticSearchInit()
		err = publicSettings.InitI18n(conf.I18n, userService.NewCurrentUserRetriever())
		if err != nil {
			log.Fatal(err.Error())
		}
		err = search.Configure(&conf.Search)
		if err != nil {
			log.Fatal(err.Error())
		}
		signals.Handle()
		if len(config.Conf.Torrents.FileStorage) > 0 {
			err := os.MkdirAll(config.Conf.Torrents.FileStorage, 0700)
			if err != nil {
				log.Fatal(err.Error())
			}
		}
		RunServer(conf)
	}
}
