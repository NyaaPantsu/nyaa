package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortMode_Parse(t *testing.T) {
	assert := assert.New(t)
	sort := SortMode(0)
	sort.Parse("")
	assert.Equal(Date, sort, "Should be default to Date")
	sort.Parse("lol10")
	assert.Equal(Date, sort, "Should be default to Date")
	sort.Parse("1")
	assert.Equal(SortMode(1), sort)
	sort.Parse("10")
	assert.Equal(Date, sort)
}

func TestSortMode_ToDBField(t *testing.T) {
	assert := assert.New(t)
	sort := SortMode(0)
	sort.Parse("")
	assert.Equal("date", sort.ToDBField(), "Should be default to Date")
	sort.Parse("lol10")
	assert.Equal("date", sort.ToDBField(), "Should be default to Date")
	sort.Parse("1")
	assert.Equal("torrent_name", sort.ToDBField())
	sort.Parse("10")
	assert.Equal("date", sort.ToDBField())
}

func TestSortMode_ToESField(t *testing.T) {
	assert := assert.New(t)
	sort := SortMode(0)
	sort.Parse("")
	assert.Equal("date", sort.ToESField(), "Should be default to Date")
	sort.Parse("lol10")
	assert.Equal("date", sort.ToESField(), "Should be default to Date")
	sort.Parse("1")
	assert.Equal("name.raw", sort.ToESField())
	sort.Parse("10")
	assert.Equal("date", sort.ToESField())
}
