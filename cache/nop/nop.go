package nop

import (
	"github.com/ewhal/nyaa/common"
	"github.com/ewhal/nyaa/model"
)

type NopCache struct {
}

func (c *NopCache) Get(key common.SearchParam, fn func() ([]model.Torrent, int, error)) ([]model.Torrent, int, error) {
	return fn()
}

func (c *NopCache) ClearAll() {

}

// New creates a new Cache that does NOTHING :D
func New() *NopCache {
	return &NopCache{}
}
