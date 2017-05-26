package router

import (
	"net/http"
)

// FaqHandler : Controller for FAQ view page
func FaqHandler(w http.ResponseWriter, r *http.Request) {
	ftv := newCommonVariables(r)
	err := faqTemplate.ExecuteTemplate(w, "index.html", ftv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
