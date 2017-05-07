package router

import (
	"net/http"

	"github.com/ewhal/nyaa/service/captcha"
	"github.com/gorilla/mux"
)

var Router *mux.Router

func init() {
	Router = mux.NewRouter()

	cssHandler := http.FileServer(http.Dir("./public/css/"))
	jsHandler := http.FileServer(http.Dir("./public/js/"))
	imgHandler := http.FileServer(http.Dir("./public/img/"))
	http.Handle("/css/", http.StripPrefix("/css/", cssHandler))
	http.Handle("/js/", http.StripPrefix("/js/", jsHandler))
	http.Handle("/img/", http.StripPrefix("/img/", imgHandler))

	// Routes,
	Router.HandleFunc("/", HomeHandler).Name("home")
	Router.HandleFunc("/page/{page:[0-9]+}", HomeHandler).Name("home_page")
	Router.HandleFunc("/search", SearchHandler).Name("search")
	Router.HandleFunc("/search/{page}", SearchHandler).Name("search_page")
	Router.HandleFunc("/api/{page}", ApiHandler).Methods("GET")
	Router.HandleFunc("/api/view/{id}", ApiViewHandler).Methods("GET")
	Router.HandleFunc("/faq", FaqHandler).Name("faq")
	Router.HandleFunc("/feed", RssHandler).Name("feed")
	Router.HandleFunc("/view/{id}", ViewHandler).Name("view_torrent")
	Router.HandleFunc("/upload", UploadHandler).Name("upload")
	Router.HandleFunc("/user/register", UserRegisterFormHandler).Name("user_register").Methods("GET")
	Router.HandleFunc("/user/login", UserLoginFormHandler).Name("user_login").Methods("GET")
	Router.HandleFunc("/verify/email/{token}", UserVerifyEmailHandler).Name("user_verify").Methods("GET")
	Router.HandleFunc("/user/register", UserRegisterPostHandler).Name("user_register_post").Methods("POST")
	Router.HandleFunc("/user/login", UserLoginPostHandler).Name("user_login_post").Methods("POST")
	Router.HandleFunc("/user/{id}", UserProfileHandler).Name("user_profile")
	Router.PathPrefix("/captcha").Methods("GET").HandlerFunc(captcha.ServeFiles)

	Router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
}
