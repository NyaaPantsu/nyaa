package nyaafeeds

import (
	"strings"
)

// ConvertToCat : Convert a torznab cat to our cat
func ConvertToCat(cat string) (string, string) {
	c := strings.Split(cat, "0000")
	if len(c) < 2 {
		return cat, ""
	}
	return c[0], c[1]
}

// ConvertFromCat : Convert a cat to a torznab cat
func ConvertFromCat(category string) (cat string) {
	c := strings.Split(category, "_")
	if len(c) < 2 {
		cat = c[0] + "0000"
		return
	}
	cat = c[0] + "0000" + c[1]
	return
}
