package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSizeBytes_Parse(t *testing.T) {
	assert := assert.New(t)
	size := SizeBytes(0)
	assert.Equal(false, size.Parse("", ""), "Should be false")
	assert.Zero(size, "Should be zero")
	assert.Equal(false, size.Parse("lol10", "k"), "Should be false")
	assert.Zero(size, "Should be zero")
	assert.Equal(true, size.Parse("13", "b"), "Should be false")
	assert.Equal(SizeBytes(13), size, "Should be equal to 13")
	assert.Equal(true, size.Parse("13", "k"), "Should be false")
	assert.Equal(SizeBytes(13312), size, "Should be equal to 13kb")
}
