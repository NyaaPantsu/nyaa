package router

import (
	"net/http"

	"github.com/NyaaPantsu/nyaa/util/languages"
	"github.com/gorilla/mux"
)

func FaqHandler(w http.ResponseWriter, r *http.Request) {
	languages.SetTranslationFromRequest(faqTemplate, r)
	ftv := FaqTemplateVariables{
		Navigation: NewNavigation(),
		Search:     NewSearchForm(),
		User:       GetUser(r),
		URL:        r.URL,
		Route:      mux.CurrentRoute(r),
	}
	err := faqTemplate.ExecuteTemplate(w, "index.html", ftv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
