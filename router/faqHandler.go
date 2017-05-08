package router

import (
	"net/http"

	"github.com/ewhal/nyaa/util/languages"
	"github.com/gorilla/mux"
)

func FaqHandler(w http.ResponseWriter, r *http.Request) {
	searchForm := NewSearchForm()
	searchForm.HideAdvancedSearch = true

	languages.SetTranslationFromRequest(faqTemplate, r, "en-us")
	err := faqTemplate.ExecuteTemplate(w, "index.html", FaqTemplateVariables{Navigation{}, searchForm, GetUser(r), r.URL, mux.CurrentRoute(r)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
