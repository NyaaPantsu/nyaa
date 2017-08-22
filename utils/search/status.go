package search

const (
	// ShowAll by default show all torrents
	ShowAll Status = iota
	// Normal torrent
	Normal
	// FilterRemakes filter torrent remakes
	FilterRemakes
	// Trusted trusted torrents
	Trusted
	// APlus torrents not used anymore
	APlus
)

// Status torrent status
type Status uint8

// String convert a status to a string
func (st *Status) String() string {
	switch *st {
	case Normal:
		return "1"
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
	case "1":
		*st = Normal
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

// ToESQuery prepare an ES statement for status
func (st *Status) ToESQuery() string {
	if *st != ShowAll {
		if *st != FilterRemakes {
			// Only show torrents with status over the one specified
			return "status:>" + st.String()
		}
		/* From the old nyaa behavior, FilterRemake means everything BUT
		* remakes
		 */
		return "!status:" + st.String()
	}
	return ""
}

// ToDBQuery prepare a DB statement for status
func (st *Status) ToDBQuery() string {
	if *st != ShowAll {
		if *st != FilterRemakes {
			// Only show torrents with status over the one specified
			return "status >= " + st.String()
		}
		/* From the old nyaa behavior, FilterRemake means everything BUT
		* remakes
		 */
		return "status <> " + st.String()
	}
	return ""
}
