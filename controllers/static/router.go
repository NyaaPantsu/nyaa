package staticController

import (
	"net/http"

	"github.com/NyaaPantsu/nyaa/controllers/databasedumps"
	"github.com/NyaaPantsu/nyaa/controllers/router"
)

func init() {
	// Static file handlers
	// TODO Use config from cli
	// TODO Make sure the directory exists
	router.Get().StaticFS("/css/", http.Dir("./public/css/"))
	router.Get().StaticFS("/js/", http.Dir("./public/js/"))
	router.Get().StaticFS("/img/", http.Dir("./public/img/"))
	router.Get().StaticFS("/apidoc/", http.Dir("./apidoc/"))
	router.Get().StaticFS("/dbdumps/", http.Dir(databasedumpsController.DatabaseDumpPath))
	router.Get().StaticFS("/gpg/", http.Dir(databasedumpsController.GPGPublicKeyPath))
}
