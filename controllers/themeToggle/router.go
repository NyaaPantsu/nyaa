package themeToggleController

import "github.com/NyaaPantsu/nyaa/controllers/router"

func init() {
	router.Get().Any("/dark", toggleThemeHandler)
}
