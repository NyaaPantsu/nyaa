package pprofController

import "github.com/NyaaPantsu/nyaa/controllers/router"
import "github.com/NyaaPantsu/nyaa/controllers/middlewares"

func init() {
	// Adding pprof support
	pprofRoutes := router.Get().Group("/debug/pprof", middlewares.ModMiddleware())
	{
		pprofRoutes.GET("/", PprofIndex)
		pprofRoutes.GET("/block", PprofIndex)
		pprofRoutes.GET("/heap", PprofIndex)
		pprofRoutes.GET("/profile", PprofProfile)
		pprofRoutes.POST("/symbol", PprofSymbol)
		pprofRoutes.GET("/symbol", PprofSymbol)
		pprofRoutes.GET("/trace", PprofTrace)
	}
}
