package torrentLanguages

import (
	"io/ioutil"
	"path"
	"strings"
	"testing"

	"fmt"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/utils/publicSettings"
)

// run before config/parse.go:init()
var _ = func() (_ struct{}) {
	config.ConfigPath = path.Join("..", "..", config.ConfigPath)
	config.DefaultConfigPath = path.Join("..", "..", config.DefaultConfigPath)
	config.Reload()
	config.Get().I18n.Directory = path.Join("..", "..", config.Get().I18n.Directory)
	return
}()

func TestCSSFlags(t *testing.T) {
	var retriever publicSettings.UserRetriever // not required during initialization
	err := publicSettings.InitI18n(config.Get().I18n, retriever)
	if err != nil {
		t.Errorf("failed to initialize language translations: %v", err)
	}
	languages := GetTorrentLanguages()
	flagsCSSPath := path.Join("..", "..", "public", "css", "flags", "flags.css")
	file, err := ioutil.ReadFile(flagsCSSPath)
	if err != nil {
		t.Errorf("Failed to load flags.css: %v", err)
		return
	}

	contents := string(file)
	for _, language := range languages {
		flag := publicSettings.Flag(language.Code, true)
		fmt.Printf("Finding css class for: %s (%s)\n", flag, language.Name)
		if !strings.Contains(contents, ".flag-"+flag) {
			t.Errorf("flags.css does not contains class .flag-%s. You probably need to update it.", flag)
		}
	}
}
