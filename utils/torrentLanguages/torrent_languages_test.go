package torrentLanguages

import (
	"io/ioutil"
	"path"
	"strings"
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

func TestCSSFlags(t *testing.T) {
	languages := GetTorrentLanguages()
	flagsCSSPath := path.Join("..", "..", "public", "css", "flags", "flags.css")
	file, err := ioutil.ReadFile(flagsCSSPath)
	if err != nil {
		t.Errorf("Failed to load flags.css: %v", err)
		return
	}

	contents := string(file)
	for _, language := range languages {
		flag := FlagFromLanguage(language)
		if !strings.Contains(contents, ".flag-"+flag) {
			t.Errorf("flags.css does not contains class .flag-%s. You probably need to update it.", flag)
		}
	}
}
