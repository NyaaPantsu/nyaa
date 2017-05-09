package main

import (
	"bufio"
	"flag"

	"github.com/nicksnyder/go-i18n/i18n"

	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/network"
	"github.com/ewhal/nyaa/router"
	"github.com/ewhal/nyaa/util/log"
	"github.com/ewhal/nyaa/util/signals"

	"net/http"
	"os"
	"time"
)

func initI18N() {
	/* Initialize the languages translation */
	i18n.MustLoadTranslationFile("translations/en-us.all.json")
}

func RunServer(conf *config.Config) {
	http.Handle("/", router.Router)

	// Set up server,
	srv := &http.Server{
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	l, err := network.CreateHTTPListener(conf)
	log.CheckError(err)
	if err == nil {
		log.Infof("listening on %s", l.Addr())
		err := srv.Serve(l)
		log.CheckError(err)
	}
}

func main() {
	conf := config.New()
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
		db.ORM, err = db.GormInit(conf)
		if err != nil {
			log.Fatal(err.Error())
		}
		initI18N()
		go signals.Handle()
		if len(config.TorrentFileStorage) > 0 {
			err := os.MkdirAll(config.TorrentFileStorage, 0700)
			if err != nil {
				log.Fatal(err.Error())
			}
		}
		RunServer(conf)
	}
}
