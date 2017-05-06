package main

import (
	"bufio"
	"flag"
	"fmt"

	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/router"
	"github.com/ewhal/nyaa/util/log"

	"net/http"
	"os"
	"time"
)

func RunServer(conf *config.Config) {
	http.Handle("/", router.Router)

	// Set up server,
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

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
		RunServer(conf)
	}
}
