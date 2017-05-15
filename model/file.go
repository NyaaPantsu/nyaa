package model

type File struct {
	ID        uint   `gorm:"column:file_id;primary_key"`
	TorrentID uint   `gorm:"column:torrent_id;unique_index:idx_tid_path"`
	Path      string `gorm:"column:path;unique_index:idx_tid_path"`
	Filesize  int64  `gorm:"column:filesize"`
}

// Returns the total size of memory allocated for this struct
func (f File) Size() int {
	return (2 + len(f.Path) + 2) * 8;
}

