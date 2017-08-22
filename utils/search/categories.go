package search

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
	var cats Categories
	if s != "" {
		parts := strings.Split(s, ",")
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
	}
	return cats
}

// ToESQuery prepare the ES statement for categories
func (ct *Categories) ToESQuery() string {
	if len(*ct) > 0 {
		conditionsOr := make([]string, len(*ct))
		for key, val := range *ct {
			if val.IsSubSet() {
				conditionsOr[key] = "category: " + strconv.FormatInt(int64(val.Main), 10) + " AND sub_category: " + strconv.FormatInt(int64(val.Sub), 10)
			} else {
				if val.Main > 0 {
					conditionsOr[key] = "category: " + strconv.FormatInt(int64(val.Main), 10)
				}
			}
			if key == 0 && conditionsOr[key] != "" {
				conditionsOr[key] = "(" + conditionsOr[key]
			}
			if len(*ct) > 1 {
				conditionsOr[key] = "(" + conditionsOr[key] + ")"
			}
			if key == len(*ct)-1 && conditionsOr[key] != "" {
				conditionsOr[key] = conditionsOr[key] + ")"
			}
		}

		return strings.Join(conditionsOr, " OR ")
	}
	return ""
}

// ToDBQuery prepare the DB statement for categories
func (ct *Categories) ToDBQuery() (string, []interface{}) {
	var args []interface{}
	if len(*ct) > 0 {
		conditionsOr := make([]string, len(*ct))
		for key, val := range *ct {
			if val.IsMainSet() {
				conditionsOr[key] += "category = ?"
				args = append(args, val.Main)
				if val.IsSubSet() {
					conditionsOr[key] += " AND sub_category = ?"
					args = append(args, val.Sub)
				}
				if len(*ct) > 1 {
					conditionsOr[key] = "(" + conditionsOr[key] + ")"
				}
			}
		}
		return strings.Join(conditionsOr, " OR "), args
	}
	return "", args
}
