package controllers

import (
	"net/http"

	"github.com/NyaaPantsu/nyaa/utils/captcha"
	"github.com/gin-gonic/gin"
	"github.com/justinas/nosurf"
)

// Router variable for exporting the route configuration
var Router *gin.Engine

// CSRFRouter : CSRF protection for Router variable for exporting the route configuration
var CSRFRouter *nosurf.CSRFHandler

func init() {
	Router = gin.New()
	//Router.Use(gzip.Gzip(gzip.DefaultCompression)) // FIXME I can't make it work :/
	Router.Use(gin.Logger())
	Router.Use(gin.Recovery())
	Router.Use(errorMiddleware())

	// Static file handlers
	// TODO Use config from cli
	// TODO Make sure the directory exists
	Router.StaticFS("/css/", http.Dir("./public/css/"))
	Router.StaticFS("/js/", http.Dir("./public/js/"))
	Router.StaticFS("/img/", http.Dir("./public/img/"))
	Router.StaticFS("/apidoc/", http.Dir("./apidoc/"))
	Router.StaticFS("/dbdumps/", http.Dir(DatabaseDumpPath))
	Router.StaticFS("/gpg/", http.Dir(GPGPublicKeyPath))

	// We don't need CSRF here
	Router.Any("/", SearchHandler)
	Router.Any("/p/:page", SearchHandler)
	Router.Any("/search", SearchHandler)
	Router.Any("/search/:page", SearchHandler)
	Router.Any("/verify/email/:token", UserVerifyEmailHandler)
	Router.Any("/faq", FaqHandler)
	Router.Any("/activities", ActivityListHandler)
	Router.Any("/activities/p/:page", ActivityListHandler)
	Router.Any("/feed", RSSHandler)
	Router.Any("/feed/p/:page", RSSHandler)
	Router.Any("/feed/magnet", RSSMagnetHandler)
	Router.Any("/feed/magnet/p/:page", RSSMagnetHandler)
	Router.Any("/feed/torznab", RSSTorznabHandler)
	Router.Any("/feed/torznab/api", RSSTorznabHandler)
	Router.Any("/feed/torznab/p/:page", RSSTorznabHandler)
	Router.Any("/feed/eztv", RSSEztvHandler)
	Router.Any("/feed/eztv/p/:page", RSSEztvHandler)

	// !!! This line need to have the same download location as the one define in config.TorrentStorageLink !!!
	Router.Any("/download/:hash", DownloadTorrent)

	Router.Any("/upload", UploadHandler)
	Router.POST("/login", UserLoginPostHandler)
	Router.GET("/register", UserRegisterFormHandler)
	Router.GET("/login", UserLoginFormHandler)
	Router.POST("/register", UserRegisterPostHandler)
	Router.POST("/logout", UserLogoutHandler)
	Router.GET("/notifications", UserNotificationsHandler)

	reportRoutes := Router.Group("/report")
	{
		//reporting a torrent
		reportRoutes.GET("/:id", ReportViewTorrentHandler)
		reportRoutes.POST("/:id", ReportTorrentHandler)
	}
	torrentViewRoutes := Router.Group("/view")
	{
		torrentViewRoutes.GET("/:id", ViewHandler)
		torrentViewRoutes.HEAD("/:id", ViewHeadHandler)
		torrentViewRoutes.POST("/:id", PostCommentHandler)
	}
	torrentRoutes := Router.Group("/torrent")
	{
		torrentRoutes.GET("/", TorrentEditUserPanel)
		torrentRoutes.POST("/", TorrentPostEditUserPanel)
		torrentRoutes.GET("/delete", TorrentDeleteUserPanel)
	}
	userRoutes := Router.Group("/user")
	{
		userRoutes.GET("/:id/:username", UserProfileHandler)
		userRoutes.GET("/:id/:username/follow", UserFollowHandler)
		userRoutes.GET("/:id/:username/edit", UserDetailsHandler)
		userRoutes.POST("/:id/:username/edit", UserProfileFormHandler)
		userRoutes.GET("/:id/:username/apireset", UserAPIKeyResetHandler)
		userRoutes.GET("/:id/:username/feed/*page", RSSHandler)
	}
	// We don't need CSRF here
	api := Router.Group("/api")
	{
		api.GET("", APIHandler)
		api.GET("/p/:page", APIHandler)
		api.GET("/view/:id", APIViewHandler)
		api.HEAD("/view/:id", APIViewHeadHandler)
		api.POST("/upload", APIUploadHandler)
		api.POST("/login", APILoginHandler)
		api.GET("/token/check", APICheckTokenHandler)
		api.GET("/token/refresh", APIRefreshTokenHandler)
		api.Any("/search", APISearchHandler)
		api.Any("/search/p/:page", APISearchHandler)
		api.PUT("/update", APIUpdateHandler)
	}

	// INFO Everything under /mod should be wrapped by wrapModHandler. This make
	// sure the page is only accessible by moderators
	// We don't need CSRF here
	modRoutes := Router.Group("/mod", modMiddleware())
	{
		modRoutes.Any("/", IndexModPanel)
		modRoutes.GET("/torrents", TorrentsListPanel)
		modRoutes.GET("/torrents/p/:page", TorrentsListPanel)
		modRoutes.POST("/torrents", TorrentsPostListPanel)
		modRoutes.POST("/torrents/p/:page", TorrentsPostListPanel)
		modRoutes.GET("/torrents/deleted", DeletedTorrentsModPanel)
		modRoutes.GET("/torrents/deleted/p/:page", DeletedTorrentsModPanel)
		modRoutes.POST("/torrents/deleted", DeletedTorrentsPostPanel)
		modRoutes.POST("/torrents/deleted/p/:page", DeletedTorrentsPostPanel)
		modRoutes.Any("/reports", TorrentReportListPanel)
		modRoutes.Any("/reports/p/:page", TorrentReportListPanel)
		modRoutes.Any("/users", UsersListPanel)
		modRoutes.Any("/users/p/:page", UsersListPanel)
		modRoutes.Any("/comments", CommentsListPanel)
		modRoutes.Any("/comments/p/:page", CommentsListPanel)
		modRoutes.Any("/comment", CommentsListPanel) // TODO
		modRoutes.GET("/torrent", TorrentEditModPanel)
		modRoutes.POST("/torrent", TorrentPostEditModPanel)
		modRoutes.Any("/torrent/delete", TorrentDeleteModPanel)
		modRoutes.Any("/torrent/block", TorrentBlockModPanel)
		modRoutes.Any("/report/delete", TorrentReportDeleteModPanel)
		modRoutes.Any("/comment/delete", CommentDeleteModPanel)
		modRoutes.GET("/reassign", TorrentReassignModPanel)
		modRoutes.POST("/reassign", TorrentPostReassignModPanel)
		apiMod := modRoutes.Group("/api")
		apiMod.Any("/torrents", APIMassMod)
	}

	Router.Any("/captcha/*hash", captcha.ServeFiles)

	Router.Any("/dumps", DatabaseDumpHandler)

	Router.GET("/settings", SeePublicSettingsHandler)
	Router.POST("/settings", ChangePublicSettingsHandler)

	// Adding pprof support
	pprofRoutes := Router.Group("/debug/pprof", modMiddleware())
	{
		pprofRoutes.GET("/", PprofIndex)
		pprofRoutes.GET("/block", PprofIndex)
		pprofRoutes.GET("/heap", PprofIndex)
		pprofRoutes.GET("/profile", PprofProfile)
		pprofRoutes.POST("/symbol", PprofSymbol)
		pprofRoutes.GET("/symbol", PprofSymbol)
		pprofRoutes.GET("/trace", PprofTrace)
	}
	CSRFRouter = nosurf.New(Router)
	CSRFRouter.ExemptRegexp("/api(?:/.+)*")
	CSRFRouter.ExemptRegexp("/mod(?:/.+)*")
	CSRFRouter.ExemptPath("/upload")
	CSRFRouter.ExemptPath("/user/login")
	CSRFRouter.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid CSRF tokens", http.StatusBadRequest)
	}))
	CSRFRouter.SetBaseCookie(http.Cookie{
		Path:   "/",
		MaxAge: nosurf.MaxAge,
	})

}
