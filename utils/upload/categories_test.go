package upload

import (
	"testing"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/stretchr/testify/assert"
)

func TestCategory(t *testing.T) {
	assert := assert.New(t)
	dummyTorrent := &models.Torrent{Category: 1, SubCategory: 1}
	tests := []struct {
		torrent  *models.Torrent
		platform int
		sukebei  bool
		expected string
	}{
		{dummyTorrent, anidex, false, ""},
		{dummyTorrent, ttosho, false, "5"},
		{dummyTorrent, ttosho, true, "12"},
		{dummyTorrent, nyaasi, false, "6_1"},
		{dummyTorrent, nyaasi, true, "1_1"},
		{dummyTorrent, 20, true, ""},
		{&models.Torrent{Category: 33, SubCategory: 33}, nyaasi, true, ""},
	}

	for _, test := range tests {
		if test.sukebei {
			// workaround to make the function believe we are in sukebei
			config.Get().Models.TorrentsTableName = "sukebei_torrents"
		} else {
			config.Get().Models.TorrentsTableName = "torrents"
		}
		assert.Equal(test.expected, Category(test.platform, test.torrent))
	}
}
