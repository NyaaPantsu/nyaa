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

func TestCheckTagType(t *testing.T) {
	tests := []struct {
		Type     string
		Expected bool
	}{
		{"", false},
		{"akuma06", false},
		{"quality", true},
		{"anidb", false},
	}
	for _, test := range tests {
		b := CheckTagType(test.Type)
		if b != test.Expected {
			t.Errorf("Error when checking tag type '%s', want '%t', got '%t'", test.Type, test.Expected, b)
		}
	}
}
