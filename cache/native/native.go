package native

import (
	"container/list"
	"sync"
	"time"

	"github.com/NyaaPantsu/nyaa/common"
	"github.com/NyaaPantsu/nyaa/model"
)

const expiryTime = time.Minute

// NativeCache implements cache.Cache
type NativeCache struct {
	cache     map[common.SearchParam]*list.Element
	ll        *list.List
	totalUsed int
	mu        sync.Mutex

	// Size sets the maximum size of the cache before evicting unread data in MB
	Size float64
}

// New Creates New Native Cache instance
func New(sz float64) *NativeCache {
	return &NativeCache{
		cache: make(map[common.SearchParam]*list.Element, 10),
		Size:  sz,
		ll:    list.New(),
	}
}

// Key stores the ID of either a thread or board page
type Key struct {
	LastN uint8
	Board string
	ID    uint64
}

// Single cache entry
type store struct {
	sync.Mutex  // Controls general access to the contents of the struct
	lastFetched time.Time
	key         common.SearchParam
	data        []model.Torrent
	count, size int
	n           *NativeCache
}

// Check the cache for and existing record. If miss, run fn to retrieve fresh
// values.
func (n *NativeCache) Get(key common.SearchParam, fn func() ([]model.Torrent, int, error)) (
	data []model.Torrent, count int, err error,
) {
	s := n.getStore(key)

	// Also keeps multiple requesters from simultaneously requesting the same
	// data
	s.Lock()
	defer s.Unlock()

	if s.isFresh() {
		return s.data, s.count, nil
	}

	data, count, err = fn()
	if err != nil {
		return
	}
	s.update(data, count)
	return
}

// Retrieve a store from the cache or create a new one
func (n *NativeCache) getStore(k common.SearchParam) (s *store) {
	n.mu.Lock()
	defer n.mu.Unlock()

	el := n.cache[k]
	if el == nil {
		s = &store{key: k, n: n}
		n.cache[k] = n.ll.PushFront(s)
	} else {
		n.ll.MoveToFront(el)
		s = el.Value.(*store)
	}
	return s
}

// Clear the cache. Only used for testing.
func (n *NativeCache) ClearAll() {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.ll = list.New()
	n.cache = make(map[common.SearchParam]*list.Element, 10)
}

// Update the total used memory counter and evict, if over limit
func (n *NativeCache) updateUsedSize(delta int) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.totalUsed += delta

	for n.totalUsed > int(n.Size)<<20 {
		e := n.ll.Back()
		if e == nil {
			break
		}
		s := n.ll.Remove(e).(*store)
		delete(n.cache, s.key)
		s.Lock()
		n.totalUsed -= s.size
		s.Unlock()
	}
}

// Return, if the data can still be considered fresh, without querying the DB
func (s *store) isFresh() bool {
	if s.lastFetched.IsZero() { // New store
		return false
	}
	return s.lastFetched.Add(expiryTime).After(time.Now())
}

// Stores the new values of s. Calculates and stores the new size. Passes the
// delta to the central cache to fire eviction checks.
func (s *store) update(data []model.Torrent, count int) {
	newSize := 0
	for _, d := range data {
		newSize += d.Size()
	}
	s.data = data
	s.count = count
	delta := newSize - s.size
	s.size = newSize
	s.lastFetched = time.Now()

	// In a separate goroutine, to ensure there is never any lock intersection
	go s.n.updateUsedSize(delta)
}
