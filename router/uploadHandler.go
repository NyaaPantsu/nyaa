package router

import (
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	var templates = template.Must(template.New("upload").Funcs(FuncMap).ParseFiles("templates/index.html", "templates/upload.html"))
	templates.ParseGlob("templates/_*.html") // common
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
		err = templates.ExecuteTemplate(w, "index.html", htv)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
