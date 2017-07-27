package controllers

import (
	"net/http"

	_ "github.com/NyaaPantsu/nyaa/controllers/activities"    // activities controller
	_ "github.com/NyaaPantsu/nyaa/controllers/api"           // api controller
	_ "github.com/NyaaPantsu/nyaa/controllers/captcha"       // captcha controller
	_ "github.com/NyaaPantsu/nyaa/controllers/databasedumps" // databasedumps controller
	_ "github.com/NyaaPantsu/nyaa/controllers/faq"           // faq controller
	_ "github.com/NyaaPantsu/nyaa/controllers/feed"          // feed controller
	_ "github.com/NyaaPantsu/nyaa/controllers/middlewares"   // middlewares
	_ "github.com/NyaaPantsu/nyaa/controllers/moderator"     // moderator controller
	_ "github.com/NyaaPantsu/nyaa/controllers/oauth"         // oauth2 controller
	_ "github.com/NyaaPantsu/nyaa/controllers/pprof"         // pprof controller
	_ "github.com/NyaaPantsu/nyaa/controllers/report"        // report controller
	"github.com/NyaaPantsu/nyaa/controllers/router"
	_ "github.com/NyaaPantsu/nyaa/controllers/search"   // search controller
	_ "github.com/NyaaPantsu/nyaa/controllers/settings" // settings controller
	_ "github.com/NyaaPantsu/nyaa/controllers/static"   // static files
	_ "github.com/NyaaPantsu/nyaa/controllers/torrent"  // torrent controller
	_ "github.com/NyaaPantsu/nyaa/controllers/upload"   // upload controller
	_ "github.com/NyaaPantsu/nyaa/controllers/user"     // user controller
	"github.com/justinas/nosurf"
)

// CSRFRouter : CSRF protection for Router variable for exporting the route configuration
var CSRFRouter *nosurf.CSRFHandler

func init() {
	CSRFRouter = nosurf.New(router.Get())
	CSRFRouter.ExemptRegexp("/api(?:/.+)*")
	CSRFRouter.ExemptRegexp("/mod(?:/.+)*")
	CSRFRouter.ExemptPath("/upload")
	CSRFRouter.ExemptPath("/user/login")
	CSRFRouter.ExemptPath("/oauth2/token")
	CSRFRouter.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid CSRF tokens", http.StatusBadRequest)
	}))
	CSRFRouter.SetBaseCookie(http.Cookie{
		Path:   "/",
		MaxAge: nosurf.MaxAge,
	})

}
