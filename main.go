package main

import (
	"bufio"
	"fmt"

	"github.com/ewhal/nyaa/config"
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
	conf := config.GetInstance()
	if *config.PrintDefaults {
		stdout := bufio.NewWriter(os.Stdout)
		conf.Pretty(stdout)
		stdout.Flush()
		os.Exit(0)
	} else {
		RunServer(conf)
	}
}
