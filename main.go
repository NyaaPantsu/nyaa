package main

import (
"bufio"
	"encoding/json"
	"fmt"
	"flag"
	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"

	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/service/torrent"
	"github.com/ewhal/nyaa/util/log"
	"github.com/ewhal/nyaa/config"

	"html"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var router *mux.Router
type SearchParam struct {
	Category   string
	Order      string
	Query      string
	Max        int
	Status     string
	Sort       string
}

func apiHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	page := vars["page"]
	pagenum, _ := strconv.Atoi(html.EscapeString(page))

	b := model.CategoryJson{Torrents: []model.TorrentsJson{}}
	maxPerPage := 50
	nbTorrents := 0

	torrents, nbTorrents := torrentService.GetAllTorrents(maxPerPage, maxPerPage*(pagenum-1))
	for i, _ := range torrents {
		res := torrents[i].ToJson()
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
	b := model.CategoryJson{Torrents: []model.TorrentsJson{}}

	torrent, err := torrentService.GetTorrentById(id)
	res := torrent.ToJson()
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
 	templates.ParseGlob("templates/_*.html") // common
	vars := mux.Vars(r)
	page := vars["page"]

	// db params url
	pagenum, _ := strconv.Atoi(html.EscapeString(page))
	if pagenum == 0 {
		pagenum = 1
	}
	
	b := []model.TorrentsJson{}
	
	search_param, torrents, nbTorrents := searchByQuery( r, pagenum )

	for i, _ := range torrents {
		res := torrents[i].ToJson()
		b = append(b, res)
	}

	navigationTorrents := Navigation{nbTorrents, search_param.Max, pagenum, "search_page"}
	searchForm := SearchForm{
		search_param.Query,
		search_param.Status,
		search_param.Category,
		search_param.Sort,
		search_param.Order,
	}
	htv := HomeTemplateVariables{b, torrentService.GetAllCategories(false), searchForm, navigationTorrents, r.URL, mux.CurrentRoute(r)}

	err := templates.ExecuteTemplate(w, "index.html", htv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func searchByQuery(r *http.Request, pagenum int) (SearchParam, []model.Torrents, int) {
	maxPerPage, errConv := strconv.Atoi(r.URL.Query().Get("max"))
	if errConv != nil {
		maxPerPage = 50 // default Value maxPerPage
	}
	
	search_param := SearchParam{}
	search_param.Max = maxPerPage
	search_param.Query = r.URL.Query().Get("q")
	search_param.Category = r.URL.Query().Get("c")
	search_param.Status = r.URL.Query().Get("s")
	search_param.Sort = r.URL.Query().Get("sort")
	search_param.Order = r.URL.Query().Get("order")

	catsSplit := strings.Split(search_param.Category, "_")
	// need this to prevent out of index panics
	var searchCatId, searchSubCatId string
	if len(catsSplit) == 2 {

		searchCatId = html.EscapeString(catsSplit[0])
		searchSubCatId = html.EscapeString(catsSplit[1])
	}
	if search_param.Sort == "" {
		search_param.Sort = "torrent_id"
	}
	if search_param.Order == "" {
		search_param.Order = "desc"
	}
	order_by := search_param.Sort + " " + search_param.Order

	parameters := torrentService.WhereParams{}
	conditions := []string{}
	if searchCatId != "" {
		conditions = append(conditions, "category_id = ?")
		parameters.Params = append(parameters.Params, searchCatId)
	}
	if searchSubCatId != "" {
		conditions = append(conditions, "sub_category_id = ?")
		parameters.Params = append(parameters.Params, searchSubCatId)
	}
	if search_param.Status != "" {
		if search_param.Status == "2" {
			conditions = append(conditions, "status_id != ?")
		} else {
			conditions = append(conditions, "status_id = ?")
		}
		parameters.params = append(parameters.params, search_param.Status)
	}
	searchQuerySplit := strings.Split(search_param.Query, " ")
	for i, _ := range searchQuerySplit {
		conditions = append(conditions, "torrent_name LIKE ?")
		parameters.Params = append(parameters.Params, "%"+searchQuerySplit[i]+"%")
	}

	parameters.Conditions = strings.Join(conditions[:], " AND ")
	log.Infof("SQL query is :: %s\n", parameters.Conditions)
	torrents, n := torrentService.GetTorrentsOrderBy(&parameters, order_by, maxPerPage, maxPerPage*(pagenum-1))
	return search_param, torrents, n
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	var templates = template.Must(template.New("FAQ").Funcs(funcMap).ParseFiles("templates/index.html", "templates/FAQ.html"))
 	templates.ParseGlob("templates/_*.html") // common
	err := templates.ExecuteTemplate(w, "index.html", FaqTemplateVariables{Navigation{}, NewSearchForm(), r.URL, mux.CurrentRoute(r)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func rssHandler(w http.ResponseWriter, r *http.Request) {

	_, torrents, _ := searchByQuery( r, 1 )
	created_as_time := time.Now()

	if len(torrents) > 0 {
		created_as_time = time.Unix(torrents[0].Date, 0)
	}
	feed := &feeds.Feed{
		Title:   "Nyaa Pantsu",
		Link:    &feeds.Link{Href: "https://nyaa.pantsu.cat/"},
		Created: created_as_time,
	}
	feed.Items = []*feeds.Item{}
	feed.Items = make([]*feeds.Item, len(torrents))

	for i, _ := range torrents {
		timestamp_as_time := time.Unix(torrents[0].Date, 0)
		torrent_json := torrents[i].ToJson()
		feed.Items[i] = &feeds.Item{
			// need a torrent view first
			//Id:		URL + torrents[i].Hash,
			Title:       torrents[i].Name,
			Link:        &feeds.Link{Href: string(torrent_json.Magnet)},
			Description: "",
			Created:     timestamp_as_time,
			Updated:     timestamp_as_time,
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
 	templates.ParseGlob("templates/_*.html") // common
	vars := mux.Vars(r)
	id := vars["id"]

	torrent, err := torrentService.GetTorrentById(id)
	b := torrent.ToJson()

	htv := ViewTemplateVariables{b, NewSearchForm(), Navigation{}, r.URL, mux.CurrentRoute(r)}

	err = templates.ExecuteTemplate(w, "index.html", htv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
var templates = template.Must(template.New("home").Funcs(funcMap).ParseFiles("templates/index.html", "templates/home.html"))
	templates.ParseGlob("templates/_*.html") // common
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

	b := []model.TorrentsJson{}
	torrents, nbTorrents := torrentService.GetAllTorrents(maxPerPage, maxPerPage*(pagenum-1))

	for i, _ := range torrents {
		res := torrents[i].ToJson()
		b = append(b, res)
	}

	navigationTorrents := Navigation{nbTorrents, maxPerPage, pagenum, "search_page"}
	htv := HomeTemplateVariables{b, torrentService.GetAllCategories(false), NewSearchForm(), navigationTorrents, r.URL, mux.CurrentRoute(r)}

	err := templates.ExecuteTemplate(w, "index.html", htv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func RunServer(conf *config.Config) {
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
		Addr:         fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	err := srv.ListenAndServe()
	log.CheckError(err)
}

func main() {

	conf := config.NewConfig()
	conf_bind := conf.BindFlags()
	defaults := flag.Bool("print-defaults", false, "print the default configuration file on stdout")
	flag.Parse()
	if *defaults {
		stdout := bufio.NewWriter(os.Stdout)
		conf.Pretty(stdout)
		stdout.Flush()
		os.Exit(0)
	} else {
		conf_bind()
		RunServer(conf)
	}
}
