package structs

import "github.com/NyaaPantsu/nyaa/models"

// TorrentCache torrent cache struct
type TorrentCache struct {
	Torrents []models.Torrent
	Count    int
}
