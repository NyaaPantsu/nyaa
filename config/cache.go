package config

// CacheConfig is config struct for caching strategy
type CacheConfig struct {
	Dialect string
	URL     string
	Size    float64
}

const DefaultCacheSize = 1 << 10

var DefaultCacheConfig = CacheConfig{
	Dialect: "nop",
}
