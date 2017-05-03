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
	"os"
	"strconv"
	"strings"
	"time"
)

var dbHandle *sql.DB
var templates = template.Must(template.ParseFiles("index.html"))
var debugLogger *log.Logger

type Record struct {
	Category         string    `json: "category"`
	Records          []Records `json: "records"`
	QueryRecordCount int       `json: "queryRecordCount"`
	TotalRecordCount int       `json: "totalRecordCount"`
}

type Records struct {
	Id     string `json: "id"`
	Name   string `json: "name"`
	Hash   string `json: "hash"`
	Magnet string `json: "magnet"`
}

func getDBHandle() *sql.DB {
	db, err := sql.Open("sqlite3", "./nyaa.db")
	checkErr(err)
	return db
}

func checkErr(err error) {
	if err != nil {
		debugLogger.Println("   " + err.Error())
		os.Exit(1)
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	page := vars["page"]
	pagenum, _ := strconv.Atoi(html.EscapeString(page))
	b := Record{Records: []Records{}}
	rows, err := dbHandle.Query("select torrent_id, torrent_name, torrent_hash from torrents ORDER BY torrent_id DESC LIMIT 50 offset ?", 50*pagenum-1)
	for rows.Next() {
		var id, name, hash, magnet string
		rows.Scan(&id, &name, &hash)
		magnet = "magnet:?xt=urn:btih:" + hash + "&dn=" + url.QueryEscape(name) + "&tr=udp://tracker.openbittorrent.com"
		res := Records{
			Id:     id,
			Name:   name,
			Hash:   hash,
			Magnet: magnet}

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
	b := Record{Category: "_", Records: []Records{}}
	rows, err := dbHandle.Query("select torrent_id, torrent_name, torrent_hash from torrents where torrent_id = ? ORDER BY torrent_id DESC", html.EscapeString(id))
	for rows.Next() {
		var id, name, hash, magnet string
		rows.Scan(&id, &name, &hash)
		magnet = "magnet:?xt=urn:btih:" + hash + "&dn=" + url.QueryEscape(name) + "&tr=udp://tracker.openbittorrent.com"
		res := Records{
			Id:     id,
			Name:   name,
			Hash:   hash,
			Magnet: magnet}

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
	pagenum, _ := strconv.Atoi(html.EscapeString(page))
	param1 := r.URL.Query().Get("q")
	cat := r.URL.Query().Get("c")
	param2 := strings.Split(cat, "_")[0]
	param3 := strings.Split(cat, "_")[1]
	b := Record{Category: cat, Records: []Records{}}
	rows, err := dbHandle.Query("select torrent_id, torrent_name, torrent_hash from torrents "+
		"where torrent_name LIKE ? AND category_id LIKE ? AND sub_category_id LIKE ? "+
		"ORDER BY torrent_id DESC LIMIT 50 offset ?",
		"%"+html.EscapeString(param1)+"%", html.EscapeString(param2)+"%", html.EscapeString(param3)+"%", 50*pagenum-1)
	for rows.Next() {
		var id, name, hash, magnet string
		rows.Scan(&id, &name, &hash)
		magnet = "magnet:?xt=urn:btih:" + hash + "&dn=" + url.QueryEscape(name) + "&tr=udp://tracker.openbittorrent.com"
		res := Records{
			Id:     id,
			Name:   name,
			Hash:   hash,
			Magnet: magnet}

		b.Records = append(b.Records, res)

	}
	rows.Close()

	err = templates.ExecuteTemplate(w, "index.html", &b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]
	pagenum, _ := strconv.Atoi(html.EscapeString(page))
	b := Record{Category: "_", Records: []Records{}}
	rows, err := dbHandle.Query("select torrent_id, torrent_name, torrent_hash from torrents ORDER BY torrent_id DESC LIMIT 50 offset ?", 50*pagenum-1)
	for rows.Next() {
		var id, name, hash, magnet string
		rows.Scan(&id, &name, &hash)
		magnet = "magnet:?xt=urn:btih:" + hash + "&dn=" + url.QueryEscape(name) + "&tr=udp://tracker.openbittorrent.com"
		res := Records{
			Id:     id,
			Name:   name,
			Hash:   hash,
			Magnet: magnet}

		b.Records = append(b.Records, res)

	}
	b.QueryRecordCount = 50
	b.TotalRecordCount = 1473098
	rows.Close()
	err = templates.ExecuteTemplate(w, "index.html", &b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func main() {

	dbHandle = getDBHandle()
	router := mux.NewRouter()

	// Routes,
	router.HandleFunc("/", rootHandler)
	router.HandleFunc("/page/{page}", rootHandler)
	router.HandleFunc("/search", searchHandler)
	router.HandleFunc("/search/{page}", searchHandler)
	router.HandleFunc("/api/{page}", apiHandler).Methods("GET")
	router.HandleFunc("/api/torrent/{id}", singleapiHandler).Methods("GET")
	// Set up server,
	srv := &http.Server{
		Handler:      router,
		Addr:         "localhost:9999",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	err := srv.ListenAndServe()
	checkErr(err)
}
