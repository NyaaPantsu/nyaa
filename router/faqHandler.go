package router

import (
	"html/template"
	"net/http"

	"github.com/ewhal/nyaa/templates"
	"github.com/gorilla/mux"
)

var faqTemplate = template.Must(template.New("FAQ").Funcs(FuncMap).ParseFiles("templates/index.html", "templates/FAQ.html"))

func FaqHandler(w http.ResponseWriter, r *http.Request) {
	searchForm := templates.NewSearchForm()
	searchForm.HideAdvancedSearch = true
	err := faqTemplate.ExecuteTemplate(w, "index.html", FaqTemplateVariables{templates.Navigation{}, searchForm, r.URL, mux.CurrentRoute(r)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
