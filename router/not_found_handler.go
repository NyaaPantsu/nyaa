package router

import (
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)

	nftv := NotFoundTemplateVariables{NewCommonVariables(r)}

	err := notFoundTemplate.ExecuteTemplate(w, "index.html", nftv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
