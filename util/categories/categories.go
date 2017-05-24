package categories

import (
	"github.com/NyaaPantsu/nyaa/config"
)

var categories map[string]string

func GetCategories() map[string]string {
	if categories != nil {
		return categories
	}

	if config.IsSukebei() {
		categories = config.TorrentSukebeiCategories
	} else {
		categories = config.TorrentCleanCategories
	}

	return categories
}

func CategoryExists(category string) bool {
	_, exists := GetCategories()[category]
	return exists
}

func GetCategoriesSelect(keepParent bool) map[string]string {
	categories := GetCategories()
	catSelect := make(map[string]string, len(categories))
	for k, v := range categories {
		if len(k) > 2 || keepParent {
			catSelect[v] = k
		}
	}
	return catSelect
}
