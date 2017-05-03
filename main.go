package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"html"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var dbHandle *sql.DB
var templates = template.Must(template.ParseFiles("index.html", "FAQ.html"))
var debugLogger *log.Logger
var trackers = "&tr=udp://zer0day.to:1337/announce&tr=udp://tracker.leechers-paradise.org:6969&tr=udp://explodie.org:6969&tr=udp://tracker.opentrackr.org:1337&tr=udp://tracker.coppersurfer.tk:6969"

type Record struct {
	Category         string    `json: "category"`
	Records          []Records `json: "records"`
	QueryRecordCount int       `json: "queryRecordCount"`
	TotalRecordCount int       `json: "totalRecordCount"`
}

type Records struct {
	Id     string       `json: "id"`
	Name   string       `json: "name"`
	Status int          `json: "status"`
	Hash   string       `json: "hash"`
	Magnet template.URL `json: "magnet"`
}

func getDBHandle() *sql.DB {
	db, err := sql.Open("sqlite3", "./nyaa.db")
	checkErr(err)
	return db
}

func checkErr(err error) {
	if err != nil {
		debugLogger.Println("   " + err.Error())
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	page := vars["page"]
	pagenum, _ := strconv.Atoi(html.EscapeString(page))
	b := Record{Records: []Records{}}
	rows, err := dbHandle.Query("select torrent_id, torrent_name, status_id, torrent_hash from torrents ORDER BY torrent_id DESC LIMIT 50 offset ?", 50*(pagenum-1))
	checkErr(err)
	for rows.Next() {
		var id, name, hash, magnet string
		var status int
		rows.Scan(&id, &name, &status, &hash)
		magnet = "magnet:?xt=urn:btih:" + hash + "&dn=" + url.QueryEscape(name) + trackers
		res := Records{
			Id:     id,
			Name:   name,
			Status: status,
			Hash:   hash,
			Magnet: safe(magnet)}

		b.Records = append(b.Records, res)

	}
	b.QueryRecordCount = 50
	b.TotalRecordCount = 1473098
	rows.Close()
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(b)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func singleapiHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	b := Record{Records: []Records{}}
	rows, err := dbHandle.Query("select torrent_id, torrent_name, status_id, torrent_hash from torrents where torrent_id = ? ORDER BY torrent_id DESC", html.EscapeString(id))
	checkErr(err)
	for rows.Next() {
		var id, name, hash, magnet string
		var status int
		rows.Scan(&id, &name, &status, &hash)
		magnet = "magnet:?xt=urn:btih:" + hash + "&dn=" + url.QueryEscape(name) + trackers
		res := Records{
			Id:     id,
			Name:   name,
			Status: status,
			Hash:   hash,
			Magnet: safe(magnet)}

		b.Records = append(b.Records, res)

	}
	b.QueryRecordCount = 1
	b.TotalRecordCount = 1473098
	rows.Close()
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(b)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]

	// db params url
	maxPerPage, errConv := strconv.Atoi(r.URL.Query().Get("max"))
	if errConv != nil {
		maxPerPage = 50 // default Value maxPerPage
	}
	pagenum, _ := strconv.Atoi(html.EscapeString(page))
	param1 := r.URL.Query().Get("q")
	cat := r.URL.Query().Get("c")

	var param2, param3 string
	params := strings.Split(cat, "_")
	if len(params) == 2 {
		param2 = params[0]
		param3 = params[1]
	}

	nbTorrents := 0

	b := Record{Category: cat, Records: []Records{}}
	rows, err := dbHandle.Query("select torrent_id, torrent_name, status_id, torrent_hash from torrents "+
		"where torrent_name LIKE ? AND category_id LIKE ? AND sub_category_id LIKE ? "+
		"ORDER BY torrent_id DESC LIMIT ? offset ?",
		"%"+html.EscapeString(param1)+"%", html.EscapeString(param2)+"%", html.EscapeString(param3)+"%", maxPerPage, maxPerPage*(pagenum-1))

	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		nbTorrents++
		var id, name, hash, magnet string
		var status int
		rows.Scan(&id, &name, &status, &hash)
		magnet = "magnet:?xt=urn:btih:" + hash + "&dn=" + url.QueryEscape(name) + trackers
		res := Records{
			Id:     id,
			Name:   name,
			Status: status,
			Hash:   hash,
			Magnet: safe(magnet)}

		b.Records = append(b.Records, res)

	}
	b.QueryRecordCount = maxPerPage
	b.TotalRecordCount = nbTorrents

	err = templates.ExecuteTemplate(w, "index.html", &b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func safe(s string) template.URL {
	return template.URL(s)
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	err = templates.ExecuteTemplate(w, "FAQ.html", &b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]

	// db params url
	maxPerPage, errConv := strconv.Atoi(r.URL.Query().Get("max"))
	if errConv != nil {
		maxPerPage = 50 // default Value maxPerPage
	}

	nbTorrents := 0
	pagenum, _ := strconv.Atoi(html.EscapeString(page))
	b := Record{Category: "_", Records: []Records{}}
	rows, err := dbHandle.Query("select torrent_id, torrent_name, status_id, torrent_hash from torrents ORDER BY torrent_id DESC LIMIT ? offset ?", maxPerPage, maxPerPage*(pagenum-1))
	checkErr(err)
	for rows.Next() {
		nbTorrents++
		var id, name, hash, magnet string
		var status int
		rows.Scan(&id, &name, &status, &hash)
		magnet = "magnet:?xt=urn:btih:" + hash + "&dn=" + url.QueryEscape(name) + trackers
		res := Records{
			Id:     id,
			Name:   name,
			Status: status,
			Hash:   hash,
			Magnet: safe(magnet)}

		b.Records = append(b.Records, res)

	}
	b.QueryRecordCount = maxPerPage
	b.TotalRecordCount = nbTorrents
	rows.Close()
	err = templates.ExecuteTemplate(w, "index.html", &b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func main() {

	dbHandle = getDBHandle()
	router := mux.NewRouter()

	cssHandler := http.FileServer(http.Dir("./css/"))
	jsHandler := http.FileServer(http.Dir("./js/"))
	http.Handle("/css/", http.StripPrefix("/css/", cssHandler))
	http.Handle("/js/", http.StripPrefix("/js/", jsHandler))

	// Routes,
	router.HandleFunc("/", rootHandler)
	router.HandleFunc("/page/{page}", rootHandler)
	router.HandleFunc("/search", searchHandler)
	router.HandleFunc("/search/{page}", searchHandler)
	router.HandleFunc("/api/{page}", apiHandler).Methods("GET")
	router.HandleFunc("/api/torrent/{id}", singleapiHandler).Methods("GET")
	router.HandleFunc("/faq", faqHandler)

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
