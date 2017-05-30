package publicSettings

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"path/filepath"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/nicksnyder/go-i18n/i18n/language"
)

// UserRetriever : this interface is required to prevent a cyclic import between the languages and userService package.
type UserRetriever interface {
	RetrieveCurrentUser(r *http.Request) (model.User, error)
}

// TemplateTfunc : T func used in template
type TemplateTfunc func(string, ...interface{}) template.HTML

var (
	defaultLanguage = config.DefaultI18nConfig.DefaultLanguage
	userRetriever   UserRetriever
)

// InitI18n : Initialize the languages translation
func InitI18n(conf config.I18nConfig, retriever UserRetriever) error {
	defaultLanguage = conf.DefaultLanguage
	userRetriever = retriever

	defaultFilepath := path.Join(conf.TranslationsDirectory, defaultLanguage+".all.json")
	err := i18n.LoadTranslationFile(defaultFilepath)
	if err != nil {
		panic(fmt.Sprintf("failed to load default translation file '%s': %v", defaultFilepath, err))
	}

	paths, err := filepath.Glob(path.Join(conf.TranslationsDirectory, "*.json"))
	if err != nil {
		return fmt.Errorf("failed to get translation files: %v", err)
	}

	for _, file := range paths {
		err := i18n.LoadTranslationFile(file)
		if err != nil {
			return fmt.Errorf("failed to load translation file '%s': %v", file, err)
		}
	}

	return nil
}

// GetDefaultLanguage : returns the default language from config
func GetDefaultLanguage() string {
	return defaultLanguage
}

// TfuncAndLanguageWithFallback : When go-i18n finds a language with >0 translations, it uses it as the Tfunc
// However, if said language has a missing translation, it won't fallback to the "main" language
func TfuncAndLanguageWithFallback(language string, languages ...string) (i18n.TranslateFunc, *language.Language, error) {
	fallbackLanguage := GetDefaultLanguage()

	tFunc, tLang, err1 := i18n.TfuncAndLanguage(language, languages...)
	// If fallbackLanguage fails, it will give the "id" field so we don't
	// care about the error
	fallbackT, fallbackTlang, _ := i18n.TfuncAndLanguage(fallbackLanguage)

	translateFunction := func(translationID string, args ...interface{}) string {
		if translated := tFunc(translationID, args...); translated != translationID {
			return translated
		}

		return fallbackT(translationID, args...)
	}

	if err1 != nil {
		tLang = fallbackTlang
	}

	return translateFunction, tLang, err1
}

// GetAvailableLanguages : Get languages available on the website
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

// GetDefaultTfunc : Gets T func from default language
func GetDefaultTfunc() (i18n.TranslateFunc, error) {
	return i18n.Tfunc(defaultLanguage)
}

// GetTfuncAndLanguageFromRequest : Gets the T func and chosen language from the request
func GetTfuncAndLanguageFromRequest(r *http.Request) (T i18n.TranslateFunc, Tlang *language.Language) {
	userLanguage := ""
	user, _ := getCurrentUser(r)
	if user.ID > 0 {
		userLanguage = user.Language
	}

	cookie, err := r.Cookie("lang")
	cookieLanguage := ""
	if err == nil {
		cookieLanguage = cookie.Value
	}

	// go-i18n supports the format of the Accept-Language header
	headerLanguage := r.Header.Get("Accept-Language")
	T, Tlang, _ = TfuncAndLanguageWithFallback(userLanguage, cookieLanguage, headerLanguage)
	return
}

// GetTfuncFromRequest : Gets the T func from the request
func GetTfuncFromRequest(r *http.Request) TemplateTfunc {
	T, _ := GetTfuncAndLanguageFromRequest(r)
	return func(id string, args ...interface{}) template.HTML {
		return template.HTML(fmt.Sprintf(T(id), args...))
	}
}

// GetThemeFromRequest: Gets the user selected theme from the request
func GetThemeFromRequest(r *http.Request) string {
	user, _ := getCurrentUser(r)
	if user.ID > 0 {
		return user.Theme
	}
	cookie, err := r.Cookie("theme")
	if err == nil {
		return cookie.Value
	}
	return ""
}

// GetThemeFromRequest: Gets the user selected theme from the request
func GetMascotFromRequest(r *http.Request) string {
	user, _ := getCurrentUser(r)
	if user.ID > 0 {
		return user.Mascot
	}
	cookie, err := r.Cookie("mascot")
	if err == nil {
		return cookie.Value
	}
	return "show"
}


func getCurrentUser(r *http.Request) (model.User, error) {
	if userRetriever == nil {
		return model.User{}, errors.New("failed to get current user: no user retriever set")
	}
	return userRetriever.RetrieveCurrentUser(r)
}
