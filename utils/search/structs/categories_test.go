package structs

import (
	"path"
	"reflect"
	"testing"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/utils/categories"
	"github.com/stretchr/testify/assert"
)

// run before config/parse.go:init()
var _ = func() (_ struct{}) {
	config.Configpaths[1] = path.Join("..", "..", "..", config.Configpaths[1])
	config.Configpaths[0] = path.Join("..", "..", "..", config.Configpaths[0])
	config.Reload()
	categories.InitCategories()
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

func TestCategory_IsSet(t *testing.T) {
	cat := Category{
		Main: 3,
		Sub:  13,
	}
	assert := assert.New(t)
	assert.Equal(true, cat.IsSet(), "Should be equal to true")
	cat.Main = 0
	assert.Equal(false, cat.IsSet(), "Should be equal to false")
	cat.Sub = 0
	assert.Equal(false, cat.IsSet(), "Should be equal to false")
}

func TestCategory_IsSubSet(t *testing.T) {
	cat := Category{
		Main: 3,
		Sub:  13,
	}
	assert := assert.New(t)
	assert.Equal(true, cat.IsSubSet(), "Should be equal to true")
	cat.Main = 0
	assert.Equal(true, cat.IsSubSet(), "Should be equal to true")
	cat.Sub = 0
	assert.Equal(false, cat.IsSubSet(), "Should be equal to false")
}

func TestCategory_IsMainSet(t *testing.T) {
	cat := Category{
		Main: 3,
		Sub:  13,
	}
	assert := assert.New(t)
	assert.Equal(true, cat.IsMainSet(), "Should be equal to true")
	cat.Main = 0
	assert.Equal(false, cat.IsMainSet(), "Should be equal to false")
	cat.Sub = 0
	assert.Equal(false, cat.IsMainSet(), "Should be equal to false")
}

func TestCategory_String(t *testing.T) {
	cat := Category{
		Main: 3,
		Sub:  13,
	}
	assert := assert.New(t)
	assert.Equal("3_13", cat.String(), "Should be equal to '3_13'")
	cat.Sub = 0
	assert.Equal("3_", cat.String(), "Should be equal to '3_'")
	cat.Main = 0
	assert.Equal("_", cat.String(), "Should be equal to '_'")
}
