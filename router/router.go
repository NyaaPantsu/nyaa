package router

import (
	"net/http"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/service/captcha"
	"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Router variable for exporting the route configuration
var Router *mux.Router

func init() {
	// Static file handlers
	cssHandler := http.FileServer(http.Dir("./public/css/"))
	jsHandler := http.FileServer(http.Dir("./public/js/"))
	imgHandler := http.FileServer(http.Dir("./public/img/"))
	// TODO Use config from cli
	// TODO Make sure the directory exists
	dumpsHandler := http.FileServer(http.Dir(DatabaseDumpPath))
	// TODO Use config from cli
	// TODO Make sure the directory exists
	gpgKeyHandler := http.FileServer(http.Dir(GPGPublicKeyPath))
	gzipHomeHandler := http.HandlerFunc(HomeHandler)
	gzipAPIHandler := http.HandlerFunc(APIHandler)
	gzipAPIViewHandler := http.HandlerFunc(APIViewHandler)
	gzipViewHandler := http.HandlerFunc(ViewHandler)
	gzipUserProfileHandler := http.HandlerFunc(UserProfileHandler)
	gzipUserAPIKeyResetHandler := http.HandlerFunc(UserAPIKeyResetHandler)
	gzipUserDetailsHandler := http.HandlerFunc(UserDetailsHandler)
	downloadTorrentHandler := http.HandlerFunc(DownloadTorrent)
	gzipUserProfileFormHandler := http.HandlerFunc(UserProfileFormHandler)
	gzipUserNotificationsHandler := http.HandlerFunc(UserNotificationsHandler)
	gzipDumpsHandler := handlers.CompressHandler(dumpsHandler)
	gzipGpgKeyHandler := handlers.CompressHandler(gpgKeyHandler)
	gzipDatabaseDumpHandler := handlers.CompressHandler(http.HandlerFunc(DatabaseDumpHandler))

	Router = mux.NewRouter()
	http.Handle("/css/", http.StripPrefix("/css/", cssHandler))
	http.Handle("/js/", http.StripPrefix("/js/", jsHandler))
	http.Handle("/img/", http.StripPrefix("/img/", imgHandler))
	http.Handle("/dbdumps/", http.StripPrefix("/dbdumps/", wrapHandler(gzipDumpsHandler)))
	http.Handle("/gpg/", http.StripPrefix("/gpg/", wrapHandler(gzipGpgKeyHandler)))

	// We don't need CSRF here
	Router.Handle("/", gzipHomeHandler).Name("home")
	Router.Handle("/page/{page:[0-9]+}", wrapHandler(gzipHomeHandler)).Name("home_page")
	Router.HandleFunc("/search", SearchHandler).Name("search")
	Router.HandleFunc("/search/{page}", SearchHandler).Name("search_page")
	Router.HandleFunc("/verify/email/{token}", UserVerifyEmailHandler).Name("user_verify").Methods("GET")
	Router.HandleFunc("/faq", FaqHandler).Name("faq")
	Router.HandleFunc("/feed", RSSHandler).Name("feed")
	Router.HandleFunc("/feed/{page}", RSSHandler).Name("feed_page")

	// !!! This line need to have the same download location as the one define in config.TorrentStorageLink !!!
	Router.Handle("/download/{hash}", wrapHandler(downloadTorrentHandler)).Name("torrent_download")

	// For now, no CSRF protection here, as API is not usable for uploads
	Router.HandleFunc("/upload", UploadHandler).Name("upload")
	Router.HandleFunc("/user/login", UserLoginPostHandler).Name("user_login").Methods("POST")

	torrentViewRoutes := Router.PathPrefix("/view").Subrouter()
	torrentViewRoutes.Handle("/{id}", wrapHandler(gzipViewHandler)).Methods("GET").Name("view_torrent")
	torrentViewRoutes.HandleFunc("/{id}", ViewHeadHandler).Methods("HEAD")
	torrentViewRoutes.HandleFunc("/{id}", PostCommentHandler).Methods("POST").Name("post_comment")

	torrentRoutes := Router.PathPrefix("/torrent").Subrouter()
	torrentRoutes.HandleFunc("/", TorrentEditUserPanel).Methods("GET").Name("user_torrent_edit")
	torrentRoutes.HandleFunc("/", TorrentPostEditUserPanel).Methods("POST").Name("user_torrent_edit")
	torrentRoutes.HandleFunc("/delete", TorrentDeleteUserPanel).Methods("GET").Name("user_torrent_delete")

	userRoutes := Router.PathPrefix("/user").Subrouter()
	userRoutes.HandleFunc("/register", UserRegisterFormHandler).Name("user_register").Methods("GET")
	userRoutes.HandleFunc("/login", UserLoginFormHandler).Name("user_login").Methods("GET")
	userRoutes.HandleFunc("/register", UserRegisterPostHandler).Name("user_register").Methods("POST")
	userRoutes.HandleFunc("/logout", UserLogoutHandler).Name("user_logout")
	userRoutes.Handle("/{id}/{username}", wrapHandler(gzipUserProfileHandler)).Name("user_profile").Methods("GET")
	userRoutes.HandleFunc("/{id}/{username}/follow", UserFollowHandler).Name("user_follow").Methods("GET")
	userRoutes.Handle("/{id}/{username}/edit", wrapHandler(gzipUserDetailsHandler)).Name("user_profile_details").Methods("GET")
	userRoutes.Handle("/{id}/{username}/edit", wrapHandler(gzipUserProfileFormHandler)).Name("user_profile_edit").Methods("POST")
	userRoutes.Handle("/{id}/{username}/apireset", wrapHandler(gzipUserAPIKeyResetHandler)).Name("user_profile_apireset").Methods("GET")
	userRoutes.Handle("/notifications", wrapHandler(gzipUserNotificationsHandler)).Name("user_notifications")
	userRoutes.HandleFunc("/{id}/{username}/feed", RSSHandler).Name("feed_user")
	userRoutes.HandleFunc("/{id}/{username}/feed/{page}", RSSHandler).Name("feed_user_page")

	// Please make EnableSecureCSRF to false when testing locally
	if config.EnableSecureCSRF {
		userRoutes.Handle("/", csrf.Protect(config.CSRFTokenHashKey)(userRoutes))
		torrentRoutes.Handle("/", csrf.Protect(config.CSRFTokenHashKey)(torrentRoutes))
		torrentViewRoutes.Handle("/", csrf.Protect(config.CSRFTokenHashKey)(torrentViewRoutes))
	} else {
		userRoutes.Handle("/", csrf.Protect(config.CSRFTokenHashKey, csrf.Secure(false))(userRoutes))
		torrentRoutes.Handle("/", csrf.Protect(config.CSRFTokenHashKey, csrf.Secure(false))(torrentRoutes))
		torrentViewRoutes.Handle("/", csrf.Protect(config.CSRFTokenHashKey, csrf.Secure(false))(torrentViewRoutes))
	}

	// We don't need CSRF here
	api := Router.PathPrefix("/api").Subrouter()
	api.Handle("", wrapHandler(gzipAPIHandler)).Methods("GET")
	api.Handle("/", wrapHandler(gzipAPIHandler)).Methods("GET")
	api.Handle("/{page:[0-9]*}", wrapHandler(gzipAPIHandler)).Methods("GET")
	api.Handle("/view/{id}", wrapHandler(gzipAPIViewHandler)).Methods("GET")
	api.HandleFunc("/view/{id}", APIViewHeadHandler).Methods("HEAD")
	api.HandleFunc("/upload", APIUploadHandler).Methods("POST")
	api.HandleFunc("/search", APISearchHandler)
	api.HandleFunc("/search/{page}", APISearchHandler)
	api.HandleFunc("/update", APIUpdateHandler).Methods("PUT")

	// INFO Everything under /mod should be wrapped by wrapModHandler. This make
	// sure the page is only accessible by moderators
	// We don't need CSRF here
	// TODO Find a native mux way to add a 'prehook' for route /mod
	Router.HandleFunc("/mod", wrapModHandler(IndexModPanel)).Name("mod_index")
	Router.HandleFunc("/mod/torrents", wrapModHandler(TorrentsListPanel)).Name("mod_tlist").Methods("GET")
	Router.HandleFunc("/mod/torrents/{page:[0-9]+}", wrapModHandler(TorrentsListPanel)).Name("mod_tlist_page").Methods("GET")
	Router.HandleFunc("/mod/torrents", wrapModHandler(TorrentsPostListPanel)).Methods("POST")
	Router.HandleFunc("/mod/torrents/{page:[0-9]+}", wrapModHandler(TorrentsPostListPanel)).Methods("POST")
	Router.HandleFunc("/mod/torrents/deleted", wrapModHandler(DeletedTorrentsModPanel)).Name("mod_tlist_deleted").Methods("GET")
	Router.HandleFunc("/mod/torrents/deleted/{page:[0-9]+}", wrapModHandler(DeletedTorrentsModPanel)).Name("mod_tlist_deleted_page").Methods("GET")
	Router.HandleFunc("/mod/torrents/deleted", wrapModHandler(DeletedTorrentsPostPanel)).Name("mod_tlist_deleted").Methods("POST")
	Router.HandleFunc("/mod/torrents/deleted/{page:[0-9]+}", wrapModHandler(DeletedTorrentsPostPanel)).Name("mod_tlist_deleted_page").Methods("POST")
	Router.HandleFunc("/mod/reports", wrapModHandler(TorrentReportListPanel)).Name("mod_trlist")
	Router.HandleFunc("/mod/reports/{page}", wrapModHandler(TorrentReportListPanel)).Name("mod_trlist_page")
	Router.HandleFunc("/mod/users", wrapModHandler(UsersListPanel)).Name("mod_ulist")
	Router.HandleFunc("/mod/users/{page}", wrapModHandler(UsersListPanel)).Name("mod_ulist_page")
	Router.HandleFunc("/mod/comments", wrapModHandler(CommentsListPanel)).Name("mod_clist")
	Router.HandleFunc("/mod/comments/{page}", wrapModHandler(CommentsListPanel)).Name("mod_clist_page")
	Router.HandleFunc("/mod/comment", wrapModHandler(CommentsListPanel)).Name("mod_cedit") // TODO
	Router.HandleFunc("/mod/torrent/", wrapModHandler(TorrentEditModPanel)).Name("mod_tedit").Methods("GET")
	Router.HandleFunc("/mod/torrent/", wrapModHandler(TorrentPostEditModPanel)).Name("mod_ptedit").Methods("POST")
	Router.HandleFunc("/mod/torrent/delete", wrapModHandler(TorrentDeleteModPanel)).Name("mod_tdelete")
	Router.HandleFunc("/mod/torrent/block", wrapModHandler(TorrentBlockModPanel)).Name("mod_tblock")
	Router.HandleFunc("/mod/report/delete", wrapModHandler(TorrentReportDeleteModPanel)).Name("mod_trdelete")
	Router.HandleFunc("/mod/comment/delete", wrapModHandler(CommentDeleteModPanel)).Name("mod_cdelete")
	Router.HandleFunc("/mod/reassign", wrapModHandler(TorrentReassignModPanel)).Name("mod_treassign").Methods("GET")
	Router.HandleFunc("/mod/reassign", wrapModHandler(TorrentPostReassignModPanel)).Name("mod_treassign").Methods("POST")

	apiMod := Router.PathPrefix("/mod/api").Subrouter()
	apiMod.HandleFunc("/torrents", wrapModHandler(APIMassMod)).Name("mod_tapi").Methods("POST")

	//reporting a torrent
	Router.HandleFunc("/report/{id}", ReportTorrentHandler).Methods("POST").Name("torrent_report")

	Router.PathPrefix("/captcha").Methods("GET").HandlerFunc(captcha.ServeFiles)

	Router.Handle("/dumps", gzipDatabaseDumpHandler).Name("dump").Methods("GET")

	Router.HandleFunc("/settings", SeePublicSettingsHandler).Methods("GET").Name("see_languages")
	Router.HandleFunc("/settings", ChangePublicSettingsHandler).Methods("POST").Name("see_languages")

	Router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
}
