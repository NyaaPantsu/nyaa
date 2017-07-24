package oauth

import (
	"github.com/NyaaPantsu/nyaa/utils/log"
	"github.com/gin-gonic/gin"
	"github.com/ory/fosite"
)

func introspectionEndpoint(c *gin.Context) {
	ctx := fosite.NewContext()
	mySessionData := newSession("", "")
	ir, err := oauth2.NewIntrospectionRequest(ctx, c.Request, mySessionData)
	if err != nil {
		log.Errorf("Error occurred in NewAuthorizeRequest: %s\n", err)
		oauth2.WriteIntrospectionError(c.Writer, err)
		return
	}

	oauth2.WriteIntrospectionResponse(c.Writer, ir)
}
