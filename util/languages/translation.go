package languages

import (
    "github.com/nicksnyder/go-i18n/i18n"
    "html/template"
    )

    func SetTranslation(language string, template *template.Template) {
    T, _ := i18n.Tfunc(language)
    template.Funcs(map[string]interface{}{
        "T": T,
    })
    }