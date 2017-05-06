package router

import (
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	var templates = template.Must(template.New("upload").Funcs(FuncMap).ParseFiles("templates/index.html", "templates/upload.html"))
	templates.ParseGlob("templates/_*.html") // common

	var uploadForm UploadForm
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		uploadForm = UploadForm{
			r.Form.Get("name"),
			r.Form.Get("magnet"),
			r.Form.Get("c"),
			r.Form.Get("desc"),
		}
		//validate name + hash
		//add to db and redirect depending on result
	}

	htv := UploadTemplateVariables{uploadForm, NewSearchForm(), Navigation{}, r.URL, mux.CurrentRoute(r)}

	err := templates.ExecuteTemplate(w, "index.html", htv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
