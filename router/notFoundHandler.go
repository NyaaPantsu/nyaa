package router

import (
	"net/http"

	"github.com/ewhal/nyaa/util/languages"
	"github.com/gorilla/mux"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)

	languages.SetTranslationFromRequest(notFoundTemplate, r, "en-us")
	err := notFoundTemplate.ExecuteTemplate(w, "index.html", NotFoundTemplateVariables{Navigation{}, NewSearchForm(), GetUser(r), r.URL, mux.CurrentRoute(r)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
