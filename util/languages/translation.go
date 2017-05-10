package languages

import (
	"fmt"
	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/ewhal/nyaa/service/user"
	"html/template"
	"net/http"
)

// When go-i18n finds a language with >0 translations, it uses it as the Tfunc
// However, if said language has a missing translation, it won't fallback to the "main" language
func TfuncWithFallback(language string, languages ...string) (i18n.TranslateFunc, error) {
	// Use the last language on the args as the fallback one.
	fallbackLanguage := language
	if languages != nil {
		fallbackLanguage = languages[len(languages)-1]
	}

	T, err1 := i18n.Tfunc(language, languages...)
	fallbackT, err2 := i18n.Tfunc(fallbackLanguage)

	if err1 != nil && err2 != nil {
		// fallbackT is still a valid function even with the error, it returns translationID.
		return fallbackT, err2;
	}

	return func(translationID string, args ...interface{}) string {
		if translated := T(translationID, args...); translated != translationID {
			return translated
		}

		return fallbackT(translationID, args...)
	}, nil
}

func GetAvailableLanguages() (languages map[string]string) {
	languages = make(map[string]string)
	var T i18n.TranslateFunc
	for _, languageTag := range i18n.LanguageTags() {
		T, _ = i18n.Tfunc(languageTag)
		/* Translation files should have an ID with the translated language name.
		   If they don't, just use the languageTag */
		if languageName := T("language_name"); languageName != "language_name" {
			languages[languageTag] = languageName;
		} else {
			languages[languageTag] = languageTag
		}
	}
	return
}

func SetTranslation(tmpl *template.Template, language string, languages ...string) i18n.TranslateFunc {
	T, _ := TfuncWithFallback(language, languages...)
	tmpl.Funcs(map[string]interface{}{
		"T": func(str string, args ...interface{}) template.HTML {
			return template.HTML(fmt.Sprintf(T(str), args...))
		},
	})
	return T
}

func SetTranslationFromRequest(tmpl *template.Template, r *http.Request, defaultLanguage string) i18n.TranslateFunc {
	userLanguage := ""
	user, _, err := userService.RetrieveCurrentUser(r)
	if err == nil {
		userLanguage = user.Language;
	}

	cookie, err := r.Cookie("lang")
	cookieLanguage := ""
	if err == nil {
		cookieLanguage = cookie.Value
	}

	// go-i18n supports the format of the Accept-Language header, thankfully.
	headerLanguage := r.Header.Get("Accept-Language")
	return SetTranslation(tmpl, userLanguage, cookieLanguage, headerLanguage, defaultLanguage)
}
