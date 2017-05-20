package router

import (
	"net/http"

	"github.com/NyaaPantsu/nyaa/util/languages"
	"github.com/gorilla/mux"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)

	nftv := NotFoundTemplateVariables{
		Navigation: NewNavigation(),
		Search:     NewSearchForm(),
		T:          languages.GetTfuncFromRequest(r),
		User:       GetUser(r),
		URL:        r.URL,
		Route:      mux.CurrentRoute(r),
	}

	err := notFoundTemplate.ExecuteTemplate(w, "index.html", nftv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
