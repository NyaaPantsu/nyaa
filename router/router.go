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
	gzipUserProfileFormHandler := handlers.CompressHandler(http.HandlerFunc(UserProfileFormHandler))

	Router = mux.NewRouter()

	// Routes
	http.Handle("/css/", http.StripPrefix("/css/", gzipCSSHandler))
	http.Handle("/js/", http.StripPrefix("/js/", gzipJSHandler))
	http.Handle("/img/", http.StripPrefix("/img/", imgHandler))
	Router.Handle("/", gzipHomeHandler).Name("home")
	Router.Handle("/page/{page:[0-9]+}", gzipHomeHandler).Name("home_page")
	Router.Handle("/search", gzipSearchHandler).Name("search")
	Router.Handle("/search/{page}", gzipSearchHandler).Name("search_page")
	Router.Handle("/api", gzipAPIHandler).Methods("GET")
	Router.Handle("/api/{page:[0-9]*}", gzipAPIHandler).Methods("GET")
	Router.Handle("/api/view/{id}", gzipAPIViewHandler).Methods("GET")
	Router.Handle("/api/upload", gzipAPIUploadHandler).Methods("POST")
	Router.Handle("/api/update", gzipAPIUpdateHandler).Methods("PUT")
	Router.Handle("/faq", gzipFaqHandler).Name("faq")
	Router.Handle("/feed", gzipRSSHandler).Name("feed")
	Router.Handle("/view/{id}", gzipViewHandler).Methods("GET").Name("view_torrent")
	Router.HandleFunc("/view/{id}", PostCommentHandler).Methods("POST").Name("post_comment")
	Router.Handle("/upload", gzipUploadHandler).Name("upload")
	Router.Handle("/user/register", gzipUserRegisterFormHandler).Name("user_register").Methods("GET")
	Router.Handle("/user/login", gzipUserLoginFormHandler).Name("user_login").Methods("GET")
	Router.Handle("/verify/email/{token}", gzipUserVerifyEmailHandler).Name("user_verify").Methods("GET")
	Router.Handle("/user/register", gzipUserRegisterPostHandler).Name("user_register").Methods("POST")
	Router.Handle("/user/login", gzipUserLoginPostHandler).Name("user_login").Methods("POST")
	Router.Handle("/user/logout", gzipUserLogoutHandler).Name("user_logout")
	Router.Handle("/user/{id}/{username}", gzipUserProfileHandler).Name("user_profile").Methods("GET")
	Router.Handle("/user/{id}/{username}", gzipUserProfileFormHandler).Name("user_profile").Methods("POST")
	Router.PathPrefix("/captcha").Methods("GET").HandlerFunc(captcha.ServeFiles)

	Router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
}
