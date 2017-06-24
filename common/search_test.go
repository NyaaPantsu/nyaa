package common

import (
	"path"
	"testing"

	"reflect"

	"github.com/NyaaPantsu/nyaa/config"
)

// run before config/parse.go:init()
var _ = func() (_ struct{}) {
	config.ConfigPath = path.Join("..", config.ConfigPath)
	config.DefaultConfigPath = path.Join("..", config.DefaultConfigPath)
	config.Parse()
	return
}()

func TestParseCategories(t *testing.T) {
	cat := ParseCategories("")
	if len(cat) > 0 {
		t.Fatal("ParseCategories with empty arg doesn't return an empty array")
	}
	cat = ParseCategories("5050")
	if len(cat) > 0 {
		t.Fatal("ParseCategories with wrong arg doesn't return an empty array")
	}
	cat = ParseCategories("50_50")
	if len(cat) > 0 {
		t.Fatal("ParseCategories with wrong arg doesn't return an empty array")
	}
	cat = ParseCategories("3_13")
	if len(cat) == 0 {
		t.Fatal("ParseCategories with good arg return an empty array")
	}
	cat = ParseCategories("3_13,3_5")
	catEqual := []*Category{
		&Category{
			Main: 3,
			Sub:  13,
		},
		&Category{
			Main: 3,
			Sub:  5,
		},
	}
	if !reflect.DeepEqual(cat, catEqual) {
		t.Fatal("ParseCategories with good arg doesn't return the right array")
	}
	cat = ParseCategories("3_13,3_5,5_50")
	if !reflect.DeepEqual(cat, catEqual) {
		t.Fatal("ParseCategories doesn't filter the wrong categories")
	}

}
