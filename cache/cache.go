package cache

import (
	"container/list"
	"sync"
	"time"

	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/util/search"
)

var (
	cache     = make(map[search.SearchParam]*list.Element, 10)
	ll        = list.New()
	totalUsed int
	mu        sync.Mutex

	// Mutable for quicker testing
	expiryTime = time.Second * 60

	// Size sets the maximum size of the cache before evicting unread data in MB
	Size float64 = 1 << 10
)

// Key stores the ID of either a thread or board page
type Key struct {
	LastN uint8
	Board string
	ID    uint64
}

// Single cache entry
type store struct {
	// Controls general access to the contents of the struct, except for size
	sync.Mutex
	lastFetched time.Time
	key         search.SearchParam
	data        []model.Torrents
	size        int
}

// Check the cache for and existing record. If miss, run fn to retrieve fresh
// values.
func CheckCache(key search.SearchParam, fn func() ([]model.Torrents, error)) (
	[]model.Torrents, error,
) {
	s := getStore(key)

	// Also keeps multiple requesters from simultaneously requesting the same
	// data
	s.Lock()
	defer s.Unlock()

	if s.isFresh() {
		return s.data, nil
	}

	data, err := fn()
	if err != nil {
		return nil, err
	}
	s.lastFetched = time.Now()
	return data, nil
}

// Retrieve a store from the cache or create a new one
func getStore(k search.SearchParam) (s *store) {
	mu.Lock()
	defer mu.Unlock()

	el := cache[k]
	if el == nil {
		s = &store{key: k}
		cache[k] = ll.PushFront(s)
	} else {
		ll.MoveToFront(el)
		s = el.Value.(*store)
	}
	return s
}

// Clear the cache. Only used for testing.
func Clear() {
	mu.Lock()
	defer mu.Unlock()

	ll = list.New()
	cache = make(map[search.SearchParam]*list.Element, 10)
}

// Update the total used memory counter and evict, if over limit
func updateUsedSize(delta int) {
	mu.Lock()
	defer mu.Unlock()

	totalUsed += delta

	for totalUsed > int(Size)*(1<<20) {
		s := ll.Remove(ll.Back()).(*store)
		delete(cache, s.key)
		totalUsed -= s.size
	}
}

// Return, if the data can still be considered fresh, without querying the DB
func (s *store) isFresh() bool {
	return s.lastFetched.Add(expiryTime).Before(time.Now())
}

// Stores the new values of s. Calculates and stores the new size. Passes the
// delta to the central cache to fire eviction checks.
func (s *store) update(data []model.Torrents) {
	newSize := 0
	for _, d := range data {
		newSize += d.Size()
	}
	s.data = data
	delta := newSize - s.size
	s.size = newSize

	// Technically it is possible to update the size even when the store is
	// already evicted, but that should never happen, unless you have a very
	// small cache, very large stored datasets and a lot of traffic.
	updateUsedSize(delta)
}
