package categories

import (
	"github.com/NyaaPantsu/nyaa/config"
)

var categories map[string]string

// GetCategories : function to get all categories depending on the actual website from config/categories.go
func GetCategories() map[string]string {
	if categories != nil {
		return categories
	}

	if config.IsSukebei() {
		categories = config.Conf.Torrents.SukebeiCategories
	} else {
		categories = config.Conf.Torrents.CleanCategories
	}

	return categories
}

// CategoryExists : Check if a category exist in config
func CategoryExists(category string) bool {
	_, exists := GetCategories()[category]
	return exists
}

// GetCategoriesSelect : Format categories in map ordered alphabetically
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
