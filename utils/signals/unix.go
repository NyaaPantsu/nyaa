// +build !win32

package signals

import (
	"os"
	"os/signal"
	"syscall"
)

func Handle() {
	chnl := make(chan os.Signal)
	signal.Notify(chnl, syscall.SIGHUP, os.Interrupt)
	go func(chnl chan os.Signal) {
		for sig := range chnl {
			switch sig {
			case syscall.SIGHUP:
				handleReload()
			case os.Interrupt:
				interrupted()
				// XXX: put unix specific cleanup here as needed
				return
			}
		}
	}(chnl)
}
