package upload

import (
	"fmt"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
)

// Convert automatically our sukebei cats to platform specified Hentai cats
var sukebeiCategories = []map[string]string{
	ttosho: {
		"1_1": "12",
		"1_2": "12",
		"1_3": "14",
		"1_4": "13",
		"1_5": "4",
		"2_1": "4",
		"2_2": "15",
	},
	nyaasi: {
		"1_1": "1_1",
		"1_2": "1_2",
		"1_3": "1_3",
		"1_4": "1_4",
		"1_5": "1_5",
		"2_1": "2_1",
		"2_2": "2_2",
	},
}

var normalCategories = []map[string]string{
	ttosho: {
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
	},
	nyaasi: {
		"3_12": "1_1",
		"3_5":  "1_2",
		"3_13": "1_3",
		"3_6":  "1_4",
		"2_3":  "2_1",
		"2_4":  "2_2",
		"4_7":  "3_1",
		"4_8":  "3_4",
		"4_14": "3_3",
		"5_9":  "4_1",
		"5_10": "4_2",
		"5_18": "4_3",
		"5_11": "4_4",
		"6_15": "5_1",
		"6_16": "5_2",
		"1_1":  "6_1",
		"1_2":  "6_2",
	},
}

// Category returns the category converted from nyaa one to tosho one
func Category(platform int, t *models.Torrent) string {
	cat := fmt.Sprintf("%d_%d", t.Category, t.SubCategory)
	// if we are in sukebei, there are some categories
	if config.IsSukebei() {
		// check that platform exist in our map for sukebei categories
		if platform < len(sukebeiCategories) {
			// return the remaped category if it exists
			if val, ok := sukebeiCategories[platform][cat]; ok {
				return val
			}
		}
	}
	// check that platform exist in our map
	if platform >= len(normalCategories) {
		return ""
	}
	// return the remaped category if it exists
	if val, ok := normalCategories[platform][cat]; ok {
		return val
	}
	return ""
}
