// +build win32

package signals

import (
	"os"
	"os/signal"
)

func Handle() {
	// TODO: Something about SIGHUP for Windows, stdin could be used
	chnl := make(chan os.Signal)
	signal.Notify(chnl, os.Interrupt)
	go func(chnl chan os.Signal) {
		for sig := range chnl {
			switch sig {
			case os.Interrupt:
				// this also closes listeners
				interrupted()
				// XXX: put any win32 specific cleanup here as needed
				return
			}
		}
	}(chnl)
}
