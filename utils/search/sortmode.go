package search

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
	MaxIota
)

// SortField is a database column in ES/DB
type SortField struct {
	ES string
	DB string
}

var sortFields = []SortField{
	{"id", config.Get().Models.TorrentsTableName + ".torrent_id"},
	{"name.raw", "torrent_name"},
	{"date", "date"},
	{"downloads", "downloads"},
	{"filesize", "filesize"},
	{"seeders", "seeders"},
	{"leechers", "leechers"},
	{"completed", "completed"},
}

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
	return s.toField().ES
}

// ToDBField convert a sortmode to use with database
func (s *SortMode) ToDBField() string {
	return s.toField().DB
}

// Private function to convert sormode to a field struct
func (s *SortMode) toField() SortField {
	// if sortmode is within range
	if *s >= MaxIota || *s < 1 {
		s.Parse(config.Get().Torrents.Order)
	}
	return sortFields[int(*s)]
}
