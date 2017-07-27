package oauth

import (
	"github.com/NyaaPantsu/nyaa/utils/oauth2"
	"github.com/gin-gonic/gin"
	"github.com/ory/fosite"
)

func revokeEndpoint(c *gin.Context) {
	// This context will be passed to all methods.
	ctx := fosite.NewContext()

	// This will accept the token revocation request and validate various parameters.
	err := oauth2.Oauth2.NewRevocationRequest(ctx, c.Request)

	// All done, send the response.
	oauth2.Oauth2.WriteRevocationResponse(c.Writer, err)
}
