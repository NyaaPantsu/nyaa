package model

import (
	"github.com/NyaaPantsu/nyaa/config"
	"github.com/zeebo/bencode"
)

type File struct {
	ID           uint   `gorm:"column:file_id;primary_key"`
	TorrentID    uint   `gorm:"column:torrent_id;unique_index:idx_tid_path"`
	// this path is bencode'd, call Path() to obtain
	BencodedPath string `gorm:"column:path;unique_index:idx_tid_path"`
	Filesize     int64  `gorm:"column:filesize"`
}

func (f File) TableName() string {
	return config.FilesTableName
}

// Returns the total size of memory allocated for this struct
func (f File) Size() int {
	return (2 + len(f.BencodedPath) + 1) * 8;
}

func (f *File) Path() (out []string) {
	bencode.DecodeString(f.BencodedPath, &out)
	return
}

func (f *File) SetPath(path []string) error {
	encoded, err := bencode.EncodeString(path)
	if err != nil {
		return err
	}

	f.BencodedPath = encoded
	return nil
}

