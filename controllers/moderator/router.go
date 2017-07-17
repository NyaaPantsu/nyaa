package moderatorController

import (
	"github.com/NyaaPantsu/nyaa/controllers/middlewares"
	"github.com/NyaaPantsu/nyaa/controllers/router"
)

func init() {
	modRoutes := router.Get().Group("/mod", middlewares.ModMiddleware())
	{
		modRoutes.Any("/", IndexModPanel)
		modRoutes.GET("/torrents", TorrentsListPanel)
		modRoutes.GET("/torrents/p/:page", TorrentsListPanel)
		modRoutes.POST("/torrents", TorrentsPostListPanel)
		modRoutes.POST("/torrents/p/:page", TorrentsPostListPanel)
		modRoutes.GET("/torrents/deleted", DeletedTorrentsModPanel)
		modRoutes.GET("/torrents/deleted/p/:page", DeletedTorrentsModPanel)
		modRoutes.POST("/torrents/deleted", DeletedTorrentsPostPanel)
		modRoutes.POST("/torrents/deleted/p/:page", DeletedTorrentsPostPanel)
		modRoutes.Any("/reports", TorrentReportListPanel)
		modRoutes.Any("/reports/p/:page", TorrentReportListPanel)
		modRoutes.Any("/users", UsersListPanel)
		modRoutes.Any("/users/p/:page", UsersListPanel)
		modRoutes.Any("/comments", CommentsListPanel)
		modRoutes.Any("/comments/p/:page", CommentsListPanel)
		modRoutes.Any("/comment", CommentsListPanel) // TODO
		modRoutes.GET("/torrent", TorrentEditModPanel)
		modRoutes.POST("/torrent", TorrentPostEditModPanel)
		modRoutes.Any("/torrent/delete", TorrentDeleteModPanel)
		modRoutes.Any("/torrent/block", TorrentBlockModPanel)
		modRoutes.Any("/report/delete", TorrentReportDeleteModPanel)
		modRoutes.Any("/comment/delete", CommentDeleteModPanel)
		modRoutes.GET("/reassign", TorrentReassignModPanel)
		modRoutes.POST("/reassign", TorrentPostReassignModPanel)
		apiMod := modRoutes.Group("/api")
		apiMod.Any("/torrents", APIMassMod)
	}
}
