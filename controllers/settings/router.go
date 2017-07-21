package settingsController

import "github.com/NyaaPantsu/nyaa/controllers/router"

func init() {
	router.Get().GET("/settings", SeePublicSettingsHandler)
	router.Get().POST("/settings", ChangePublicSettingsHandler)
}
