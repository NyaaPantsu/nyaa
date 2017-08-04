package structs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentifier(t *testing.T) {
	torrentParam := &TorrentParam{}
	assert := assert.New(t)
	assert.Empty(torrentParam.Identifier(), "It should be empty")
}
