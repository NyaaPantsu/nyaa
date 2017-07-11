package publicSettings

import (
	"path"
	"testing"

	"github.com/nicksnyder/go-i18n/i18n"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"

	"fmt"

	"strings"

	"github.com/NyaaPantsu/nyaa/config"
)

// run before config/parse.go:init()
var _ = func() (_ struct{}) {
	config.ConfigPath = path.Join("..", "..", config.ConfigPath)
	config.DefaultConfigPath = path.Join("..", "..", config.DefaultConfigPath)
	config.Reload()
	return
}()

func TestInitI18n(t *testing.T) {
	conf := config.Get().I18n
	conf.Directory = path.Join("..", "..", conf.Directory)
	var retriever UserRetriever // not required during initialization

	err := InitI18n(conf, retriever)
	if err != nil {
		t.Errorf("failed to initialize language translations: %v", err)
	}
}

func TestLanguages(t *testing.T) {
	conf := config.Get().I18n
	conf.Directory = path.Join("..", "..", conf.Directory)
	var retriever UserRetriever // not required during initialization
	err := InitI18n(conf, retriever)
	if err != nil {
		t.Errorf("failed to initialize language translations: %v", err)
	}
	displayLang := language.Make(conf.DefaultLanguage)
	if displayLang.String() == "und" {
		t.Errorf("Couldn't find the language display for the language %s", displayLang.String())
	}
	n := display.Tags(displayLang)

	tags := i18n.LanguageTags()
	for _, languageTag := range tags {
		// The matcher will match Swiss German to German.
		lang := getParentTag(languageTag)
		if lang.String() == "und" {
			t.Errorf("Couldn't find the language root for the language %s", languageTag)
		}
		fmt.Printf("Name of the language natively: %s\n", strings.Title(display.Self.Name(lang)))
		fmt.Printf("Name of the language in %s: %s\n", displayLang.String(), n.Name(lang))
	}
}

func TestTranslate(t *testing.T) {
	conf := config.Get().I18n
	conf.Directory = path.Join("..", "..", conf.Directory)
	var retriever UserRetriever // not required during initialization
	err := InitI18n(conf, retriever)
	if err != nil {
		t.Errorf("failed to initialize language translations: %v", err)
	}

	T, _ := GetDefaultTfunc()
	test := []map[string]string{
		{
			"test":   "",
			"result": "",
		},
		{
			"test":   "fr-fr",
			"result": "French (France)",
		},
		{
			"test":   "fr",
			"result": "French",
		},
		{
			"test":   "fredfef",
			"result": "",
		},
	}
	for _, langTest := range test {
		result := Translate(langTest["test"], T("language_code"))
		if result != langTest["result"] {
			t.Errorf("Result from Translate function different from the expected: have '%s', wants '%s'", result, langTest["result"])
		}
	}
}
