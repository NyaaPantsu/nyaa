package router

import (
	"net/http"

	"github.com/ewhal/nyaa/service/captcha"
	// "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var Router *mux.Router

func init() {
	// Static file handlers
	cssHandler := http.FileServer(http.Dir("./public/css/"))
	jsHandler := http.FileServer(http.Dir("./public/js/"))
	imgHandler := http.FileServer(http.Dir("./public/img/"))
	gzipHomeHandler := http.HandlerFunc(HomeHandler)
	gzipAPIHandler := http.HandlerFunc(ApiHandler)
	gzipAPIViewHandler := http.HandlerFunc(ApiViewHandler)
	gzipViewHandler := http.HandlerFunc(ViewHandler)
	gzipUserProfileHandler := http.HandlerFunc(UserProfileHandler)
	gzipUserDetailsHandler := http.HandlerFunc(UserDetailsHandler)
	gzipUserProfileFormHandler := http.HandlerFunc(UserProfileFormHandler)
/*
	// Enable GZIP compression for all handlers except imgHandler and captcha
	gzipCSSHandler := cssHandler)
	gzipJSHandler:= jsHandler)
	gzipSearchHandler:= http.HandlerFunc(SearchHandler)
	gzipAPIUploadHandler := http.HandlerFunc(ApiUploadHandler)
	gzipAPIUpdateHandler := http.HandlerFunc(ApiUpdateHandler)
	gzipFaqHandler := http.HandlerFunc(FaqHandler)
	gzipRSSHandler := http.HandlerFunc(RSSHandler)
	gzipUploadHandler := http.HandlerFunc(UploadHandler)
	gzipUserRegisterFormHandler := http.HandlerFunc(UserRegisterFormHandler)
	gzipUserLoginFormHandler := http.HandlerFunc(UserLoginFormHandler)
	gzipUserVerifyEmailHandler := http.HandlerFunc(UserVerifyEmailHandler)
	gzipUserRegisterPostHandler := http.HandlerFunc(UserRegisterPostHandler)
	gzipUserLoginPostHandler := http.HandlerFunc(UserLoginPostHandler)
	gzipUserLogoutHandler := http.HandlerFunc(UserLogoutHandler)
	gzipUserFollowHandler := http.HandlerFunc(UserFollowHandler)

	gzipIndexModPanel := http.HandlerFunc(IndexModPanel)
	gzipTorrentsListPanel := http.HandlerFunc(TorrentsListPanel)
	gzipTorrentReportListPanel := http.HandlerFunc(TorrentReportListPanel)
	gzipUsersListPanel := http.HandlerFunc(UsersListPanel)
	gzipCommentsListPanel := http.HandlerFunc(CommentsListPanel)
	gzipTorrentEditModPanel := http.HandlerFunc(TorrentEditModPanel)
	gzipTorrentPostEditModPanel := http.HandlerFunc(TorrentPostEditModPanel)
	gzipCommentDeleteModPanel := http.HandlerFunc(CommentDeleteModPanel)
	gzipTorrentDeleteModPanel := http.HandlerFunc(TorrentDeleteModPanel)
	gzipTorrentReportDeleteModPanel := http.HandlerFunc(TorrentReportDeleteModPanel)*/

	//gzipTorrentReportCreateHandler := http.HandlerFunc(CreateTorrentReportHandler)
	//gzipTorrentReportDeleteHandler := http.HandlerFunc(DeleteTorrentReportHandler)
	//gzipTorrentDeleteHandler := http.HandlerFunc(DeleteTorrentHandler)

	Router = mux.NewRouter()

	// Routes
	http.Handle("/css/", http.StripPrefix("/css/", cssHandler))
	http.Handle("/js/", http.StripPrefix("/js/", jsHandler))
	http.Handle("/img/", http.StripPrefix("/img/", imgHandler))
	Router.Handle("/", wrapHandler(gzipHomeHandler)).Name("home")
	Router.Handle("/page/{page:[0-9]+}", wrapHandler(gzipHomeHandler)).Name("home_page")
	Router.HandleFunc("/search", SearchHandler).Name("search")
	Router.HandleFunc("/search/{page}", SearchHandler).Name("search_page")
	Router.Handle("/api", wrapHandler(gzipAPIHandler)).Methods("GET")
	Router.Handle("/api/{page:[0-9]*}", wrapHandler(gzipAPIHandler)).Methods("GET")
	Router.Handle("/api/view/{id}", wrapHandler(gzipAPIViewHandler)).Methods("GET")
	Router.HandleFunc("/api/upload", ApiUploadHandler).Methods("POST")
	Router.HandleFunc("/api/update", ApiUpdateHandler).Methods("PUT")
	Router.HandleFunc("/faq", FaqHandler).Name("faq")
	Router.HandleFunc("/feed", RSSHandler).Name("feed")
	Router.Handle("/view/{id}", wrapHandler(gzipViewHandler)).Methods("GET").Name("view_torrent")
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

	Router.HandleFunc("/mod", IndexModPanel).Name("mod_index")
	Router.HandleFunc("/mod/torrents", TorrentsListPanel).Name("mod_tlist")
	Router.HandleFunc("/mod/torrents/{page}", TorrentsListPanel).Name("mod_tlist_page")
	Router.HandleFunc("/mod/reports", TorrentReportListPanel).Name("mod_trlist")
	Router.HandleFunc("/mod/reports/{page}", TorrentReportListPanel).Name("mod_trlist_page")
	Router.HandleFunc("/mod/users", UsersListPanel).Name("mod_ulist")
	Router.HandleFunc("/mod/users/{page}", UsersListPanel).Name("mod_ulist_page")
	Router.HandleFunc("/mod/comments", CommentsListPanel).Name("mod_clist")
	Router.HandleFunc("/mod/comments/{page}", CommentsListPanel).Name("mod_clist_page")
	Router.HandleFunc("/mod/comment", CommentsListPanel).Name("mod_cedit") // TODO
	Router.HandleFunc("/mod/torrent/", TorrentEditModPanel).Name("mod_tedit").Methods("GET")
	Router.HandleFunc("/mod/torrent/", TorrentPostEditModPanel).Name("mod_ptedit").Methods("POST")
	Router.HandleFunc("/mod/torrent/delete", TorrentDeleteModPanel).Name("mod_tdelete")
	Router.HandleFunc("/mod/report/delete", TorrentReportDeleteModPanel).Name("mod_trdelete")
	Router.HandleFunc("/mod/comment/delete", CommentDeleteModPanel).Name("mod_cdelete")

	//reporting a torrent
	Router.HandleFunc("/report/{id}", ReportTorrentHandler).Methods("POST").Name("post_comment")

	Router.PathPrefix("/captcha").Methods("GET").HandlerFunc(captcha.ServeFiles)

	//Router.HandleFunc("/report/create", gzipTorrentReportCreateHandler).Name("torrent_report_create").Methods("POST")
	// TODO Allow only moderators to access /moderation/*
	//Router.HandleFunc("/moderation/report/delete", gzipTorrentReportDeleteHandler).Name("torrent_report_delete").Methods("POST")
	//Router.HandleFunc("/moderation/torrent/delete", gzipTorrentDeleteHandler).Name("torrent_delete").Methods("POST")

	Router.HandleFunc("/language", SeeLanguagesHandler).Methods("GET").Name("see_languages")
	Router.HandleFunc("/language", ChangeLanguageHandler).Methods("POST").Name("change_language")

	Router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
}
