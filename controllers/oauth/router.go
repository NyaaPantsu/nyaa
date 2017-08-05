package oauth

import "github.com/NyaaPantsu/nyaa/controllers/router"

func init() {
	oauth2Routes := router.Get().Group("/oauth2")
	{
		oauth2Routes.Any("/auth", authEndpoint)
		oauth2Routes.Any("/token", tokenEndpoint)

		// revoke tokens
		oauth2Routes.Any("/revoke", revokeEndpoint)
		oauth2Routes.Any("/introspect", introspectionEndpoint)
	}
}
