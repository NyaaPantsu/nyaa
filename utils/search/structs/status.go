package structs

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

// String convert a status to a string
func (st *Status) String() string {
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

// Parse a string to a status
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
