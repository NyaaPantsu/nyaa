package search

import (
	"strings"

	"github.com/NyaaPantsu/nyaa/config"
)

// Tags struct for search
type Tags config.ArrayString

// Parse a tag list separated by commas to Tags struct
func (ta *Tags) Parse(str string) bool {
	if str == "" {
		return false
	}
	arr := strings.Split(str, ",")
	for _, tag := range arr {
		if tag != "" {
			*ta = append(*ta, tag)
		}
	}
	return len(*ta) > 0
}

// ToESQuery prepare the ES statement for tags
func (ta Tags) ToESQuery() string {
	if len(ta) > 0 && ta[0] != "" {
		return "tags:" + strings.Join(ta, " tags:")
	}
	return ""
}

// ToDBQuery prepare the DB statement for tags
func (ta Tags) ToDBQuery() (string, []string) {
	if len(ta) > 0 {
		conditionsAnd := make([]string, len(ta))
		for key := range ta {
			conditionsAnd[key] = "tags = ?"
		}
		return strings.Join(conditionsAnd, " AND "), ta
	}
	return "", []string{}
}
