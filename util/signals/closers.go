package signals

import (
	"io"
	"sync"
)

var (
	closeAccess sync.Mutex
	closers     []io.Closer
)

// RegisterCloser adds an io.Closer to close on interrupt
func RegisterCloser(c io.Closer) {
	closeAccess.Lock()
	closers = append(closers, c)
	closeAccess.Unlock()
}

func closeClosers() {
	closeAccess.Lock()
	for idx := range closers {
		closers[idx].Close()
	}
	closeAccess.Unlock()
}
