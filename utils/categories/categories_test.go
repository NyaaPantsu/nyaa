package categories

import (
	"path"
	"strings"

	"testing"

	"reflect"

	"github.com/NyaaPantsu/nyaa/config"
)

// run before config/parse.go:init()
var _ = func() (_ struct{}) {
	config.ConfigPath = path.Join("..", "..", config.ConfigPath)
	config.DefaultConfigPath = path.Join("..", "..", config.DefaultConfigPath)
	config.Parse()
	return
}()

func TestGetCategories(t *testing.T) {
	cats := make(map[string]string)

	for _, v := range All() {
		cats[v.ID] = v.Name
	}
	if len(cats) == 0 {
		t.Skip("Couldn't load categories to test Categories")
	}

	if !reflect.DeepEqual(cats, config.Conf.Torrents.CleanCategories) && !reflect.DeepEqual(cats, config.Conf.Torrents.SukebeiCategories) {
		t.Error("Categories doesn't correspond to the configured ones")
	}
}

func TestCategoryExists(t *testing.T) {
	if Exists("k") {
		t.Error("Category that shouldn't exist return true")
	}
}

func TestGetCategoriesSelect(t *testing.T) {
	cats := GetSelect(true, false)
	for _, value := range cats {
		split := strings.Split(value.ID, "_")
		if len(split) != 2 {
			t.Errorf("The category %s doesn't have only one underscore", value.Name)
		}
		if split[1] != "" {
			t.Errorf("The function doesn't filter out child categories, expected '', got %s", split[1])
		}
	}
	cats = GetSelect(false, true)
	for _, value := range cats {
		split := strings.Split(value.ID, "_")
		if len(split) != 2 {
			t.Errorf("The category %s doesn't have only one underscore", value.Name)
		}
		if split[1] == "" {
			t.Error("The function doesn't filter out parent categories, expected a string, got nothing")
		}
	}
	cats = GetSelect(true, true)
	if len(cats) != len(All()) {
		t.Errorf("Same amount of categories isn't return when no filter applied")
	}
}
