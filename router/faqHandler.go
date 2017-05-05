package router


import(
	"net/http"
	"html/template"
	"github.com/gorilla/mux"
)

func FaqHandler(w http.ResponseWriter, r *http.Request) {
	var templates = template.Must(template.New("FAQ").Funcs(FuncMap).ParseFiles("templates/index.html", "templates/FAQ.html"))
 	templates.ParseGlob("templates/_*.html") // common
	err := templates.ExecuteTemplate(w, "index.html", FaqTemplateVariables{Navigation{}, NewSearchForm(), r.URL, mux.CurrentRoute(r)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}