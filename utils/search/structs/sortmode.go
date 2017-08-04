package structs

import "github.com/NyaaPantsu/nyaa/config"

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

// SortMode selected sort mode
type SortMode uint8

// Parse a string to sortMode
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

// ToESField convert a sortmode to use with ES
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

// ToDBField convert a sortmode to use with database
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
