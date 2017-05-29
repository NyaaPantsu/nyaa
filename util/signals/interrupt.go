package signals

import (
	"github.com/NyaaPantsu/nyaa/router"
	"github.com/NyaaPantsu/nyaa/util/log"
	"sync"
)

// registered interrupt callbacks.
// currently only used to gracefully close the server.
var intEvents struct {
	lock sync.Mutex
	funcs []func()
}

func OnInterrupt(fn func()) {
	intEvents.lock.Lock()
	intEvents.funcs = append(intEvents.funcs, fn)
	intEvents.lock.Unlock()
}

func handleReload() {
	log.Info("Got SIGHUP")
	router.ReloadTemplates()
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
