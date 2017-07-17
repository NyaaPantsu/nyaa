package reportController

import "github.com/NyaaPantsu/nyaa/controllers/router"

func init() {
	reportRoutes := router.Get().Group("/report")
	{
		//reporting a torrent
		reportRoutes.GET("/:id", ReportViewTorrentHandler)
		reportRoutes.POST("/:id", ReportTorrentHandler)
	}
}
