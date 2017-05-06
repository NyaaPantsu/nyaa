package router

import (
	"html/template"
	"net/http"

	"github.com/ewhal/nyaa/service/torrent"
	"github.com/ewhal/nyaa/templates"
	"github.com/gorilla/mux"
)

var viewTemplate = template.Must(template.New("view").Funcs(FuncMap).ParseFiles("templates/index.html", "templates/view.html"))

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	torrent, err := torrentService.GetTorrentById(id)
	b := torrent.ToJson()

	htv := ViewTemplateVariables{b, templates.NewSearchForm(), templates.Navigation{}, r.URL, mux.CurrentRoute(r)}

	err = viewTemplate.ExecuteTemplate(w, "index.html", htv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
