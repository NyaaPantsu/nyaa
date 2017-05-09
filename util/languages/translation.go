package languages

import (
	"fmt"
	"github.com/nicksnyder/go-i18n/i18n"
	"html/template"
	"net/http"
)

func SetTranslation(tmpl *template.Template, language string, languages ...string) i18n.TranslateFunc {
	T, _ := i18n.Tfunc(language, languages...)
	tmpl.Funcs(map[string]interface{}{
		"T": func(str string, args ...interface{}) template.HTML {
			return template.HTML(fmt.Sprintf(T(str), args...))
		},
	})
	return T
}

func SetTranslationFromRequest(tmpl *template.Template, r *http.Request, defaultLanguage string) i18n.TranslateFunc {
	cookie, err := r.Cookie("lang")
	cookieLanguage := ""
	if err == nil {
		cookieLanguage = cookie.Value
	}
	// go-i18n supports the format of the Accept-Language header, thankfully.
	headerLanguage := r.Header.Get("Accept-Language")
	return SetTranslation(tmpl, cookieLanguage, headerLanguage, defaultLanguage)
}
