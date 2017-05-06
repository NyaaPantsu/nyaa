package router

import (
	"github.com/ewhal/nyaa/service/torrent"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	var templates = template.Must(template.New("view").Funcs(FuncMap).ParseFiles("templates/index.html", "templates/view.html"))
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
