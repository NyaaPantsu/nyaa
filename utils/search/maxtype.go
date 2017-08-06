package search

import (
	"strconv"

	"github.com/NyaaPantsu/nyaa/config"
)

type maxType uint32

// Parse the maximum number of result
func (m *maxType) Parse(s string) {
	max, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		// If we can't convert the limit from url to int, we set it to default value
		max = uint64(config.Get().Navigation.TorrentsPerPage)
	} else if max > uint64(config.Get().Navigation.MaxTorrentsPerPage) {
		// If the maximum value is greater than the maximum set in config.yml, we overwrite it with the configured one
		// Stops someone to make an unthinkable query of max=huge number
		max = uint64(config.Get().Navigation.MaxTorrentsPerPage)
	}
	*m = maxType(max)
}
