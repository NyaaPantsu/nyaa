package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

func FaqHandler(w http.ResponseWriter, r *http.Request) {
	searchForm := NewSearchForm()
	searchForm.HideAdvancedSearch = true
	err := faqTemplate.ExecuteTemplate(w, "index.html", FaqTemplateVariables{Navigation{}, searchForm, r.URL, mux.CurrentRoute(r)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
