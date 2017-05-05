package main

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"html"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var db *gorm.DB
var router *mux.Router
var debugLogger *log.Logger
var trackers = "&tr=udp://zer0day.to:1337/announce&tr=udp://tracker.leechers-paradise.org:6969&tr=udp://explodie.org:6969&tr=udp://tracker.opentrackr.org:1337&tr=udp://tracker.coppersurfer.tk:6969&tr=http://tracker.baka-sub.cf/announce"

func getDBHandle() *gorm.DB {
	dbInit, err := gorm.Open("sqlite3", "./nyaa.db")

	// Migrate the schema of Torrents
	dbInit.AutoMigrate(&Torrents{}, &Categories{}, &Sub_Categories{}, &Statuses{})

	checkErr(err)
	return dbInit
}

func checkErr(err error) {
	if err != nil {
		debugLogger.Println("   " + err.Error())
	}
}

func unZlib(description []byte) string {
	if len(description) > 0 {
		b := bytes.NewReader(description)
		log.Println(b)
		z, err := zlib.NewReader(b)
		checkErr(err)
		defer z.Close()
		p, err := ioutil.ReadAll(z)
		checkErr(err)
		return string(p)
	}
	return ""
}

func apiHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	page := vars["page"]
	pagenum, _ := strconv.Atoi(html.EscapeString(page))

	b := CategoryJson{Torrents: []TorrentsJson{}}
	maxPerPage := 50
	nbTorrents := 0

	torrents, nbTorrents := getAllTorrents(maxPerPage, maxPerPage*(pagenum-1))
	for i, _ := range torrents {
		res := torrents[i].toJson()
		b.Torrents = append(b.Torrents, res)
	}

	b.QueryRecordCount = maxPerPage
	b.TotalRecordCount = nbTorrents
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(b)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func apiViewHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	b := CategoryJson{Torrents: []TorrentsJson{}}

	torrent, err := getTorrentById(id)
	res := torrent.toJson()
	b.Torrents = append(b.Torrents, res)

	b.QueryRecordCount = 1
	b.TotalRecordCount = 1
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(b)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	var templates = template.Must(template.New("home").Funcs(funcMap).ParseFiles("templates/index.html", "templates/home.html"))
	vars := mux.Vars(r)
	page := vars["page"]

	// db params url
	maxPerPage, errConv := strconv.Atoi(r.URL.Query().Get("max"))
	if errConv != nil {
		maxPerPage = 50 // default Value maxPerPage
	}
	pagenum, _ := strconv.Atoi(html.EscapeString(page))
	if pagenum == 0 {
		pagenum = 1
	}
	searchQuery := r.URL.Query().Get("q")
	cat := r.URL.Query().Get("c")
	stat := r.URL.Query().Get("s")
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	catsSplit := strings.Split(cat, "_")
	// need this to prevent out of index panics
	var searchCatId, searchSubCatId string
	if len(catsSplit) == 2 {

		searchCatId = html.EscapeString(catsSplit[0])
		searchSubCatId = html.EscapeString(catsSplit[1])
	}
	if sort == "" {
		sort = "torrent_id"
	}
	if order == "" {
		order = "desc"
	}
	order_by := sort + " " + order

	nbTorrents := 0

	b := []TorrentsJson{}

	parameters := createWhereParams("torrent_name LIKE ? AND status_id LIKE ? AND category_id LIKE ? AND sub_category_id LIKE ?",
		"%"+searchQuery+"%", stat+"%", searchCatId+"%", searchSubCatId+"%")
	torrents, nbTorrents := getTorrentsOrderBy(&parameters, order_by, maxPerPage, maxPerPage*(pagenum-1))

	for i, _ := range torrents {
		res := torrents[i].toJson()
		b = append(b, res)
	}

	navigationTorrents := Navigation{nbTorrents, maxPerPage, pagenum, "search_page"}
	searchForm := SearchForm{searchQuery, stat, cat, sort, order}
	htv := HomeTemplateVariables{b, getAllCategories(false), searchForm, navigationTorrents, r.URL, mux.CurrentRoute(r)}

	err := templates.ExecuteTemplate(w, "index.html", htv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func safe(s string) template.URL {
	return template.URL(s)
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	var templates = template.Must(template.New("FAQ").Funcs(funcMap).ParseFiles("templates/index.html", "templates/FAQ.html"))
	err := templates.ExecuteTemplate(w, "index.html", FaqTemplateVariables{Navigation{}, NewSearchForm(), r.URL, mux.CurrentRoute(r)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func rssHandler(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//category := vars["c"]

	// db params url
	//maxPerPage := 50 // default Value maxPerPage

	torrents := getFeeds()
	created := time.Now().String()
	if len(torrents) > 0 {
		created = torrents[0].Timestamp
	}
	created_as_time, err := time.Parse("2006-01-02 15:04:05", created)
	if err == nil {

	}
	feed := &feeds.Feed{
		Title:   "Nyaa Pantsu",
		Link:    &feeds.Link{Href: "https://nyaa.pantsu.cat/"},
		Created: created_as_time,
	}
	feed.Items = []*feeds.Item{}
	feed.Items = make([]*feeds.Item, len(torrents))

	for i, _ := range torrents {
		timestamp_as_time, err := time.Parse("2006-01-02 15:04:05", torrents[i].Timestamp)
		if err == nil {
			feed.Items[i] = &feeds.Item{
				// need a torrent view first
				//Id:		URL + torrents[i].Hash,
				Title:       torrents[i].Name,
				Link:        &feeds.Link{Href: string(torrents[i].Magnet)},
				Description: "",
				Created:     timestamp_as_time,
				Updated:     timestamp_as_time,
			}
		}
	}

	rss, err := feed.ToRss()
	if err == nil {
		w.Write([]byte(rss))
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func viewHandler(w http.ResponseWriter, r *http.Request) {
	var templates = template.Must(template.ParseFiles("templates/index.html", "templates/view.html"))
	vars := mux.Vars(r)
	id := vars["id"]

	torrent, err := getTorrentById(id)
	b := torrent.toJson()

	htv := ViewTemplateVariables{b, NewSearchForm(), Navigation{}, r.URL, mux.CurrentRoute(r)}

	err = templates.ExecuteTemplate(w, "index.html", htv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	var templates = template.Must(template.New("home").Funcs(funcMap).ParseFiles("templates/index.html", "templates/home.html"))
	vars := mux.Vars(r)
	page := vars["page"]

	// db params url
	maxPerPage, errConv := strconv.Atoi(r.URL.Query().Get("max"))
	if errConv != nil {
		maxPerPage = 50 // default Value maxPerPage
	}

	nbTorrents := 0
	pagenum, _ := strconv.Atoi(html.EscapeString(page))
	if pagenum == 0 {
		pagenum = 1
	}

	b := []TorrentsJson{}
	torrents, nbTorrents := getAllTorrents(maxPerPage, maxPerPage*(pagenum-1))

	for i, _ := range torrents {
		res := torrents[i].toJson()
		b = append(b, res)
	}

	navigationTorrents := Navigation{nbTorrents, maxPerPage, pagenum, "search_page"}
	htv := HomeTemplateVariables{b, getAllCategories(false), NewSearchForm(), navigationTorrents, r.URL, mux.CurrentRoute(r)}

	err := templates.ExecuteTemplate(w, "index.html", htv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func main() {

	db = getDBHandle()
	router = mux.NewRouter()

	cssHandler := http.FileServer(http.Dir("./css/"))
	jsHandler := http.FileServer(http.Dir("./js/"))
	imgHandler := http.FileServer(http.Dir("./img/"))
	http.Handle("/css/", http.StripPrefix("/css/", cssHandler))
	http.Handle("/js/", http.StripPrefix("/js/", jsHandler))
	http.Handle("/img/", http.StripPrefix("/img/", imgHandler))

	// Routes,
	router.HandleFunc("/", rootHandler).Name("home")
	router.HandleFunc("/page/{page:[0-9]+}", rootHandler).Name("home_page")
	router.HandleFunc("/search", searchHandler).Name("search")
	router.HandleFunc("/search/{page}", searchHandler).Name("search_page")
	router.HandleFunc("/api/{page}", apiHandler).Methods("GET")
	router.HandleFunc("/api/view/{id}", apiViewHandler).Methods("GET")
	router.HandleFunc("/faq", faqHandler).Name("faq")
	router.HandleFunc("/feed.xml", rssHandler)
	router.HandleFunc("/view/{id}", viewHandler).Name("view_torrent")

	http.Handle("/", router)

	// Set up server,
	srv := &http.Server{
		Addr:         "localhost:9999",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	err := srv.ListenAndServe()
	checkErr(err)
}
