package metainfo

// this file is from https://github.com/majestrate/XD

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeebo/bencode"
)

// FilePath type
type FilePath []string

// FilePath : get filepath
func (f FilePath) FilePath() string {
	return filepath.Join(f...)
}

// Open : open file using base path
func (f FilePath) Open(base string) (*os.File, error) {
	return os.OpenFile(filepath.Join(base, f.FilePath()), os.O_RDWR|os.O_CREATE, 0600)
}

// FileInfo struct
type FileInfo struct {
	// length of file
	Length uint64 `bencode:"length"`
	// relative path of file
	Path FilePath `bencode:"path"`
	// md5sum
	Sum []byte `bencode:"md5sum,omitempty"`
}

// Info : info section of torrent file
type Info struct {
	// length of pices in bytes
	PieceLength uint32 `bencode:"piece length"`
	// piece data
	Pieces []byte `bencode:"pieces"`
	// name of root file
	Path string `bencode:"name"`
	// file metadata
	Files []FileInfo `bencode:"files,omitempty"`
	// private torrent
	Private *int64 `bencode:"private,omitempty"`
	// length of file in signle file mode
	Length uint64 `bencode:"length,omitempty"`
	// md5sum
	Sum []byte `bencode:"md5sum,omitempty"`
}

// GetFiles : get fileinfos from this info section
func (i Info) GetFiles() (infos []FileInfo) {
	if i.Length > 0 {
		infos = append(infos, FileInfo{
			Length: i.Length,
			Path:   FilePath([]string{i.Path}),
			Sum:    i.Sum,
		})
	} else {
		infos = append(infos, i.Files...)
	}
	return
}

// NumPieces : length of Info
func (i Info) NumPieces() uint32 {
	return uint32(len(i.Pieces) / 20)
}

// TorrentFile : a torrent file
type TorrentFile struct {
	Info         Info       `bencode:"info"`
	Announce     string     `bencode:"announce"`
	AnnounceList [][]string `bencode:"announce-list"`
	Created      uint64     `bencode:"created"`
	Comment      []byte     `bencode:"comment"`
	CreatedBy    []byte     `bencode:"created by"`
	Encoding     []byte     `bencode:"encoding"`
}

// TotalSize : get total size of files from torrent info section
func (tf *TorrentFile) TotalSize() uint64 {
	if tf.IsSingleFile() {
		return tf.Info.Length
	}
	total := uint64(0)
	for _, f := range tf.Info.Files {
		total += f.Length
	}
	return total
}

// GetAllAnnounceURLS : get all trackers url
func (tf *TorrentFile) GetAllAnnounceURLS() (l []string) {
	l = make([]string, 0, 64)
	if len(tf.Announce) > 0 {
		l = append(l, tf.Announce)
	}
	for _, al := range tf.AnnounceList {
		for _, a := range al {
			if len(a) > 0 {
				l = append(l, a)
			}
		}
	}
	return
}

// TorrentName : return torrent name
func (tf *TorrentFile) TorrentName() string {
	return tf.Info.Path
}

// IsPrivate : return true if this torrent is private otherwise return false
func (tf *TorrentFile) IsPrivate() bool {
	return tf.Info.Private != nil && *tf.Info.Private == 1
}

// IsSingleFile : return true if this torrent is for a single file
func (tf *TorrentFile) IsSingleFile() bool {
	return tf.Info.Length > 0
}

// Encode : bencode this file via an io.Writer
func (tf *TorrentFile) Encode(w io.Writer) (err error) {
	enc := bencode.NewEncoder(w)
	err = enc.Encode(tf)
	return
}

// Decode : load from an io.Reader
func (tf *TorrentFile) Decode(r io.Reader) (err error) {
	dec := bencode.NewDecoder(r)
	err = dec.Decode(tf)
	return
}

type torrentRaw struct {
	InfoRaw bencode.RawMessage `bencode:"info"`
}

// DecodeInfohash : Decode and calculate the info hash
func DecodeInfohash(r io.Reader) (hash string, err error) {
	var t torrentRaw
	d := bencode.NewDecoder(r)
	err = d.Decode(&t)
	if err != nil {
		return
	}

	s := sha1.New()
	_, err = s.Write(t.InfoRaw)
	if err != nil {
		return
	}
	rawHash := s.Sum(nil)

	hash = strings.ToUpper(hex.EncodeToString(rawHash))
	return
}
