package common

import "strconv"

type Status uint8

const (
	ShowAll Status = iota
	FilterRemakes
	Trusted
	APlus
)

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
