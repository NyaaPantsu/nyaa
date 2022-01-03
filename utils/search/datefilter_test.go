package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDateFilter_Parse(t *testing.T) {
	assert := assert.New(t)
	date := DateFilter("")
	assert.Equal(false, date.Parse(""), "Should be false")
	assert.Empty(string(date), "Should be empty")

	assert.Equal(false, date.Parse("lol"), "Should be false")
	assert.Empty(string(date), "Should be empty")

	assert.Equal(false, date.Parse("05/06/1486"), "Should be false")
	assert.Empty(string(date), "Should be empty")

	assert.Equal(true, date.Parse("2017/08/01"), "Should be true")
	assert.Equal("2017-08-01", string(date), "Should be empty")
}
