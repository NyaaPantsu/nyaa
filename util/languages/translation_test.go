package languages

import (
	"path"
	"testing"

	"github.com/ewhal/nyaa/config"
)

func TestInitI18n(t *testing.T) {
	conf := config.DefaultI18nConfig
	conf.TranslationsDirectory = path.Join("..", "..", conf.TranslationsDirectory)

	err := InitI18n(conf)
	if err != nil {
		t.Errorf("failed to initialize language translations: %v", err)
	}
}
