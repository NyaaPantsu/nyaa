package router

import (
	"html"
	"html/template"
	"net/http"
	"strconv"

	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/service/torrent"
	"github.com/ewhal/nyaa/templates"
	"github.com/gorilla/mux"
)

var homeTemplate = template.Must(template.New("home").Funcs(FuncMap).ParseFiles("templates/index.html", "templates/home.html"))

func HomeHandler(w http.ResponseWriter, r *http.Request) {
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

	navigationTorrents := templates.Navigation{nbTorrents, maxPerPage, pagenum, "search_page"}
	htv := HomeTemplateVariables{b, torrentService.GetAllCategories(false), templates.NewSearchForm(), navigationTorrents, r.URL, mux.CurrentRoute(r)}

	err := homeTemplate.ExecuteTemplate(w, "index.html", htv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
