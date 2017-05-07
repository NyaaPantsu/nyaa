package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/nicksnyder/go-i18n/i18n"

	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/router"
	"github.com/ewhal/nyaa/util/log"
	"github.com/ewhal/nyaa/util/signals"

	"net/http"
	"os"
	"time"
)

func initI18N() {
	/* Initialize the languages translation */
	i18n.MustLoadTranslationFile("service/user/locale/en-us.all.json")
}

func RunServer(conf *config.Config) {
	http.Handle("/", router.Router)

	// Set up server,
	addr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	srv := &http.Server{
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Infof("listening on %s", addr)

	err := srv.ListenAndServe()
	log.CheckError(err)
}

func main() {
	conf := config.NewConfig()
	process_flags := conf.BindFlags()
	defaults := flag.Bool("print-defaults", false, "print the default configuration file on stdout")
	flag.Parse()
	if *defaults {
		stdout := bufio.NewWriter(os.Stdout)
		conf.Pretty(stdout)
		stdout.Flush()
		os.Exit(0)
	} else {
		err := process_flags()
		if err != nil {
			log.CheckError(err)
		}
		db.ORM, _ = db.GormInit(conf)
		initI18N()
		go signals.Handle()
		if len(config.TorrentFileStorage) > 0 {
			os.MkdirAll(config.TorrentFileStorage, 0755)
		}
		RunServer(conf)
	}
}
