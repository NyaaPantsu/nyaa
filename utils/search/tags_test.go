package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTags_Parse(t *testing.T) {
	assert := assert.New(t)
	tags := Tags{}
	ok := tags.Parse("")
	assert.Empty(tags)
	assert.False(ok)
	tags = Tags{}
	ok = tags.Parse("dd")
	assert.Equal(Tags{"dd"}, tags)
	assert.True(ok)
	tags = Tags{}
	ok = tags.Parse("dd,,,,")
	assert.Equal(Tags{"dd"}, tags, "Should remove empty tags")
	assert.True(ok)
	tags = Tags{}
	ok = tags.Parse("dd,,ff,,ss")
	assert.Equal(Tags{"dd", "ff", "ss"}, tags, "Should remove empty tags and keep the other")
	assert.True(ok)
}
func TestTags_ToDBQuery(t *testing.T) {
	tags := Tags{"dd", "ff"}
	assert := assert.New(t)
	search, args := tags.ToDBQuery()
	assert.Equal("tags = ? AND tags = ?", search, "Should be equal")
	assert.Equal([]string{"dd", "ff"}, args, "Should be equal")
	tags = Tags{"dd"}
	search, args = tags.ToDBQuery()
	assert.Equal("tags = ?", search, "Should be equal to 'tags = ?'")
	assert.Equal([]string{"dd"}, args, "Should be equal to '3_13'")
	tags = Tags{}
	search, args = tags.ToDBQuery()
	assert.Empty(search, "Should be empty")
	assert.Empty(args, "Should be empty")
}

func TestTags_ToESQuery(t *testing.T) {
	tags := Tags{"dd", "ff"}
	assert := assert.New(t)
	assert.Equal("tags:dd tags:ff", tags.ToESQuery(), "Should be equal")

	tags = Tags{"dd"}
	assert.Equal("tags:dd", tags.ToESQuery(), "Should be equal")

	tags = Tags{}
	assert.Empty(tags.ToESQuery(), "Should be empty")
}
