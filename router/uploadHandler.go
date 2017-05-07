package router

import (
	"html/template"
	"net/http"

	"github.com/ewhal/nyaa/service/captcha"
	"github.com/gorilla/mux"
)

var uploadTemplate = template.Must(template.New("upload").Funcs(FuncMap).ParseFiles("templates/index.html", "templates/upload.html"))

func init() {
	template.Must(uploadTemplate.ParseGlob("templates/_*.html")) // common
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case "POST":
		var form UploadForm
		defer r.Body.Close()
		err = form.ExtractInfo(r)
		if err == nil {
			// validate name + hash
			// authenticate captcha
			// add to db and redirect depending on result
		}
	case "GET":
		htv := UploadTemplateVariables{
			Upload: UploadForm{
				CaptchaID: captcha.GetID(r.RemoteAddr),
			},
			Search: NewSearchForm(),
			URL:    r.URL,
			Route:  mux.CurrentRoute(r),
		}
		err = uploadTemplate.ExecuteTemplate(w, "index.html", htv)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
