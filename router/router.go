package router

import (
	"net/http"

	"github.com/NyaaPantsu/nyaa/service/captcha"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var Router *mux.Router

func init() {
	// Static file handlers
	cssHandler := http.FileServer(http.Dir("./public/css/"))
	jsHandler := http.FileServer(http.Dir("./public/js/"))
	imgHandler := http.FileServer(http.Dir("./public/img/"))
	// TODO Use config from cli
	// TODO Make sure the directory exists
	dumpsHandler  := http.FileServer(http.Dir(DatabaseDumpPath))
	// TODO Use config from cli
	// TODO Make sure the directory exists
	gpgKeyHandler := http.FileServer(http.Dir(GPGPublicKeyPath))
	gzipHomeHandler := http.HandlerFunc(HomeHandler)
	gzipAPIHandler := http.HandlerFunc(ApiHandler)
	gzipAPIViewHandler := http.HandlerFunc(ApiViewHandler)
	gzipViewHandler := http.HandlerFunc(ViewHandler)
	gzipUserProfileHandler := http.HandlerFunc(UserProfileHandler)
	gzipUserDetailsHandler := http.HandlerFunc(UserDetailsHandler)
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
	Router.Handle("/", gzipHomeHandler).Name("home")
	Router.Handle("/page/{page:[0-9]+}", wrapHandler(gzipHomeHandler)).Name("home_page")
	Router.HandleFunc("/search", SearchHandler).Name("search")
	Router.HandleFunc("/search/{page}", SearchHandler).Name("search_page")
	Router.Handle("/api", wrapHandler(gzipAPIHandler)).Methods("GET")
	Router.Handle("/api/{page:[0-9]*}", wrapHandler(gzipAPIHandler)).Methods("GET")
	Router.Handle("/api/view/{id}", wrapHandler(gzipAPIViewHandler)).Methods("GET")
	Router.HandleFunc("/api/view/{id}", ApiViewHeadHandler).Methods("HEAD")
	Router.HandleFunc("/api/upload", ApiUploadHandler).Methods("POST")
	Router.HandleFunc("/api/update", ApiUpdateHandler).Methods("PUT")
	Router.HandleFunc("/faq", FaqHandler).Name("faq")
	Router.HandleFunc("/feed", RSSHandler).Name("feed")
	Router.HandleFunc("/feed/{page}", RSSHandler).Name("feed_page")
	Router.Handle("/view/{id}", wrapHandler(gzipViewHandler)).Methods("GET").Name("view_torrent")
	Router.HandleFunc("/view/{id}", ViewHeadHandler).Methods("HEAD")
	Router.HandleFunc("/view/{id}", PostCommentHandler).Methods("POST").Name("post_comment")
	Router.HandleFunc("/upload", UploadHandler).Name("upload")
	Router.HandleFunc("/user/register", UserRegisterFormHandler).Name("user_register").Methods("GET")
	Router.HandleFunc("/user/login", UserLoginFormHandler).Name("user_login").Methods("GET")
	Router.HandleFunc("/verify/email/{token}", UserVerifyEmailHandler).Name("user_verify").Methods("GET")
	Router.HandleFunc("/user/register", UserRegisterPostHandler).Name("user_register").Methods("POST")
	Router.HandleFunc("/user/login", UserLoginPostHandler).Name("user_login").Methods("POST")
	Router.HandleFunc("/user/logout", UserLogoutHandler).Name("user_logout")
	Router.Handle("/user/{id}/{username}", wrapHandler(gzipUserProfileHandler)).Name("user_profile").Methods("GET")
	Router.HandleFunc("/user/{id}/{username}/follow", UserFollowHandler).Name("user_follow").Methods("GET")
	Router.Handle("/user/{id}/{username}/edit", wrapHandler(gzipUserDetailsHandler)).Name("user_profile_details").Methods("GET")
	Router.Handle("/user/{id}/{username}/edit", wrapHandler(gzipUserProfileFormHandler)).Name("user_profile_edit").Methods("POST")
	Router.Handle("/user/notifications", wrapHandler(gzipUserNotificationsHandler)).Name("user_notifications")
	Router.HandleFunc("/user/{id}/{username}/feed", RSSHandler).Name("feed_user")
	Router.HandleFunc("/user/{id}/{username}/feed/{page}", RSSHandler).Name("feed_user_page")

	// INFO Everything under /mod should be wrapped by WrapModHandler. This make
	// sure the page is only accessible by moderators
	// TODO Find a native mux way to add a 'prehook' for route /mod
	Router.HandleFunc("/mod",                 WrapModHandler(IndexModPanel)).Name("mod_index")
	Router.HandleFunc("/mod/torrents",        WrapModHandler(TorrentsListPanel)).Name("mod_tlist").Methods("GET")
	Router.HandleFunc("/mod/torrents/{page}", WrapModHandler(TorrentsListPanel)).Name("mod_tlist_page").Methods("GET")
	Router.HandleFunc("/mod/torrents", WrapModHandler(TorrentsPostListPanel)).Methods("POST")
	Router.HandleFunc("/mod/torrents/{page}", WrapModHandler(TorrentsPostListPanel)).Methods("POST")
	Router.HandleFunc("/mod/reports",         WrapModHandler(TorrentReportListPanel)).Name("mod_trlist")
	Router.HandleFunc("/mod/reports/{page}",  WrapModHandler(TorrentReportListPanel)).Name("mod_trlist_page")
	Router.HandleFunc("/mod/users",           WrapModHandler(UsersListPanel)).Name("mod_ulist")
	Router.HandleFunc("/mod/users/{page}",    WrapModHandler(UsersListPanel)).Name("mod_ulist_page")
	Router.HandleFunc("/mod/comments",        WrapModHandler(CommentsListPanel)).Name("mod_clist")
	Router.HandleFunc("/mod/comments/{page}", WrapModHandler(CommentsListPanel)).Name("mod_clist_page")
	Router.HandleFunc("/mod/comment",         WrapModHandler(CommentsListPanel)).Name("mod_cedit") // TODO
	Router.HandleFunc("/mod/torrent/",        WrapModHandler(TorrentEditModPanel)).Name("mod_tedit").Methods("GET")
	Router.HandleFunc("/mod/torrent/",        WrapModHandler(TorrentPostEditModPanel)).Name("mod_ptedit").Methods("POST")
	Router.HandleFunc("/mod/torrent/delete",  WrapModHandler(TorrentDeleteModPanel)).Name("mod_tdelete")
	Router.HandleFunc("/mod/report/delete",   WrapModHandler(TorrentReportDeleteModPanel)).Name("mod_trdelete")
	Router.HandleFunc("/mod/comment/delete",  WrapModHandler(CommentDeleteModPanel)).Name("mod_cdelete")
	Router.HandleFunc("/mod/reassign",        WrapModHandler(TorrentReassignModPanel)).Name("mod_treassign").Methods("GET")
	Router.HandleFunc("/mod/reassign",        WrapModHandler(TorrentPostReassignModPanel)).Name("mod_treassign").Methods("POST")

	//reporting a torrent
	Router.HandleFunc("/report/{id}", ReportTorrentHandler).Methods("POST").Name("torrent_report")

	Router.PathPrefix("/captcha").Methods("GET").HandlerFunc(captcha.ServeFiles)

	Router.Handle("/dumps", gzipDatabaseDumpHandler).Name("dump").Methods("GET")

	Router.HandleFunc("/language", SeeLanguagesHandler).Methods("GET").Name("see_languages")
	Router.HandleFunc("/language", ChangeLanguageHandler).Methods("POST").Name("change_language")

	Router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
}
