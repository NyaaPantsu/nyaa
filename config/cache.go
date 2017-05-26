package config

// CacheConfig is config struct for caching strategy
type CacheConfig struct {
	Dialect string
	URL     string
	Size    float64
}

// DefaultCacheSize : Size by default for the cache
const DefaultCacheSize = 1 << 10

// DefaultCacheConfig : Config by default for the cache
var DefaultCacheConfig = CacheConfig{
	Dialect: "nop",
}
