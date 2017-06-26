package router

import (
	"net/http"

	"github.com/NyaaPantsu/nyaa/service/captcha"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/justinas/nosurf"
)

// Router variable for exporting the route configuration
var Router *gin.Engine

// CSRFRouter : CSRF protection for Router variable for exporting the route configuration
var CSRFRouter *nosurf.CSRFHandler

func init() {
	Router = gin.Default()
	Router.Use(gzip.Gzip(gzip.DefaultCompression))

	// Static file handlers
	// TODO Use config from cli
	// TODO Make sure the directory exists
	Router.StaticFS("/css/", http.Dir("./public/css/"))
	Router.StaticFS("/js/", http.Dir("./public/js/"))
	Router.StaticFS("/img/", http.Dir("./public/img/"))
	Router.StaticFS("/dbdumps/", http.Dir(DatabaseDumpPath))
	Router.StaticFS("/gpg/", http.Dir(GPGPublicKeyPath))

	// We don't need CSRF here
	Router.Any("/", SearchHandler).Use(errorMiddleware())
	Router.Any("/page/:page", SearchHandler).Use(errorMiddleware())
	Router.Any("/search", SearchHandler).Use(errorMiddleware())
	Router.Any("/search/:page", SearchHandler).Use(errorMiddleware())
	Router.Any("/verify/email/:token", UserVerifyEmailHandler).Use(errorMiddleware())
	Router.Any("/faq", FaqHandler).Use(errorMiddleware())
	Router.Any("/activities", ActivityListHandler).Use(errorMiddleware())
	Router.Any("/feed", RSSHandler)
	Router.Any("/feed/:page", RSSHandler)
	Router.Any("/feed/magnet", RSSMagnetHandler)
	Router.Any("/feed/magnet/:page", RSSMagnetHandler)
	Router.Any("/feed/torznab", RSSTorznabHandler)
	Router.Any("/feed/torznab/api", RSSTorznabHandler)
	Router.Any("/feed/torznab/:page", RSSTorznabHandler)
	Router.Any("/feed/eztv", RSSEztvHandler)
	Router.Any("/feed/eztv/:page", RSSEztvHandler)

	// !!! This line need to have the same download location as the one define in config.TorrentStorageLink !!!
	Router.Any("/download/:hash", DownloadTorrent)

	// For now, no CSRF protection here, as API is not usable for uploads
	Router.Any("/upload", UploadHandler)
	Router.Any("/user/login", UserLoginPostHandler)

	torrentViewRoutes := Router.Group("/view", errorMiddleware())
	{
		torrentViewRoutes.GET("/:id", ViewHandler)
		torrentViewRoutes.HEAD("/:id", ViewHeadHandler)
		torrentViewRoutes.POST("/:id", PostCommentHandler)
	}
	torrentRoutes := Router.Group("/torrent", errorMiddleware())
	{
		torrentRoutes.GET("/", TorrentEditUserPanel)
		torrentRoutes.POST("/", TorrentPostEditUserPanel)
		torrentRoutes.GET("/delete", TorrentDeleteUserPanel)
	}
	userRoutes := Router.Group("/user").Use(errorMiddleware())
	{
		userRoutes.GET("/register", UserRegisterFormHandler)
		userRoutes.GET("/login", UserLoginFormHandler)
		userRoutes.POST("/register", UserRegisterPostHandler)
		userRoutes.POST("/logout", UserLogoutHandler)
		userRoutes.GET("/:id/:username", UserProfileHandler)
		userRoutes.GET("/:id/:username/follow", UserFollowHandler)
		userRoutes.GET("/:id/:username/edit", UserDetailsHandler)
		userRoutes.POST("/:id/:username/edit", UserProfileFormHandler)
		userRoutes.GET("/:id/:username/apireset", UserAPIKeyResetHandler)
		userRoutes.GET("/notifications", UserNotificationsHandler)
		userRoutes.GET("/:id/:username/feed/*page", RSSHandler)
	}
	// We don't need CSRF here
	api := Router.Group("/api")
	{
		api.GET("", APIHandler)
		api.GET("/", APIHandler)
		api.GET("/{page:[0-9]*}", APIHandler)
		api.GET("/view/:id", APIViewHandler)
		api.HEAD("/view/:id", APIViewHeadHandler)
		api.POST("/upload", APIUploadHandler)
		api.POST("/login", APILoginHandler)
		api.GET("/token/check", APICheckTokenHandler)
		api.GET("/token/refresh", APIRefreshTokenHandler)
		api.Any("/search", APISearchHandler)
		api.Any("/search/{page}", APISearchHandler)
		api.PUT("/update", APIUpdateHandler)
	}

	// INFO Everything under /mod should be wrapped by wrapModHandler. This make
	// sure the page is only accessible by moderators
	// We don't need CSRF here
	// TODO Find a native mux way to add a 'prehook' for route /mod
	modRoutes := Router.Group("/mod", errorMiddleware(), modMiddleware())
	{
		modRoutes.Any("/mod", IndexModPanel)
		modRoutes.Any("/mod/torrents", TorrentsListPanel)
		modRoutes.Any("/mod/torrents/:page", TorrentsListPanel)
		modRoutes.POST("/mod/torrents", TorrentsPostListPanel)
		modRoutes.POST("/mod/torrents/:page", TorrentsPostListPanel)
		modRoutes.Any("/mod/torrents/deleted", DeletedTorrentsModPanel)
		modRoutes.Any("/mod/torrents/deleted/:page", DeletedTorrentsModPanel)
		modRoutes.Any("/mod/torrents/deleted", DeletedTorrentsPostPanel)
		modRoutes.Any("/mod/torrents/deleted/:page", DeletedTorrentsPostPanel)
		modRoutes.Any("/mod/reports", TorrentReportListPanel)
		modRoutes.Any("/mod/reports/{page}", TorrentReportListPanel)
		modRoutes.Any("/mod/users", UsersListPanel)
		modRoutes.Any("/mod/users/{page}", UsersListPanel)
		modRoutes.Any("/mod/comments", CommentsListPanel)
		modRoutes.Any("/mod/comments/{page}", CommentsListPanel)
		modRoutes.Any("/mod/comment", CommentsListPanel) // TODO
		modRoutes.Any("/mod/torrent/", TorrentEditModPanel)
		modRoutes.Any("/mod/torrent/", TorrentPostEditModPanel)
		modRoutes.Any("/mod/torrent/delete", TorrentDeleteModPanel)
		modRoutes.Any("/mod/torrent/block", TorrentBlockModPanel)
		modRoutes.Any("/mod/report/delete", TorrentReportDeleteModPanel)
		modRoutes.Any("/mod/comment/delete", CommentDeleteModPanel)
		modRoutes.Any("/mod/reassign", TorrentReassignModPanel)
		modRoutes.Any("/mod/reassign", TorrentPostReassignModPanel)
		apiMod := modRoutes.Group("/mod/api")
		apiMod.Any("/torrents", APIMassMod)
	}
	//reporting a torrent
	Router.POST("/report/:id", ReportTorrentHandler)

	Router.Any("/captcha", captcha.ServeFiles)

	Router.Any("/dumps", DatabaseDumpHandler)

	Router.GET("/settings", SeePublicSettingsHandler)
	Router.POST("/settings", ChangePublicSettingsHandler)

	Router.Use(errorMiddleware())

	CSRFRouter = nosurf.New(Router)
	CSRFRouter.ExemptRegexp("/api(?:/.+)*")
	CSRFRouter.ExemptPath("/mod")
	CSRFRouter.ExemptPath("/upload")
	CSRFRouter.ExemptPath("/user/login")
}
