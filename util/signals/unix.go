// +build !win32

package signals

import (
	"github.com/NyaaPantsu/nyaa/router"
	"github.com/NyaaPantsu/nyaa/util/log"
	"os"
	"os/signal"
	"syscall"
)

func handleReload() {
	log.Info("Got SIGHUP")
	router.ReloadTemplates()
	log.Info("reloaded templates")

}

// Handle signals
// returns when done
func Handle() {
	chnl := make(chan os.Signal)
	signal.Notify(chnl, syscall.SIGHUP, os.Interrupt)
	for {
		sig, ok := <-chnl
		if !ok {
			break
		}
		switch sig {
		case syscall.SIGHUP:
			handleReload()
			break
		case os.Interrupt:
			interrupted()
			return
		default:
			break
		}
	}
}

// unix implementation of interrupt
// called in interrupted()
func handleInterrupt() {
	// XXX: put unix specific cleanup here as needed
}
