package categories

import (
	"sort"

	"github.com/NyaaPantsu/nyaa/config"
)

// Category is a struct defining a category
type Category struct {
	ID   string
	Name string
}

// Cateogories is a struct defining an array of categories
type Categories []Category

var categories Categories
var Index map[string]int

func initCategories() {
	var cats map[string]string
	if config.IsSukebei() {
		cats = config.Conf.Torrents.SukebeiCategories
	} else {
		cats = config.Conf.Torrents.CleanCategories
	}

	// Sorting categories alphabetically
	var index []string
	ids := make(map[string]string)
	Index = make(map[string]int, len(cats))
	for id, name := range cats {
		index = append(index, name)
		ids[name] = id
	}
	sort.Strings(index)

	// Creating index of categories
	for k, name := range index {
		categories = append(categories, Category{ids[name], name})
		Index[ids[name]] = k
	}
}

// All : function to get all categories depending on the actual website from config/categories.go
func All() Categories {
	if len(categories) == 0 {
		initCategories()
	}
	return categories
}

// Get : function to get a category by the key in the index array
func Get(key int) Category {
	return All()[key]
}

// Get : function to get a category by the id of the category from the database
func GetByID(id string) (Category, bool) {
	if key, ok := Index[id]; ok {
		return All()[key], true
	}
	return Category{"", ""}, false
}

// Exists : Check if a category exist in config
func (cats Categories) Exists(id string) bool {
	for _, cat := range cats {
		if cat.ID == id {
			return true
		}
	}
	return false
}

// Exists : Check if a category exist in config
func Exists(category string) bool {
	return All().Exists(category)
}

// GetSelect : Format categories in map ordered alphabetically
func GetSelect(keepParent bool, keepChild bool) Categories {
	var catSelect Categories
	for _, v := range All() {
		if (keepParent && keepChild) || (len(v.ID) > 2 && !keepParent) || (len(v.ID) <= 2 && !keepChild) {
			catSelect = append(catSelect, v)
		}
	}
	return catSelect
}
