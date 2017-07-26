package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// C global cache variable
var C *cache.Cache

func init() {
	C = cache.New(5*time.Minute, 10*time.Minute)
}
