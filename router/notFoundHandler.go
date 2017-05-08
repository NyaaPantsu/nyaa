package router

import (
	"html/template"
	"net/http"

	"github.com/ewhal/nyaa/util/languages"
	"github.com/gorilla/mux"
)

var notFoundTemplate = template.Must(template.New("NotFound").Funcs(FuncMap).ParseFiles("templates/index.html", "templates/404.html"))

func init() {
	// common
	template.Must(notFoundTemplate.ParseGlob("templates/_*.html"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)

	searchForm := NewSearchForm()
	searchForm.HideAdvancedSearch = true

	languages.SetTranslationFromRequest(notFoundTemplate, r, "en-us")
	err := notFoundTemplate.ExecuteTemplate(w, "index.html", NotFoundTemplateVariables{Navigation{}, searchForm, GetUser(r), r.URL, mux.CurrentRoute(r)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
