package router

import(
	"github.com/ewhal/nyaa/util/search"
	"net/http"
	"html/template"
	"github.com/gorilla/mux"
	"html"
	"strconv"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/service/torrent"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	var templates = template.Must(template.New("home").Funcs(FuncMap).ParseFiles("templates/index.html", "templates/home.html"))
 	templates.ParseGlob("templates/_*.html") // common
	vars := mux.Vars(r)
	page := vars["page"]

	// db params url
	pagenum, _ := strconv.Atoi(html.EscapeString(page))
	if pagenum == 0 {
		pagenum = 1
	}
	
	b := []model.TorrentsJson{}
	
	search_param, torrents, nbTorrents := search.SearchByQuery( r, pagenum )

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