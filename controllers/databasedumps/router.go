package databasedumpsController

import "github.com/NyaaPantsu/nyaa/controllers/router"

func init() {
	router.Get().Any("/dumps", DatabaseDumpHandler)
}
