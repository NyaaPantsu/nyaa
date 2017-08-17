package tagsValidator

import (
	"path"
	"testing"

	"github.com/NyaaPantsu/nyaa/config"
)

// run before config/parse.go:init()
var _ = func() (_ struct{}) {
	config.Configpaths[1] = path.Join("..", "..", "..", config.Configpaths[1])
	config.Configpaths[0] = path.Join("..", "..", "..", config.Configpaths[0])
	config.Reload()
	return
}()

func TestCheck(t *testing.T) {
	tests := []struct {
		Tag      []string
		Expected bool
	}{
		{[]string{"", ""}, false},
		{[]string{"akuma06", ""}, false},
		{[]string{"videoquality", "full_hd"}, true},
		{[]string{"anidbid", ""}, false},
		{[]string{"anidbid", "20"}, true},
	}
	for _, test := range tests {
		b := Check(test.Tag[0], test.Tag[1])
		if b != test.Expected {
			t.Errorf("Error when checking tag type '%v', want '%t', got '%t'", test.Tag, test.Expected, b)
		}
	}
}
