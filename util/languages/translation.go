package languages

import (
    "github.com/nicksnyder/go-i18n/i18n"
    "html/template"
    )

    func SetTranslation(language string, tmpl *template.Template) {
    T, _ := i18n.Tfunc(language)
    tmpl.Funcs(map[string]interface{}{
        "T": func (str string)  template.HTML {
        	return template.HTML(T(str))
        	},
    })
    }