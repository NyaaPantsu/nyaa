package router

import (
	"github.com/gorilla/mux"

	"net/http"
)

var Router *mux.Router

func init() {
	Router = mux.NewRouter()

	cssHandler := http.FileServer(http.Dir("./public/css/"))
	jsHandler := http.FileServer(http.Dir("./public/js/"))
	imgHandler := http.FileServer(http.Dir("./public/img/"))
	http.Handle("/css/", http.StripPrefix("/public/css/", cssHandler))
	http.Handle("/js/", http.StripPrefix("/public/js/", jsHandler))
	http.Handle("/img/", http.StripPrefix("/public/img/", imgHandler))

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
	Router.HandleFunc("/user/register", UserRegisterFormHandler).Name("user_register")
}
