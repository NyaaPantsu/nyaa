package torrentController

import "github.com/NyaaPantsu/nyaa/controllers/router"

func init() {
	router.Get().Any("/download/:hash", DownloadTorrent)

	torrentRoutes := router.Get().Group("/torrent")
	{
		torrentRoutes.GET("/", TorrentEditUserPanel)
		torrentRoutes.POST("/", TorrentPostEditUserPanel)
		torrentRoutes.GET("/tag", ViewFormTag)
		torrentRoutes.POST("/tag", ViewFormTag)
		torrentRoutes.GET("/tag/add", AddTag)
		torrentRoutes.GET("/tag/remove", DeleteTag)
		torrentRoutes.GET("/delete", TorrentDeleteUserPanel)
	}
	torrentViewRoutes := router.Get().Group("/view")
	{
		torrentViewRoutes.GET("/:id", ViewHandler)
		torrentViewRoutes.HEAD("/:id", ViewHeadHandler)
		torrentViewRoutes.POST("/:id", PostCommentHandler)
	}
}
