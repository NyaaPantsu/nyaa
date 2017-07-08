package structs

import "github.com/NyaaPantsu/nyaa/models"

type TorrentCache struct {
	Torrents []models.Torrent
	Count    int
}
