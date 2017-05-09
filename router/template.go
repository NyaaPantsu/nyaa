package router

import (
	"html/template"
	"path/filepath"
)

var TemplateDir = "templates"

var homeTemplate, searchTemplate, faqTemplate, uploadTemplate, viewTemplate, viewRegisterTemplate, viewLoginTemplate, viewRegisterSuccessTemplate, viewVerifySuccessTemplate, viewProfileTemplate, viewProfileEditTemplate, viewUserDeleteTemplate, notFoundTemplate *template.Template

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
		templateLoader{
			templ: &viewRegisterTemplate,
			name:  "user_register",
			file:  "user/register.html",
		},
		templateLoader{
			templ: &viewRegisterSuccessTemplate,
			name:  "user_register_success",
			file:  "user/signup_success.html",
		},
		templateLoader{
			templ: &viewVerifySuccessTemplate,
			name:  "user_verify_success",
			file:  "user/verify_success.html",
		},
		templateLoader{
			templ: &viewLoginTemplate,
			name:  "user_login",
			file:  "user/login.html",
		},
		templateLoader{
			templ: &viewProfileTemplate,
			name:  "user_profile",
			file:  "user/profile.html",
		},
		templateLoader{
			templ: &viewProfileEditTemplate,
			name:  "user_profile",
			file:  "user/profile_edit.html",
		},
		templateLoader{
			templ: &viewUserDeleteTemplate,
			name:  "user_delete",
			file:  "user/delete_success.html",
		},
		templateLoader{
			templ: &notFoundTemplate,
			name:  "404",
			file:  "404.html",
		},
	}
	for _, templ := range templs {
		t := template.Must(template.New(templ.name).Funcs(FuncMap).ParseFiles(filepath.Join(TemplateDir, "index.html"), filepath.Join(TemplateDir, templ.file)))
		t = template.Must(t.ParseGlob(filepath.Join("templates", "_*.html")))

		*templ.templ = t
	}
}
