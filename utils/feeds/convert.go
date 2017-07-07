package nyaafeeds

import (
	"strconv"
	"strings"

	"github.com/NyaaPantsu/nyaa/utils/categories"
)

// ConvertToCat : Convert a torznab cat to our cat
func ConvertToCat(cat string) string {
	if cat == "" {
		return ""
	}

	cats := strings.Split(cat, ",")
	var returnCat []string
	for _, val := range cats {
		localeCat := convertCat(val)
		if localeCat != "" {
			returnCat = append(returnCat, localeCat)
		}
	}
	return strings.Join(returnCat, ",")
}

func convertCat(cat string) string {
	if cat == "5070" {
		return "3_5"
	}

	if len(cat) < 6 {
		return ""
	}

	cI, _ := strconv.Atoi(cat[2:4])
	subI, _ := strconv.Atoi(cat[4:6])

	c := strconv.Itoa(cI)
	sub := strconv.Itoa(subI)

	if categories.Exists(c + "_" + sub) {
		return c + "_" + sub
	}

	return ""
}

// ConvertFromCat : Convert a cat to a torznab cat
func ConvertFromCat(category string) (cat string) {
	c := strings.Split(category, "_")
	if len(c[0]) < 2 {
		c[0] = "0" + c[0]
	}
	if len(c) < 2 || c[1] == "" {
		cat = "10" + c[0] + "00"
		return
	}
	if len(c[1]) < 2 {
		c[1] = "0" + c[1]
	}
	cat = "10" + c[0] + c[1]
	return
}
