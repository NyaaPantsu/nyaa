package router

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

var faqTemplate = template.Must(template.New("FAQ").Funcs(FuncMap).ParseFiles("templates/index.html", "templates/FAQ.html"))

func init() {
	// common
	template.Must(faqTemplate.ParseGlob("templates/_*.html"))
}

func FaqHandler(w http.ResponseWriter, r *http.Request) {
	searchForm := NewSearchForm()
	searchForm.HideAdvancedSearch = true
	err := faqTemplate.ExecuteTemplate(w, "index.html", FaqTemplateVariables{Navigation{}, searchForm, r.URL, mux.CurrentRoute(r)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
