package oauth2

import (
	"context"
	"time"

	"github.com/ory/fosite"
	"github.com/pkg/errors"
)

// HandleTokenEndpointRequest implements
// * https://tools.ietf.org/html/rfc6749#section-4.1.3 (everything)
func (c *AuthorizeExplicitGrantHandler) HandleTokenEndpointRequest(ctx context.Context, request fosite.AccessRequester) error {
	// grant_type REQUIRED.
	// Value MUST be set to "authorization_code".
	if !request.GetGrantTypes().Exact("authorization_code") {
		return errors.WithStack(errors.WithStack(fosite.ErrUnknownRequest))
	}

	if !request.GetClient().GetGrantTypes().Has("authorization_code") {
		return errors.Wrap(errors.WithStack(fosite.ErrInvalidGrant), "The client is not allowed to use grant type authorization_code")
	}

	code := request.GetRequestForm().Get("code")
	signature := c.AuthorizeCodeStrategy.AuthorizeCodeSignature(code)
	authorizeRequest, err := c.CoreStorage.GetAuthorizeCodeSession(ctx, signature, request.GetSession())
	if errors.Cause(err) == fosite.ErrNotFound {
		return errors.Wrap(fosite.ErrInvalidRequest, err.Error())
	} else if err != nil {
		return errors.Wrap(errors.WithStack(fosite.ErrServerError), err.Error())
	}

	// The authorization server MUST verify that the authorization code is valid
	// This needs to happen after store retrieval for the session to be hydrated properly
	if err := c.AuthorizeCodeStrategy.ValidateAuthorizeCode(ctx, request, code); err != nil {
		return errors.Wrap(errors.WithStack(fosite.ErrInvalidRequest), err.Error())
	}

	// Override scopes
	request.SetRequestedScopes(authorizeRequest.GetRequestedScopes())

	// The authorization server MUST ensure that the authorization code was issued to the authenticated
	// confidential client, or if the client is public, ensure that the
	// code was issued to "client_id" in the request,
	if authorizeRequest.GetClient().GetID() != request.GetClient().GetID() {
		return errors.Wrap(errors.WithStack(fosite.ErrInvalidRequest), "Client ID mismatch")
	}

	// ensure that the "redirect_uri" parameter is present if the
	// "redirect_uri" parameter was included in the initial authorization
	// request as described in Section 4.1.1, and if included ensure that
	// their values are identical.
	forcedRedirectURI := authorizeRequest.GetRequestForm().Get("redirect_uri")
	if forcedRedirectURI != "" && forcedRedirectURI != request.GetRequestForm().Get("redirect_uri") {
		return errors.Wrap(errors.WithStack(fosite.ErrInvalidRequest), "Redirect URI mismatch")
	}

	// Checking of POST client_id skipped, because:
	// If the client type is confidential or the client was issued client
	// credentials (or assigned other authentication requirements), the
	// client MUST authenticate with the authorization server as described
	// in Section 3.2.1.
	request.SetSession(authorizeRequest.GetSession())
	request.GetSession().SetExpiresAt(fosite.AccessToken, time.Now().Add(c.AccessTokenLifespan))
	return nil
}

func (c *AuthorizeExplicitGrantHandler) PopulateTokenEndpointResponse(ctx context.Context, requester fosite.AccessRequester, responder fosite.AccessResponder) error {
	// grant_type REQUIRED.
	// Value MUST be set to "authorization_code".
	if !requester.GetGrantTypes().Exact("authorization_code") {
		return errors.WithStack(fosite.ErrUnknownRequest)
	}

	code := requester.GetRequestForm().Get("code")
	signature := c.AuthorizeCodeStrategy.AuthorizeCodeSignature(code)
	authorizeRequest, err := c.CoreStorage.GetAuthorizeCodeSession(ctx, signature, requester.GetSession())
	if err != nil {
		return errors.Wrap(errors.WithStack(fosite.ErrServerError), err.Error())
	} else if err := c.AuthorizeCodeStrategy.ValidateAuthorizeCode(ctx, requester, code); err != nil {
		// This needs to happen after store retrieval for the session to be hydrated properly
		return errors.Wrap(errors.WithStack(fosite.ErrInvalidRequest), err.Error())
	}

	for _, scope := range authorizeRequest.GetGrantedScopes() {
		requester.GrantScope(scope)
	}

	access, accessSignature, err := c.AccessTokenStrategy.GenerateAccessToken(ctx, requester)
	if err != nil {
		return errors.Wrap(errors.WithStack(fosite.ErrServerError), err.Error())
	}

	var refresh, refreshSignature string
	if authorizeRequest.GetGrantedScopes().Has("offline") {
		refresh, refreshSignature, err = c.RefreshTokenStrategy.GenerateRefreshToken(ctx, requester)
		if err != nil {
			return errors.Wrap(errors.WithStack(fosite.ErrServerError), err.Error())
		}
	}

	if err := c.CoreStorage.DeleteAuthorizeCodeSession(ctx, signature); err != nil {
		return errors.Wrap(errors.WithStack(fosite.ErrServerError), err.Error())
	} else if err := c.CoreStorage.CreateAccessTokenSession(ctx, accessSignature, requester); err != nil {
		return errors.Wrap(errors.WithStack(fosite.ErrServerError), err.Error())
	} else if refreshSignature != "" {
		if err := c.CoreStorage.CreateRefreshTokenSession(ctx, refreshSignature, requester); err != nil {
			return errors.Wrap(errors.WithStack(fosite.ErrServerError), err.Error())
		}
	}

	responder.SetAccessToken(access)
	responder.SetTokenType("bearer")
	responder.SetExpiresIn(getExpiresIn(requester, fosite.AccessToken, c.AccessTokenLifespan, time.Now()))
	responder.SetScopes(requester.GetGrantedScopes())
	if refresh != "" {
		responder.SetExtra("refresh_token", refresh)
	}

	return nil
}
