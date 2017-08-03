package moderatorController

import (
	"github.com/NyaaPantsu/nyaa/controllers/middlewares"
	"github.com/NyaaPantsu/nyaa/controllers/router"
)

func init() {
	// We create a subroute group to apply a modmiddleware on all routes inside it
	// The following routes all start at /mod
	modRoutes := router.Get().Group("/mod", middlewares.ModMiddleware())
	{
		/* Panel index route */
		modRoutes.Any("/", IndexModPanel)

		/* Torrent Listing routes */
		modRoutes.GET("/torrents", TorrentsListPanel)
		modRoutes.GET("/torrents/p/:page", TorrentsListPanel)
		modRoutes.POST("/torrents", TorrentsPostListPanel)
		modRoutes.POST("/torrents/p/:page", TorrentsPostListPanel)

		/* Deleted torrents listing routes */
		modRoutes.GET("/torrents/deleted", DeletedTorrentsModPanel)
		modRoutes.GET("/torrents/deleted/p/:page", DeletedTorrentsModPanel)
		modRoutes.POST("/torrents/deleted", DeletedTorrentsPostPanel)
		modRoutes.POST("/torrents/deleted/p/:page", DeletedTorrentsPostPanel)

		/* Report listing routes */
		modRoutes.Any("/reports", TorrentReportListPanel)
		modRoutes.Any("/reports/p/:page", TorrentReportListPanel)

		/* User listing routes */
		modRoutes.Any("/users", UsersListPanel)
		modRoutes.Any("/users/p/:page", UsersListPanel)

		/* Comments listing routes */
		modRoutes.Any("/comments", CommentsListPanel)
		modRoutes.Any("/comments/p/:page", CommentsListPanel)
		modRoutes.Any("/comment", CommentsListPanel) // TODO Edit comment view

		/* Announcement listing routes */
		modRoutes.Any("/announcement", listAnnouncements)
		modRoutes.Any("/announcement/p/:page", listAnnouncements)

		/* Torrent edit view */
		modRoutes.GET("/torrent", TorrentEditModPanel)
		modRoutes.POST("/torrent", TorrentPostEditModPanel)

		/* Torrent delete routes */
		modRoutes.Any("/torrent/delete", TorrentDeleteModPanel)

		/* Announcement edit view */
		modRoutes.GET("/announcement/form", addAnnouncement)
		modRoutes.POST("/announcement/form", postAnnouncement)

		/* Announcement delete routes */
		modRoutes.Any("/announcement/delete", deleteAnnouncement)

		/* Torrent lock/unlock route */
		modRoutes.Any("/torrent/block", TorrentBlockModPanel)

		/* Tags delete route */
		modRoutes.Any("/tags/delete", DeleteTagsModPanel)

		/* Report delete route */
		modRoutes.Any("/report/delete", TorrentReportDeleteModPanel)

		/* Comment delete route */
		modRoutes.Any("/comment/delete", CommentDeleteModPanel)

		/* Reassign form routes */
		modRoutes.GET("/reassign", TorrentReassignModPanel)
		modRoutes.POST("/reassign", TorrentPostReassignModPanel)

		/* Oauth clients listing routes */
		modRoutes.GET("/oauth_client", clientsListPanel)
		modRoutes.GET("/oauth_client/p/:page", clientsListPanel)

		/* Oauth client delete route */
		modRoutes.GET("/oauth_client/delete", clientsDeleteModPanel)

		/* Oauth client edit routes */
		modRoutes.GET("/oauth_client/form", formClientController)
		modRoutes.POST("/oauth_client/form", formPostClientController)

		/* Mod Api routes */
		apiMod := modRoutes.Group("/api")
		apiMod.Any("/torrents", APIMassMod)
	}
}
