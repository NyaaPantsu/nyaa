package middlewares

import "github.com/NyaaPantsu/nyaa/controllers/router"

func init() {
	router.Get().Use(CSP(), ErrorMiddleware())
}
