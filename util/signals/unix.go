// +build !win32

package signals

import (
	"github.com/ewhal/nyaa/router"
	"github.com/ewhal/nyaa/util/log"
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
	signal.Notify(chnl, syscall.SIGHUP)
	for {
		sig, ok := <-chnl
		if !ok {
			break
		}
		switch sig {
		case syscall.SIGHUP:
			handleReload()
			break
		default:
			break
		}
	}
}
