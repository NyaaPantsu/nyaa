package cache

import (
	"github.com/NyaaPantsu/nyaa/cache/memcache"
	"github.com/NyaaPantsu/nyaa/cache/native"
	"github.com/NyaaPantsu/nyaa/cache/nop"
	"github.com/NyaaPantsu/nyaa/common"
	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/model"
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
