package router

import (
	"html/template"
	"path/filepath"
)

var TemplateDir = "templates"

var homeTemplate, searchTemplate, faqTemplate, uploadTemplate, viewTemplate *template.Template

type templateLoader struct {
	templ **template.Template
	file  string
	name  string
}

// ReloadTemplates reloads templates on runtime
func ReloadTemplates() {
	templs := []templateLoader{
		templateLoader{
			templ: &homeTemplate,
			name:  "home",
			file:  "home.html",
		},
		templateLoader{
			templ: &searchTemplate,
			name:  "search",
			file:  "home.html",
		},
		templateLoader{
			templ: &uploadTemplate,
			name:  "upload",
			file:  "upload.html",
		},
		templateLoader{
			templ: &faqTemplate,
			name:  "FAQ",
			file:  "FAQ.html",
		},
		templateLoader{
			templ: &viewTemplate,
			name:  "view",
			file:  "view.html",
		},
	}
	for _, templ := range templs {
		t := template.Must(template.New(templ.name).Funcs(FuncMap).ParseFiles(filepath.Join(TemplateDir, "index.html"), filepath.Join(TemplateDir, templ.file)))
		t = template.Must(t.ParseGlob(filepath.Join("templates", "_*.html")))

		*templ.templ = t
	}
}
