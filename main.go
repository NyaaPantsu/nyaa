package main

import (
"bufio"
	"fmt"
	"flag"

	"github.com/ewhal/nyaa/util/log"
	"github.com/ewhal/nyaa/router"
	"github.com/ewhal/nyaa/config"


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
	conf_bind := conf.BindFlags()
	defaults := flag.Bool("print-defaults", false, "print the default configuration file on stdout")
	flag.Parse()
	if *defaults {
		stdout := bufio.NewWriter(os.Stdout)
		conf.Pretty(stdout)
		stdout.Flush()
		os.Exit(0)
	} else {
		conf_bind()
		RunServer(conf)
	}
}
