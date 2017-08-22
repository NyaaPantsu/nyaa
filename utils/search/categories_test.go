package search

import (
	"path"
	"testing"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/utils/categories"
	"github.com/stretchr/testify/assert"
)

// run before config/parse.go:init()
var _ = func() (_ struct{}) {
	config.Configpaths[1] = path.Join("..", "..", config.Configpaths[1])
	config.Configpaths[0] = path.Join("..", "..", config.Configpaths[0])
	config.Reload()
	categories.InitCategories()
	return
}()

func TestParseCategories(t *testing.T) {
	assert := assert.New(t)
	cat := ParseCategories("")
	assert.Empty(cat, "ParseCategories with empty arg doesn't return an empty array")
	cat = ParseCategories("5050")
	assert.Empty(cat, "ParseCategories with wrong arg doesn't return an empty array")

	cat = ParseCategories("50_50")
	assert.Empty(cat, "ParseCategories with wrong arg doesn't return an empty array")

	cat = ParseCategories("3_13")
	catEqual := []*Category{
		&Category{
			Main: 3,
			Sub:  13,
		},
	}
	assert.Equal(catEqual, cat, "ParseCategories with good arg doesn't return the right array")

	cat = ParseCategories("_")
	assert.Empty(cat, "Should be empty")

	cat = ParseCategories("3_13,3_5")
	catEqual = []*Category{
		&Category{
			Main: 3,
			Sub:  13,
		},
		&Category{
			Main: 3,
			Sub:  5,
		},
	}
	assert.Equal(catEqual, cat, "ParseCategories with good arg doesn't return the right array")

	cat = ParseCategories("3_13,3_5,5_50")
	assert.Equal(catEqual, cat, "ParseCategories with good arg doesn't filter the wrong categories")
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

func TestCategories_ToDBQuery(t *testing.T) {
	cat := Categories{&Category{
		Main: 3,
		Sub:  13,
	},
		&Category{
			Main: 3,
			Sub:  5,
		},
	}
	assert := assert.New(t)
	search, args := cat.ToDBQuery()
	assert.Equal("(category = ? AND sub_category = ?) OR (category = ? AND sub_category = ?)", search, "Should be equal")
	assert.Equal([]interface{}{uint8(3), uint8(13), uint8(3), uint8(5)}, args, "Should be equal")
	cat = Categories{&Category{
		Main: 3,
		Sub:  13,
	},
	}
	search, args = cat.ToDBQuery()
	assert.Equal("category = ? AND sub_category = ?", search, "Should be equal to 'category = ? AND sub_category = ?'")
	assert.Equal([]interface{}{uint8(3), uint8(13)}, args, "Should be equal to '3_13'")
	cat = Categories{&Category{
		Main: 3,
		Sub:  0,
	},
	}
	search, args = cat.ToDBQuery()
	assert.Equal("category = ?", search, "Should be equal to 'category = ?'")
	assert.Equal([]interface{}{uint8(3)}, args, "Should be equal to '3'")
	cat = Categories{&Category{
		Main: 0,
		Sub:  0,
	},
	}
	search, args = cat.ToDBQuery()
	assert.Empty(search, "Should be empty")
	assert.Empty(args, "Should be empty")
}

func TestCategories_ToESQuery(t *testing.T) {
	cat := Categories{&Category{
		Main: 3,
		Sub:  13,
	},
		&Category{
			Main: 3,
			Sub:  5,
		},
	}
	assert := assert.New(t)
	assert.Equal("((category: 3 AND sub_category: 13) OR (category: 3 AND sub_category: 5))", cat.ToESQuery(), "Should be equal")
	cat = Categories{&Category{
		Main: 3,
		Sub:  13,
	},
	}
	assert.Equal("(category: 3 AND sub_category: 13)", cat.ToESQuery(), "Should be equal")
	cat = Categories{&Category{
		Main: 3,
		Sub:  0,
	},
	}
	assert.Equal("(category: 3)", cat.ToESQuery(), "Should be equal")
	cat = Categories{&Category{
		Main: 0,
		Sub:  0,
	},
	}
	assert.Equal("", cat.ToESQuery(), "Should be empty")
}
