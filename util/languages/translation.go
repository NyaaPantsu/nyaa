package languages

import (
	"fmt"
	"github.com/ewhal/nyaa/service/user"
	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/nicksnyder/go-i18n/i18n/language"
	"html/template"
	"net/http"
)

// When go-i18n finds a language with >0 translations, it uses it as the Tfunc
// However, if said language has a missing translation, it won't fallback to the "main" language
func TfuncAndLanguageWithFallback(language string, languages ...string) (i18n.TranslateFunc, *language.Language, error) {
	// Use the last language on the args as the fallback one.
	fallbackLanguage := language
	if languages != nil {
		fallbackLanguage = languages[len(languages)-1]
	}

	T, Tlang, err1 := i18n.TfuncAndLanguage(language, languages...)
	fallbackT, fallbackTlang, err2 := i18n.TfuncAndLanguage(fallbackLanguage)

	if err1 != nil && err2 != nil {
		// fallbackT is still a valid function even with the error, it returns translationID.
		return fallbackT, fallbackTlang, err2
	}

	return func(translationID string, args ...interface{}) string {
		if translated := T(translationID, args...); translated != translationID {
			return translated
		}

		return fallbackT(translationID, args...)
	}, Tlang, nil
}

func GetAvailableLanguages() (languages map[string]string) {
	languages = make(map[string]string)
	var T i18n.TranslateFunc
	for _, languageTag := range i18n.LanguageTags() {
		T, _ = i18n.Tfunc(languageTag)
		/* Translation files should have an ID with the translated language name.
		   If they don't, just use the languageTag */
		if languageName := T("language_name"); languageName != "language_name" {
			languages[languageTag] = languageName
		} else {
			languages[languageTag] = languageTag
		}
	}
	return
}

func setTranslation(tmpl *template.Template, T i18n.TranslateFunc) {
	tmpl.Funcs(map[string]interface{}{
		"T": func(str string, args ...interface{}) template.HTML {
			return template.HTML(fmt.Sprintf(T(str), args...))
		},
		"Ts": func(str string, args ...interface{}) string {
			return fmt.Sprintf(T(str), args...)
		},
	})
}

func GetTfuncAndLanguageFromRequest(r *http.Request, defaultLanguage string) (T i18n.TranslateFunc, Tlang *language.Language) {
	userLanguage := ""
	user, _, err := userService.RetrieveCurrentUser(r)
	if err == nil {
		userLanguage = user.Language
	}

	cookie, err := r.Cookie("lang")
	cookieLanguage := ""
	if err == nil {
		cookieLanguage = cookie.Value
	}

	// go-i18n supports the format of the Accept-Language header, thankfully.
	headerLanguage := r.Header.Get("Accept-Language")
	T, Tlang, _ = TfuncAndLanguageWithFallback(userLanguage, cookieLanguage, headerLanguage, defaultLanguage)
	return
}


func SetTranslationFromRequest(tmpl *template.Template, r *http.Request, defaultLanguage string) i18n.TranslateFunc {
	r.Header.Add("Vary", "Accept-Encoding")
	T, _ := GetTfuncAndLanguageFromRequest(r, defaultLanguage)
	setTranslation(tmpl, T)
	return T
}
