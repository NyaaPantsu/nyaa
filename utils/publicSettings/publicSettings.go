package publicSettings

import (
	"errors"
	"fmt"
	"html/template"
	"path"
	"path/filepath"
	"strings"
	"net/http"

	"sort"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/utils/timeHelper"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/nicksnyder/go-i18n/i18n/language"

	glang "golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

// UserRetriever : this interface is required to prevent a cyclic import between the languages and userService package.
type UserRetriever interface {
	RetrieveCurrentUser(c *gin.Context) (*models.User, error)
}

// Language localization language struct
type Language struct {
	Name string
	Code string
	Tag  string
}

// Languages Array of Language
type Languages []Language

// TemplateTfunc : T func used in template
type TemplateTfunc func(string, ...interface{}) template.HTML

var (
	defaultLanguage = config.Get().I18n.DefaultLanguage
	userRetriever   UserRetriever
	languages       Languages
)

// InitI18n : Initialize the languages translation
func InitI18n(conf config.I18nConfig, retriever UserRetriever) error {
	defaultLanguage = conf.DefaultLanguage
	userRetriever = retriever

	defaultFilepath := path.Join(conf.Directory, defaultLanguage+".all.json")
	err := i18n.LoadTranslationFile(defaultFilepath)
	if err != nil {
		panic(fmt.Sprintf("failed to load default translation file '%s': %v", defaultFilepath, err))
	}

	paths, err := filepath.Glob(path.Join(conf.Directory, "*.json"))
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

// GetAvailableLanguages : Get languages available on the website, languages are parsed once at runtime
func GetAvailableLanguages() Languages {
	if len(languages) > 0 {
		return languages
	}
	// Need this to sort out languages alphabetically by language tag
	var codes []string
	for _, languageTag := range i18n.LanguageTags() {
		codes = append(codes, languageTag)
	}
	languages = ParseLanguages(codes)
	return languages
}

// GetParentTag returns the highest parent of a language (e.g. fr-fr -> fr)
func GetParentTag(languageTag string) glang.Tag {
	lang := glang.Make(languageTag)
	for !lang.Parent().IsRoot() {
		lang = lang.Parent()
	}
	return lang
}

// ParseLanguages takes a list of language codes and convert them in languages object
func ParseLanguages(codes []string) Languages {
	var langs Languages
	sort.Strings(codes)
	// Now build languages array
	for _, languageTag := range codes {
		lang := GetParentTag(languageTag)
		langs = append(langs, Language{strings.Title(display.Self.Name(glang.Make(languageTag))), lang.String(), languageTag})
	}
	return langs
}

// GetDefaultTfunc : Gets T func from default language
func GetDefaultTfunc() (i18n.TranslateFunc, error) {
	return i18n.Tfunc(defaultLanguage)
}

// GetTfuncAndLanguageFromRequest : Gets the T func and chosen language from the request
func GetTfuncAndLanguageFromRequest(c *gin.Context) (T i18n.TranslateFunc, Tlang *language.Language) {
	userLanguage := ""
	user, _ := getCurrentUser(c)
	if user.ID > 0 {
		userLanguage = user.Language
	}

	cookie, err := c.Cookie("lang")
	cookieLanguage := ""
	if err == nil {
		cookieLanguage = cookie
	}

	// go-i18n supports the format of the Accept-Language header
	headerLanguage := c.Request.Header.Get("Accept-Language")
	T, Tlang, _ = TfuncAndLanguageWithFallback(userLanguage, cookieLanguage, headerLanguage)
	return
}

// GetTfuncFromRequest : Gets the T func from the request
func GetTfuncFromRequest(c *gin.Context) TemplateTfunc {
	T, _ := GetTfuncAndLanguageFromRequest(c)
	return func(id string, args ...interface{}) template.HTML {
		return template.HTML(fmt.Sprintf(T(id), args...))
	}
}

// GetThemeFromRequest : Gets the user selected theme from the request
func GetThemeFromRequest(c *gin.Context) string {
	user, _ := getCurrentUser(c)
	if user.ID > 0 {
		return user.Theme
	}
	cookie, err := c.Cookie("theme")
	if err == nil {
		return cookie
	}
	return ""
}

// GetAltColorsFromRequest : Return whether user has enabled alt colors or not
func GetAltColorsFromRequest(c *gin.Context) bool {
	user, _ := getCurrentUser(c)
	if user.ID > 0 {
		return user.AltColors != "false"
		//Doing this in order to make it return true should the field be empty
	}
	cookie, err := c.Cookie("altColors")
	if err == nil {
		return cookie == "true"
	}
	return true
}

// GetOldNavFromRequest : Return whether user has enabled old navigation or not
func GetOldNavFromRequest(c *gin.Context) bool {
	user, _ := getCurrentUser(c)
	if user.ID > 0 {
		return user.OldNav != "false"
		//Doing this in order to make it return true should the field be empty
	}
	cookie, err := c.Cookie("oldNav")
	if err == nil {
		return cookie == "true"
	}
	return true
}

// GetMascotFromRequest : Return whether user has enabled mascot or not
func GetMascotFromRequest(c *gin.Context) string {
	user, _ := getCurrentUser(c)
	if user.ID > 0 {
		return user.Mascot
	}
	cookie, err := c.Cookie("mascot")
	if err == nil {
		return cookie
	}
	return "show"
}

// GetMascotURLFromRequest : Get the user selected mascot url from the request.
// Returns an empty string if not set.
func GetMascotURLFromRequest(c *gin.Context) string {
	user, _ := getCurrentUser(c)
	if user.ID > 0 {
		return user.MascotURL
	}

	cookie, err := c.Cookie("mascot_url")
	if err == nil {
		return cookie
	}

	return ""
}

func GetEUCookieFromRequest(c *gin.Context) bool {
	_, err := c.Cookie("EU_Cookie")
	if err == nil {
		return true
		//Cookie exists, everything good
	}

	http.SetCookie(c.Writer, &http.Cookie{Name: "EU_Cookie", Value: "true", Domain: getDomainName(), Path: "/", Expires: timeHelper.FewDaysLater(365)})	
	return false
	//Cookie doesn't exist, we create it to prevent the message from popping up anymore after that and return false
}

func getCurrentUser(c *gin.Context) (*models.User, error) {
	if userRetriever == nil {
		return &models.User{}, errors.New("failed to get current user: no user retriever set")
	}
	return userRetriever.RetrieveCurrentUser(c)
}

// Exist evaluate if a language exists or not (language code or language name)
func (langs Languages) Exist(name string) bool {
	for _, language := range langs {
		if language.Code == name || language.Name == name {
			return true
		}
	}

	return false
}

// Translate accepts a languageCode in string and translate the language to the language from the language code
func (lang *Language) Translate(languageCode template.HTML) string {
	return Translate(lang.Tag, string(languageCode))
}

// Translate accepts a languageCode in string and translate the language to the language from the language code in to
func Translate(languageCode string, to string) string {
	langTranslate := display.Tags(GetParentTag(to))
	translated := langTranslate.Name(glang.Make(languageCode))
	if translated == "Root" {
		return ""
	}
	return translated
}

// Flag reads the language's country code and return the country's flag if national true or the international flag for the language
func (lang *Language) Flag(national bool) string {
	if !national {
		if !strings.Contains(lang.Tag, ",") {
			return Flag(lang.Tag, false)
		}
		return Flag(lang.Code, false)
	}
	return lang.Code
}

// Flag reads the language's country code and return the country's flag if national true or the international flag for the language
func Flag(languageCode string, parent bool) string {
	lang := glang.Make(languageCode)
	if parent {
		lang = GetParentTag(languageCode)
	}

	languageSplit := strings.Split(lang.String(), "-")
	if len(languageSplit) > 1 {
		return languageSplit[1]
	}
	return lang.String()
}

func getDomainName() string {
	domain := config.Get().Cookies.DomainName
	if config.Get().Environment == "DEVELOPMENT" {
		domain = ""
	}
	return domain
}
