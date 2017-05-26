package util

import (
	"net/http"
	"runtime/debug"

	"github.com/NyaaPantsu/nyaa/util/log"
)

// SendError : Send an error as response
func SendError(w http.ResponseWriter, err error, code int) {
	log.Warnf("%s:\n%s\n", err, debug.Stack())
	http.Error(w, err.Error(), code)
}
