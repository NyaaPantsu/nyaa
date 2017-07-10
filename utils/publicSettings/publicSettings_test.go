package publicSettings

import (
	"path"
	"testing"

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
