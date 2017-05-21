package router

import (
	"net/http"
)

func FaqHandler(w http.ResponseWriter, r *http.Request) {
	ftv := FaqTemplateVariables{NewCommonVariables(r)}
	err := faqTemplate.ExecuteTemplate(w, "index.html", ftv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
