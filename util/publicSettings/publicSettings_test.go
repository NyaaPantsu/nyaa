package publicSettings

import (
	"path"
	"testing"

	"github.com/NyaaPantsu/nyaa/config"
)

func TestInitI18n(t *testing.T) {
	conf := config.Conf.I18n
	conf.Directory = path.Join("..", "..", conf.Directory)
	var retriever UserRetriever // not required during initialization

	err := InitI18n(conf, retriever)
	if err != nil {
		t.Errorf("failed to initialize language translations: %v", err)
	}
}
