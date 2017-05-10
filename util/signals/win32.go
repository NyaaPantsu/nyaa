// +build win32

package signals

import (
	"os"
	"os/signal"
)

func Handle() {
	// TODO: Something about SIGHUP for Windows

	chnl := make(chan os.Signal)
	signal.Notify(chnl, os.Interrupt)
	for {
		sig, ok := <-chnl
		if !ok {
			break
		}
		switch sig {
		case os.Interrupt:
			// this also closes listeners
			interrupted()
			return
		default:
			break
		}
	}
}

// win32 interrupt handler
// called in interrupted()
func handleInterrupt() {
	// XXX: put any win32 specific cleanup here as needed
}
