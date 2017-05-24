package common

import (
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

func (st *Status) Parse(s string) {
	switch s {
	case "1":
		*st = FilterRemakes
		break
	case "2":
		*st = Trusted
		break
	case "3":
		*st = APlus
		break
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
		break
	case "2":
		*s = Date
		break
	case "3":
		*s = Downloads
		break
	case "4":
		*s = Size
		break
	case "5":
		*s = Seeders
		break
	case "6":
		*s = Leechers
		break
	case "7":
		*s = Completed
		break
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
		return "name"
	case Date:
		return "name"
	case Downloads:
		return "downloads"
	case Size:
		return "size"
	case Seeders:
		return "seeders"
	case Leechers:
		return "leechers"
	case Completed:
		return "completed"
	}
	return "date"
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
	return c.Main != 0 && c.Sub != 0
}

// Parse sets category by string
// returns true if string is valid otherwise returns false
func (c *Category) Parse(s string) (ok bool) {
	parts := strings.Split(s, "_")
	if len(parts) == 2 {
		tmp, err := strconv.ParseUint(parts[0], 10, 8)
		if err == nil {
			c.Main = uint8(tmp)
			tmp, err = strconv.ParseUint(parts[1], 10, 8)
			if err == nil {
				c.Sub = uint8(tmp)
				ok = true
			}
		}
	}
	return
}

// deprecated for TorrentParam
type SearchParam struct {
	Order    bool // True means acsending
	Status   Status
	Sort     SortMode
	Category Category
	Page     int
	UserID   uint
	Max      uint
	NotNull  string
	Query    string
}
