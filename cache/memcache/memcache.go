package memcache

import (
	"github.com/ewhal/nyaa/common"
	"github.com/ewhal/nyaa/model"
)

type Memcache struct {
}

func (c *Memcache) Get(key common.SearchParam, r func() ([]model.Torrent, int, error)) (torrents []model.Torrent, num int, err error) {
	return
}

func (c *Memcache) ClearAll() {

}

func New() *Memcache {
	return &Memcache{}
}
