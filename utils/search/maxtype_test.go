package search

import (
	"testing"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/stretchr/testify/assert"
)

func TestMaxType_Parse(t *testing.T) {
	var max maxType
	defMax := maxType(config.Get().Navigation.TorrentsPerPage)
	max.Parse("")
	assert := assert.New(t)
	assert.Equal(defMax, max, "Should be equal")

	max.Parse("100")
	assert.Equal(maxType(100), max, "Should be equal to 100")

	max.Parse("lol10")
	assert.Equal(defMax, max, "Should be equal")
}
