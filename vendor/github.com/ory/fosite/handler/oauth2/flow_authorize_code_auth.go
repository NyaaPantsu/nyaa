package oauth2

import (
	"strings"
	"time"

	"fmt"

	"context"

	"github.com/ory/fosite"
	"github.com/pkg/errors"
)

// AuthorizeExplicitGrantTypeHandler is a response handler for the Authorize Code grant using the explicit grant type
// as defined in https://tools.ietf.org/html/rfc6749#section-4.1
type AuthorizeExplicitGrantHandler struct {
	AccessTokenStrategy   AccessTokenStrategy
	RefreshTokenStrategy  RefreshTokenStrategy
	AuthorizeCodeStrategy AuthorizeCodeStrategy
	CoreStorage           CoreStorage

	// AuthCodeLifespan defines the lifetime of an authorize code.
	AuthCodeLifespan time.Duration

	// AccessTokenLifespan defines the lifetime of an access token.
	AccessTokenLifespan time.Duration

	ScopeStrategy fosite.ScopeStrategy
}

func (c *AuthorizeExplicitGrantHandler) HandleAuthorizeEndpointRequest(ctx context.Context, ar fosite.AuthorizeRequester, resp fosite.AuthorizeResponder) error {
	// This let's us define multiple response types, for example open id connect's id_token
	if !ar.GetResponseTypes().Exact("code") {
		return nil
	}

	if !ar.GetClient().GetResponseTypes().Has("code") {
		return errors.WithStack(fosite.ErrInvalidGrant)
	}

	if !fosite.IsRedirectURISecure(ar.GetRedirectURI()) {
		return errors.Wrap(fosite.ErrInvalidRequest, "Redirect URL is using an insecure protocol, http is only allowed for hosts with suffix `localhost`, for example: http://myapp.localhost/")
	}

	client := ar.GetClient()
	for _, scope := range ar.GetRequestedScopes() {
		if !c.ScopeStrategy(client.GetScopes(), scope) {
			return errors.Wrap(fosite.ErrInvalidScope, fmt.Sprintf("The client is not allowed to request scope %s", scope))
		}
	}

	return c.IssueAuthorizeCode(ctx, ar, resp)
}

func (c *AuthorizeExplicitGrantHandler) IssueAuthorizeCode(ctx context.Context, ar fosite.AuthorizeRequester, resp fosite.AuthorizeResponder) error {
	code, signature, err := c.AuthorizeCodeStrategy.GenerateAuthorizeCode(ctx, ar)
	if err != nil {
		return errors.Wrap(fosite.ErrServerError, err.Error())
	}

	ar.GetSession().SetExpiresAt(fosite.AuthorizeCode, time.Now().Add(c.AuthCodeLifespan))
	if err := c.CoreStorage.CreateAuthorizeCodeSession(ctx, signature, ar); err != nil {
		return errors.Wrap(fosite.ErrServerError, err.Error())
	}

	resp.AddQuery("code", code)
	resp.AddQuery("state", ar.GetState())
	resp.AddQuery("scope", strings.Join(ar.GetGrantedScopes(), " "))
	ar.SetResponseTypeHandled("code")
	return nil
}
