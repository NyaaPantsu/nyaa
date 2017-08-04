package signals

import (
	"sync"

	"github.com/NyaaPantsu/nyaa/utils/log"
)

// registered interrupt callbacks.
// currently only used to gracefully close the server.
var intEvents struct {
	lock  sync.Mutex
	funcs []func()
}

// OnInterrupt handles signal interupts
func OnInterrupt(fn func()) {
	intEvents.lock.Lock()
	intEvents.funcs = append(intEvents.funcs, fn)
	intEvents.lock.Unlock()
}

func handleReload() {
	log.Info("Got SIGHUP")
	//router.ReloadTemplates()
	log.Info("reloaded templates")
}

// handle interrupt signal, platform independent
func interrupted() {
	intEvents.lock.Lock()
	for _, fn := range intEvents.funcs {
		fn()
	}
	intEvents.lock.Unlock()
}
