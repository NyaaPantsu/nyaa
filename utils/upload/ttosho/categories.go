package ttoshoConfig

import (
	"fmt"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
)

// Convert automatically our sukebei cats to anidex Hentai cats
var sukebeiCategories = map[string]string{
	"1_1": "12",
	"1_2": "12",
	"1_3": "14",
	"1_4": "13",
	"1_5": "4",
	"2_1": "4",
	"2_2": "15",
}

var normalCategories = map[string]string{
	"3_12": "1",
	"3_5":  "1",
	"3_13": "10",
	"3_6":  "7",
	"2_3":  "2",
	"2_4":  "2",
	"4_7":  "3",
	"4_8":  "7",
	"4_14": "10",
	"5_9":  "8",
	"5_10": "8",
	"5_18": "10",
	"5_11": "7",
	"6_15": "5",
	"6_16": "5",
	"1_1":  "5",
	"1_2":  "5",
}

// Category returns the category converted from nyaa one to tosho one
func Category(t *models.Torrent) string {
	cat := fmt.Sprintf("%d_%d", t.Category, t.SubCategory)
	if config.IsSukebei() {
		if val, ok := sukebeiCategories[cat]; ok {
			return val
		}
	}
	if val, ok := normalCategories[cat]; ok {
		return val
	}
	return ""
}
