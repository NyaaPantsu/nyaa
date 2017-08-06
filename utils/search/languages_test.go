package search

import (
	"reflect"
	"testing"

	"github.com/NyaaPantsu/nyaa/utils/publicSettings"
	"github.com/stretchr/testify/assert"
)

func TestParseLanguages(t *testing.T) {
	langs := ParseLanguages([]string{""})
	if len(langs) > 0 {
		t.Fatal("ParseLanguages with empty arg doesn't return an empty array")
	}
	langs = ParseLanguages([]string{})
	if len(langs) > 0 {
		t.Fatal("ParseLanguages with empty array doesn't return an empty array")
	}
	langs = ParseLanguages([]string{"fr"})
	if len(langs) == 0 {
		t.Fatal("ParseLanguages with good arg return an empty array")
	}
	langs = ParseLanguages([]string{"en,fr"})
	langEqual := publicSettings.Languages{
		publicSettings.Language{
			Name: "",
			Code: "en",
		},
		publicSettings.Language{
			Name: "",
			Code: "fr",
		},
	}
	if !reflect.DeepEqual(langs, langEqual) {
		t.Fatal("ParseLanguages with good arg doesn't return the right array")
	}
	langs = ParseLanguages([]string{"en,,,,fr"})
	if !reflect.DeepEqual(langs, langEqual) {
		t.Fatal("ParseLanguages doesn't remove empty values")
	}
	langs = ParseLanguages([]string{"en", "fr"})
	if !reflect.DeepEqual(langs, langEqual) {
		t.Fatal("ParseLanguages with good arg doesn't return the right array")
	}
	langs = ParseLanguages([]string{"en", "", "", "fr"})
	if !reflect.DeepEqual(langs, langEqual) {
		t.Fatal("ParseLanguages doesn't remove empty values")
	}
}

func TestSplitsLanguages(t *testing.T) {
	assert := assert.New(t)
	expect := publicSettings.Languages{{Code: "fr"}}

	assert.Empty(splitsLanguages(""), "Should be empty")
	assert.Empty(splitsLanguages(","), "Should be empty")
	assert.Equal(expect, splitsLanguages(",fr"), "Should be empty")
}
