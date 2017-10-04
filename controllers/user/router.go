package userController

import "github.com/NyaaPantsu/nyaa/controllers/router"
import "github.com/NyaaPantsu/nyaa/controllers/feed"
import "github.com/NyaaPantsu/nyaa/controllers/search"

func init() {

	// Login
	router.Get().POST("/login", UserLoginPostHandler)
	router.Get().GET("/login", UserLoginFormHandler)

	// Register
	router.Get().GET("/register", UserRegisterFormHandler)
	router.Get().POST("/register", UserRegisterPostHandler)

	// Logout
	router.Get().POST("/logout", UserLogoutHandler)

	// Notifications
	router.Get().GET("/notifications", UserNotificationsHandler)

	// Verify Email
	router.Get().Any("/verify/email/:token", UserVerifyEmailHandler)

	// User Profile specific routes
	userRoutes := router.Get().Group("/user")
	{
		userRoutes.GET("/:id", UserProfileHandler)
		userRoutes.GET("/:id/:username", UserProfileHandler)
		userRoutes.GET("/:id/:username/follow", UserFollowHandler)
		userRoutes.GET("/:id/:username/edit", UserDetailsHandler)
		userRoutes.POST("/:id/:username/edit", UserProfileFormHandler)
		userRoutes.GET("/:id/:username/apireset", UserAPIKeyResetHandler)
		userRoutes.GET("/:id/:username/search", searchController.SearchHandler)
		userRoutes.GET("/:id/:username/search/:page", searchController.SearchHandler)
		userRoutes.GET("/:id/:username/feed", feedController.RSSHandler)
		userRoutes.GET("/:id/:username/feed/:page", feedController.RSSHandler)
	}
	
	router.Get().Any("/username", RedirectToUserSearch)
	router.Get().Any("/username/:username", UserGetFromName)
	router.Get().Any("/username/:username/search", searchController.SearchHandler)
	router.Get().Any("/username/:username/search:page", searchController.SearchHandler)
}
