package themeToggleController

import "github.com/NyaaPantsu/nyaa/controllers/router"

func init() {
	router.Get().Any("/dark", toggleThemeHandler)
	router.Get().Any("/dark/*redirect", toggleThemeHandler)
}
