package signals

import (
	"github.com/NyaaPantsu/nyaa/router"
	"github.com/NyaaPantsu/nyaa/util/log"
)

func handleReload() {
	log.Info("Got SIGHUP")
	router.ReloadTemplates()
	log.Info("reloaded templates")
}

// handle interrupt signal, platform independent
func interrupted() {
	closeClosers()
}
