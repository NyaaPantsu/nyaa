package common

import (
	humanize "github.com/dustin/go-humanize"

	"fmt"
	"strconv"
	"strings"
)

type Status uint8

const (
	ShowAll Status = iota
	FilterRemakes
	Trusted
	APlus
)

func (st *Status) ToString() string {
	switch *st {
	case FilterRemakes:
		return "1"
	case Trusted:
		return "2"
	case APlus:
		return "3"
	}
	return ""
}

func (st *Status) Parse(s string) {
	switch s {
	case "1":
		*st = FilterRemakes
	case "2":
		*st = Trusted
	case "3":
		*st = APlus
	default:
		*st = ShowAll
	}
}

type SortMode uint8

const (
	ID SortMode = iota
	Name
	Date
	Downloads
	Size
	Seeders
	Leechers
	Completed
)

func (s *SortMode) Parse(str string) {
	switch str {
	case "1":
		*s = Name
	case "2":
		*s = Date
	case "3":
		*s = Downloads
	case "4":
		*s = Size
	case "5":
		*s = Seeders
	case "6":
		*s = Leechers
	case "7":
		*s = Completed
	default:
		*s = ID
	}
}

/* INFO Always need to keep in sync with the field that are used in the
 * elasticsearch index.
 * TODO Verify the field in postgres database
 */
func (s *SortMode) ToESField() string {
	switch *s {
	case ID:
		return "id"
	case Name:
		return "name.raw"
	case Date:
		return "date"
	case Downloads:
		return "downloads"
	case Size:
		return "filesize"
	case Seeders:
		return "seeders"
	case Leechers:
		return "leechers"
	case Completed:
		return "completed"
	}
	return "id"
}

func (s *SortMode) ToDBField() string {
	switch *s {
	case ID:
		return "torrent_id"
	case Name:
		return "torrent_name"
	case Date:
		return "date"
	case Downloads:
		return "downloads"
	case Size:
		return "filesize"
	case Seeders:
		return "seeders"
	case Leechers:
		return "leechers"
	case Completed:
		return "completed"
	}
	return "id"
}

type Category struct {
	Main, Sub uint8
}

func (c Category) String() (s string) {
	if c.Main != 0 {
		s += strconv.Itoa(int(c.Main))
	}
	s += "_"
	if c.Sub != 0 {
		s += strconv.Itoa(int(c.Sub))
	}
	return
}

func (c Category) IsSet() bool {
	return c.IsMainSet() && c.IsSubSet()
}

func (c Category) IsMainSet() bool {
	return c.Main != 0
}

func (c Category) IsSubSet() bool {
	return c.Sub != 0
}

// Parse sets category by string
// returns true if string is valid otherwise returns false
func ParseCategories(s string) []*Category {
	fmt.Println("s: " + s)
	if s != "" {
		parts := strings.Split(s, ",")
		categories := make([]*Category, len(parts))
		for key, val := range parts {
			partsCat := strings.Split(val, "_")

			if len(partsCat) == 2 {
				tmp, err := strconv.ParseUint(partsCat[0], 10, 8)
				if err == nil {
					c := uint8(tmp)
					tmp, err = strconv.ParseUint(partsCat[1], 10, 8)
					var sub uint8
					if err == nil {
						sub = uint8(tmp)
					}
					categories[key] = &Category{
						Main: c,
						Sub:  sub,
					}
				}
			}
		}
		return categories
	}
	return Categories{}
}

type SizeBytes uint64

func (sz *SizeBytes) Parse(s string) bool {
	size64, err := humanize.ParseBytes(s)
	if err != nil {
		*sz = 0
		return false
	}
	*sz = SizeBytes(size64)
	return true
}

type Categories []*Category

// deprecated for TorrentParam
type SearchParam struct {
	TorrentID uint
	FromID    uint // Search for torrentID > FromID
	Order     bool // True means acsending
	Status    Status
	Sort      SortMode
	Category  Categories
	FromDate  string
	ToDate    string
	Page      int
	UserID    uint
	Max       uint
	NotNull   string
	Language  string
	MinSize   SizeBytes
	MaxSize   SizeBytes
	Query     string
}
