package structs

import (
	"math"
	"time"

	humanize "github.com/dustin/go-humanize"

	"strconv"
	"strings"

	"fmt"

	"github.com/NyaaPantsu/nyaa/config"
	catUtil "github.com/NyaaPantsu/nyaa/utils/categories"
	"github.com/NyaaPantsu/nyaa/utils/publicSettings"
)

const (
	// ShowAll by default show all torrents
	ShowAll Status = 0
	// FilterRemakes filter torrent remakes
	FilterRemakes = 2
	// Trusted trusted torrents
	Trusted = 3
	// APlus torrents not used anymore
	APlus = 4
)

// Status torrent status
type Status uint8

// SortMode selected sort mode
type SortMode uint8

// Category torrent categories
type Category struct {
	Main, Sub uint8
}

// SizeBytes size in bytes
type SizeBytes uint64

// DateFilter date to filter for
type DateFilter string

// Categories multiple torrent categories
type Categories []*Category

// TorrentParam defines all parameters that can be provided when searching for a torrent
type TorrentParam struct {
	Full      bool // True means load all members
	Order     bool // True means ascending
	Hidden    bool // True means filter hidden torrents
	Deleted   bool // False means filter deleted torrents
	Status    Status
	Sort      SortMode
	Category  Categories
	Max       uint32
	Offset    uint32
	UserID    uint32
	TorrentID uint32
	FromID    uint32
	FromDate  DateFilter
	ToDate    DateFilter
	NotNull   string // csv
	NameLike  string // csv
	Languages publicSettings.Languages
	MinSize   SizeBytes
	MaxSize   SizeBytes
}

func (p *TorrentParam) Identifier() string {
	cats := ""
	for _, v := range p.Category {
		cats += fmt.Sprintf("%d%d", v.Main, v.Sub)
	}
	languages := ""
	for _, v := range p.Languages {
		languages += fmt.Sprintf("%s%s", v.Code, v.Name)
	}
	return fmt.Sprintf("%s%s%s%d%d%d%d%d%d%d%s%s%d%d%s%t%t%t%t", p.NameLike, p.NotNull, languages, p.Max, p.Offset, p.FromID, p.MinSize, p.MaxSize, p.Status, p.Sort, p.FromDate, p.ToDate, p.UserID, p.TorrentID, cats, p.Full, p.Order, p.Hidden, p.Deleted)
}

func (st *Status) ToString() string {
	switch *st {
	case FilterRemakes:
		return "2"
	case Trusted:
		return "3"
	case APlus:
		return "4"
	}
	return ""
}

func (st *Status) Parse(s string) {
	switch s {
	case "2":
		*st = FilterRemakes
	case "3":
		*st = Trusted
	case "4":
		*st = APlus
	default:
		*st = ShowAll
	}
}

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
		*s = Date
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
	return config.Get().Torrents.Order
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

// ParseCategories sets category by string
func ParseCategories(s string) []*Category {
	if s != "" {
		parts := strings.Split(s, ",")
		var categories []*Category
		for _, val := range parts {
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
					if catUtil.Exists(partsCat[0] + "_" + partsCat[1]) {
						categories = append(categories, &Category{
							Main: c,
							Sub:  sub,
						})
					}
				}
			}
		}
		return categories
	}
	return Categories{}
}

// ParseLanguages sets languages by string
func ParseLanguages(s string) publicSettings.Languages {
	var languages publicSettings.Languages
	if s != "" {
		parts := strings.Split(s, ",")
		for _, lang := range parts {
			languages = append(languages, publicSettings.Language{Name: "", Code: lang}) // We just need the code
		}
	}
	return languages
}

func (sz *SizeBytes) Parse(s string, sizeType string) bool {
	if s == "" {
		*sz = 0
		return false
	}
	var multiplier uint64
	switch sizeType {
	case "b":
		multiplier = 1
	case "k":
		multiplier = uint64(math.Exp2(10))
	case "m":
		multiplier = uint64(math.Exp2(20))
	case "g":
		multiplier = uint64(math.Exp2(30))
	}
	size64, err := humanize.ParseBytes(s)
	if err != nil {
		*sz = 0
		return false
	}
	*sz = SizeBytes(size64 * multiplier)
	return true
}

func (d *DateFilter) Parse(s string, dateType string) bool {
	if s == "" {
		*d = ""
		return false
	}
	dateInt, err := strconv.Atoi(s)
	if err != nil {
		*d = ""
		return false
	}
	switch dateType {
	case "m":
		*d = DateFilter(time.Now().AddDate(0, -dateInt, 0).Format("2006-01-02"))
	case "y":
		*d = DateFilter(time.Now().AddDate(-dateInt, 0, 0).Format("2006-01-02"))
	default:
		*d = DateFilter(time.Now().AddDate(0, 0, -dateInt).Format("2006-01-02"))
	}
	return true
}
