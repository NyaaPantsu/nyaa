package cache

import (
	"github.com/ewhal/nyaa/cache/memcache"
	"github.com/ewhal/nyaa/cache/native"
	"github.com/ewhal/nyaa/cache/nop"
	"github.com/ewhal/nyaa/common"
	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/model"
)

// Cache defines interface for caching search results
type Cache interface {
	Get(key common.SearchParam, r func() ([]model.Torrent, int, error)) ([]model.Torrent, int, error)
	ClearAll()
}

// Impl cache implementation instance
var Impl Cache

func Configure(conf *config.CacheConfig) (err error) {
	switch conf.Dialect {
	case "native":
		Impl = native.New(conf.Size)
		return
	case "memcache":
		Impl = memcache.New()
		return
	default:
		Impl = nop.New()
	}
	return
}
