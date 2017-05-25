package router

import (
	"net/http"
)

// NotFoundHandler : Controller for displaying 404 error page
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)

	nftv := newCommonVariables(r)

	err := notFoundTemplate.ExecuteTemplate(w, "index.html", nftv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
