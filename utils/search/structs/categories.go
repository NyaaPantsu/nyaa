package structs

import (
	"strconv"
	"strings"

	"github.com/NyaaPantsu/nyaa/utils/categories"
)

// Categories multiple torrent categories
type Categories []*Category

// Category torrent categories
type Category struct {
	Main, Sub uint8
}

// String convert a category to a string
func (c Category) String() string {
	s := ""
	if c.Main != 0 {
		s += strconv.Itoa(int(c.Main))
	}
	s += "_"
	if c.Sub != 0 {
		s += strconv.Itoa(int(c.Sub))
	}
	return s
}

// IsSet check if a category is correctly set
func (c Category) IsSet() bool {
	return c.IsMainSet() && c.IsSubSet()
}

// IsMainSet check if a category is part of a main category
func (c Category) IsMainSet() bool {
	return c.Main != 0
}

// IsSubSet check if a category is a subset of a main category
func (c Category) IsSubSet() bool {
	return c.Sub != 0
}

// ParseCategories sets category by string
func ParseCategories(s string) []*Category {
	if s != "" {
		parts := strings.Split(s, ",")
		var cats []*Category
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
					if categories.Exists(partsCat[0] + "_" + partsCat[1]) {
						cats = append(cats, &Category{
							Main: c,
							Sub:  sub,
						})
					}
				}
			}
		}
		return cats
	}
	return Categories{}
}
