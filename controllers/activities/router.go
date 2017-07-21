package activitiesController

import "github.com/NyaaPantsu/nyaa/controllers/router"

func init() {
	router.Get().Any("/activities", ActivityListHandler)
	router.Get().Any("/activities/p/:page", ActivityListHandler)
}
