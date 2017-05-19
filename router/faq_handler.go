package router

import (
	"net/http"

	"github.com/NyaaPantsu/nyaa/util/languages"
	"github.com/gorilla/mux"
)

func FaqHandler(w http.ResponseWriter, r *http.Request) {
	languages.SetTranslationFromRequest(faqTemplate, r)
	err := faqTemplate.ExecuteTemplate(w, "index.html", FaqTemplateVariables{NewNavigation(), NewSearchForm(), GetUser(r), r.URL, mux.CurrentRoute(r)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
