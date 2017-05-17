package native

import (
	"sync"
	"testing"

	"github.com/NyaaPantsu/nyaa/common"
	"github.com/NyaaPantsu/nyaa/model"
)

// Basic test for deadlocks and race conditions
func TestConcurrency(t *testing.T) {
	c := New(0.000001)

	fn := func() ([]model.Torrent, int, error) {
		return []model.Torrent{{}, {}, {}}, 10, nil
	}

	var wg sync.WaitGroup
	wg.Add(300)
	for i := 0; i < 3; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				go func(j int) {
					defer wg.Done()
					k := common.SearchParam{
						Page: j,
					}
					if _, _, err := c.Get(k, fn); err != nil {
						t.Fatal(err)
					}
				}(j)
			}
		}()
	}
	wg.Wait()
}
