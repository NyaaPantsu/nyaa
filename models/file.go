package models

import (
	"strings"
	
	"github.com/NyaaPantsu/nyaa/config"
	"github.com/zeebo/bencode"
)

// File model
type File struct {
	ID        uint `gorm:"column:file_id;primary_key"`
	TorrentID uint `gorm:"column:torrent_id;unique_index:idx_tid_path"`
	// this path is bencode'd, call Path() to obtain
	BencodedPath string `gorm:"column:path;unique_index:idx_tid_path"`
	Filesize     int64  `gorm:"column:filesize"`
}

// FileJSON for file model in json
type FileJSON struct {
	Path     string `json:"path"`
	Filesize int64  `json:"filesize"`
}

// TableName : Return the name of files table
func (f File) TableName() string {
	return config.Get().Models.FilesTableName
}

// Size : Returns the total size of memory allocated for this struct
func (f File) Size() int {
	return (2 + len(f.BencodedPath) + 1) * 8
}

// Path : Returns the path to the file
func (f *File) Path() (out []string) {
	bencode.DecodeString(f.BencodedPath, &out)
	return
}

// SetPath : Set the path of the file
func (f *File) SetPath(path []string) error {
	encoded, err := bencode.EncodeString(path)
	if err != nil {
		return err
	}

	f.BencodedPath = encoded
	return nil
}

// Filename : Returns the filename of the file
func (f *File) Filename() string {
	path := f.Path()
	return path[len(path)-1]
}

// FilenameWithoutExtension : Returns the filename of the file without the extension
func (f *File) FilenameWithoutExtension() string {
	path := f.Path()
	fileName := path[len(path)-1]
	index := strings.LastIndex(fileName, ".")
	
	if index == -1 {
		return fileName
	}
	
	return fileName[:index]
}

// FilenameExtension : Returns the extension of a filename, or an empty string
func (f *File) FilenameExtension() string {
	path := f.Path()
	fileName := path[len(path)-1]
	index := strings.LastIndex(fileName, ".")
	
	if index == -1 {
		return ""
	}
	
	return fileName[index:]
}
