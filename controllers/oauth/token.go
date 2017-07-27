package oauth

import (
	"github.com/NyaaPantsu/nyaa/utils/log"
	"github.com/NyaaPantsu/nyaa/utils/oauth2"
	"github.com/gin-gonic/gin"
	"github.com/ory/fosite"
)

func tokenEndpoint(c *gin.Context) {
	// This context will be passed to all methods.
	ctx := fosite.NewContext()

	// Create an empty session object which will be passed to the request handlers
	mySessionData := oauth2.NewSession("", "")

	// This will create an access request object and iterate through the registered TokenEndpointHandlers to validate the request.
	accessRequest, err := oauth2.Oauth2.NewAccessRequest(ctx, c.Request, mySessionData)

	// Catch any errors, e.g.:
	// * unknown client
	// * invalid redirect
	// * ...
	if err != nil {
		log.Errorf("Error occurred in NewAccessRequest: %s\n", err)
		oauth2.Oauth2.WriteAccessError(c.Writer, accessRequest, err)
		return
	}

	// If this is a client_credentials grant, grant all scopes the client is allowed to perform.
	if accessRequest.GetGrantTypes().Exact("client_credentials") {
		for _, scope := range accessRequest.GetRequestedScopes() {
			if fosite.HierarchicScopeStrategy(accessRequest.GetClient().GetScopes(), scope) {
				accessRequest.GrantScope(scope)
			}
		}
	}

	// Next we create a response for the access request. Again, we iterate through the TokenEndpointHandlers
	// and aggregate the result in response.
	response, err := oauth2.Oauth2.NewAccessResponse(ctx, accessRequest)
	if err != nil {
		log.Errorf("Error occurred in NewAccessResponse: %s\n", err)
		oauth2.Oauth2.WriteAccessError(c.Writer, accessRequest, err)
		return
	}

	// All done, send the response.
	oauth2.Oauth2.WriteAccessResponse(c.Writer, accessRequest, response)

	// The client now has a valid access token
}
