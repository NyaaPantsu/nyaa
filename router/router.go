package router

import (
	"net/http"

	"github.com/ewhal/nyaa/service/captcha"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var Router *mux.Router

func init() {
	// Static file handlers
	cssHandler := http.FileServer(http.Dir("./public/css/"))
	jsHandler := http.FileServer(http.Dir("./public/js/"))
	imgHandler := http.FileServer(http.Dir("./public/img/"))

	// Enable GZIP compression for all handlers except imgHandler and captcha
	gzipCSSHandler := handlers.CompressHandler(cssHandler)
	gzipJSHandler := handlers.CompressHandler(jsHandler)
	gzipHomeHandler := handlers.CompressHandler(http.HandlerFunc(HomeHandler))
	gzipSearchHandler := handlers.CompressHandler(http.HandlerFunc(SearchHandler))
	gzipAPIHandler := handlers.CompressHandler(http.HandlerFunc(ApiHandler))
	gzipAPIViewHandler := handlers.CompressHandler(http.HandlerFunc(ApiViewHandler))
	gzipAPIUploadHandler := handlers.CompressHandler(http.HandlerFunc(ApiUploadHandler))
	gzipAPIUpdateHandler := handlers.CompressHandler(http.HandlerFunc(ApiUpdateHandler))
	gzipFaqHandler := handlers.CompressHandler(http.HandlerFunc(FaqHandler))
	gzipRSSHandler := handlers.CompressHandler(http.HandlerFunc(RSSHandler))
	gzipViewHandler := handlers.CompressHandler(http.HandlerFunc(ViewHandler))
	gzipUploadHandler := handlers.CompressHandler(http.HandlerFunc(UploadHandler))
	gzipUserRegisterFormHandler := handlers.CompressHandler(http.HandlerFunc(UserRegisterFormHandler))
	gzipUserLoginFormHandler := handlers.CompressHandler(http.HandlerFunc(UserLoginFormHandler))
	gzipUserVerifyEmailHandler := handlers.CompressHandler(http.HandlerFunc(UserVerifyEmailHandler))
	gzipUserRegisterPostHandler := handlers.CompressHandler(http.HandlerFunc(UserRegisterPostHandler))
	gzipUserLoginPostHandler := handlers.CompressHandler(http.HandlerFunc(UserLoginPostHandler))
	gzipUserLogoutHandler := handlers.CompressHandler(http.HandlerFunc(UserLogoutHandler))
	gzipUserProfileHandler := handlers.CompressHandler(http.HandlerFunc(UserProfileHandler))
	gzipUserFollowHandler := handlers.CompressHandler(http.HandlerFunc(UserFollowHandler))
	gzipUserDetailsHandler := handlers.CompressHandler(http.HandlerFunc(UserDetailsHandler))
	gzipUserProfileFormHandler := handlers.CompressHandler(http.HandlerFunc(UserProfileFormHandler))

	gzipIndexModPanel := handlers.CompressHandler(http.HandlerFunc(IndexModPanel))
	gzipTorrentsListPanel := handlers.CompressHandler(http.HandlerFunc(TorrentsListPanel))
	gzipTorrentReportListPanel := handlers.CompressHandler(http.HandlerFunc(TorrentReportListPanel))
	gzipUsersListPanel := handlers.CompressHandler(http.HandlerFunc(UsersListPanel))
	gzipCommentsListPanel := handlers.CompressHandler(http.HandlerFunc(CommentsListPanel))
	gzipTorrentEditModPanel := handlers.CompressHandler(http.HandlerFunc(TorrentEditModPanel))
	gzipTorrentPostEditModPanel := handlers.CompressHandler(http.HandlerFunc(TorrentPostEditModPanel))
	gzipCommentDeleteModPanel := handlers.CompressHandler(http.HandlerFunc(CommentDeleteModPanel))
	gzipTorrentDeleteModPanel := handlers.CompressHandler(http.HandlerFunc(TorrentDeleteModPanel))

	//gzipTorrentReportCreateHandler := handlers.CompressHandler(http.HandlerFunc(CreateTorrentReportHandler))
	//gzipTorrentReportDeleteHandler := handlers.CompressHandler(http.HandlerFunc(DeleteTorrentReportHandler))
	//gzipTorrentDeleteHandler := handlers.CompressHandler(http.HandlerFunc(DeleteTorrentHandler))

	Router = mux.NewRouter()

	// Routes
	http.Handle("/css/", http.StripPrefix("/css/", wrapHandler(gzipCSSHandler)))
	http.Handle("/js/", http.StripPrefix("/js/", wrapHandler(gzipJSHandler)))
	http.Handle("/img/", http.StripPrefix("/img/", wrapHandler(imgHandler)))
	Router.Handle("/", gzipHomeHandler).Name("home")
	Router.Handle("/page/{page:[0-9]+}", wrapHandler(gzipHomeHandler)).Name("home_page")
	Router.Handle("/search", gzipSearchHandler).Name("search")
	Router.Handle("/search/{page}", gzipSearchHandler).Name("search_page")
	Router.Handle("/api", gzipAPIHandler).Methods("GET")
	Router.Handle("/api/{page:[0-9]*}", wrapHandler(gzipAPIHandler)).Methods("GET")
	Router.Handle("/api/view/{id}", wrapHandler(gzipAPIViewHandler)).Methods("GET")
	Router.Handle("/api/upload", gzipAPIUploadHandler).Methods("POST")
	Router.Handle("/api/update", gzipAPIUpdateHandler).Methods("PUT")
	Router.Handle("/faq", gzipFaqHandler).Name("faq")
	Router.Handle("/feed", gzipRSSHandler).Name("feed")
	Router.Handle("/view/{id}", wrapHandler(gzipViewHandler)).Methods("GET").Name("view_torrent")
	Router.HandleFunc("/view/{id}", PostCommentHandler).Methods("POST").Name("post_comment")
	Router.Handle("/upload", gzipUploadHandler).Name("upload")
	Router.Handle("/user/register", gzipUserRegisterFormHandler).Name("user_register").Methods("GET")
	Router.Handle("/user/login", gzipUserLoginFormHandler).Name("user_login").Methods("GET")
	Router.Handle("/verify/email/{token}", gzipUserVerifyEmailHandler).Name("user_verify").Methods("GET")
	Router.Handle("/user/register", gzipUserRegisterPostHandler).Name("user_register").Methods("POST")
	Router.Handle("/user/login", gzipUserLoginPostHandler).Name("user_login").Methods("POST")
	Router.Handle("/user/logout", gzipUserLogoutHandler).Name("user_logout")
	Router.Handle("/user/{id}/{username}", wrapHandler(gzipUserProfileHandler)).Name("user_profile").Methods("GET")
	Router.Handle("/user/{id}/{username}/follow", gzipUserFollowHandler).Name("user_follow").Methods("GET")
	Router.Handle("/user/{id}/{username}/edit", wrapHandler(gzipUserDetailsHandler)).Name("user_profile_details").Methods("GET")
	Router.Handle("/user/{id}/{username}/edit", wrapHandler(gzipUserProfileFormHandler)).Name("user_profile_edit").Methods("POST")

	Router.Handle("/mod", gzipIndexModPanel).Name("mod_index")
	Router.Handle("/mod/torrents", gzipTorrentsListPanel).Name("mod_tlist")
	Router.Handle("/mod/torrents/{page}", gzipTorrentsListPanel).Name("mod_tlist_page")
	Router.Handle("/mod/users", gzipUsersListPanel).Name("mod_ulist")
	Router.Handle("/mod/users/{page}", gzipUsersListPanel).Name("mod_ulist_page")
	Router.Handle("/mod/comments", gzipCommentsListPanel).Name("mod_clist")
	Router.Handle("/mod/comments/{page}", gzipCommentsListPanel).Name("mod_clist_page")
	Router.Handle("/mod/comment", gzipCommentsListPanel).Name("mod_cedit") // TODO
	Router.Handle("/mod/torrent/", gzipTorrentEditModPanel).Name("mod_tedit")
	Router.Handle("/mod/torrent/", gzipTorrentPostEditModPanel).Name("mod_ptedit")
	Router.Handle("/mod/torrent/delete", gzipTorrentDeleteModPanel).Name("mod_tdelete")
	Router.Handle("/mod/comment/delete", gzipCommentDeleteModPanel).Name("mod_cdelete")

	//reporting a torrent
	Router.HandleFunc("/report/{id}", ReportTorrentHandler).Methods("POST").Name("post_comment")

	Router.PathPrefix("/captcha").Methods("GET").HandlerFunc(captcha.ServeFiles)

	//Router.Handle("/report/create", gzipTorrentReportCreateHandler).Name("torrent_report_create").Methods("POST")
	// TODO Allow only moderators to access /moderation/*
	//Router.Handle("/moderation/report/delete", gzipTorrentReportDeleteHandler).Name("torrent_report_delete").Methods("POST")
	//Router.Handle("/moderation/torrent/delete", gzipTorrentDeleteHandler).Name("torrent_delete").Methods("POST")
	Router.Handle("/mod/reports", gzipTorrentReportListPanel).Name("mod_trlist").Methods("GET")
	Router.Handle("/mod/reports/{page}", gzipTorrentReportListPanel).Name("mod_trlist_page").Methods("GET")

	Router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
}
