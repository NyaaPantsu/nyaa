package router

import (
	"net/http"

	"github.com/NyaaPantsu/nyaa/util/languages"
	"github.com/gorilla/mux"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)

	languages.SetTranslationFromRequest(notFoundTemplate, r)
	err := notFoundTemplate.ExecuteTemplate(w, "index.html", NotFoundTemplateVariables{NewNavigation(), NewSearchForm(), GetUser(r), r.URL, mux.CurrentRoute(r)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
