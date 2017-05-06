package router

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

var uploadTemplate = template.Must(template.New("upload").Funcs(FuncMap).ParseFiles("templates/index.html", "templates/upload.html"))

func init() {
	template.Must(uploadTemplate.ParseGlob("templates/_*.html")) // common
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var uploadForm UploadForm
	if r.Method == "POST" {
		defer r.Body.Close()
		err = uploadForm.ExtractInfo(r)
		if err == nil {
			//validate name + hash
			//add to db and redirect depending on result
		}
	} else if r.Method == "GET" {
		htv := UploadTemplateVariables{uploadForm, NewSearchForm(), Navigation{}, r.URL, mux.CurrentRoute(r)}
		err = uploadTemplate.ExecuteTemplate(w, "index.html", htv)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
