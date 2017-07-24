package fosite

import (
	"context"
	"net/http"
	"strings"

	"net/url"

	"github.com/pkg/errors"
)

// Implements
// * https://tools.ietf.org/html/rfc6749#section-2.3.1
//   Clients in possession of a client password MAY use the HTTP Basic
//   authentication scheme as defined in [RFC2617] to authenticate with
//   the authorization server.  The client identifier is encoded using the
//   "application/x-www-form-urlencoded" encoding algorithm per
//   Appendix B, and the encoded value is used as the username; the client
//   password is encoded using the same algorithm and used as the
//   password.  The authorization server MUST support the HTTP Basic
//   authentication scheme for authenticating clients that were issued a
//   client password.
//   Including the client credentials in the request-body using the two
//   parameters is NOT RECOMMENDED and SHOULD be limited to clients unable
//   to directly utilize the HTTP Basic authentication scheme (or other
//   password-based HTTP authentication schemes).  The parameters can only
//   be transmitted in the request-body and MUST NOT be included in the
//   request URI.
//   * https://tools.ietf.org/html/rfc6749#section-3.2.1
//   - Confidential clients or other clients issued client credentials MUST
//   authenticate with the authorization server as described in
//   Section 2.3 when making requests to the token endpoint.
//   - If the client type is confidential or the client was issued client
//   credentials (or assigned other authentication requirements), the
//   client MUST authenticate with the authorization server as described
//   in Section 3.2.1.
func (f *Fosite) NewAccessRequest(ctx context.Context, r *http.Request, session Session) (AccessRequester, error) {
	var err error
	accessRequest := NewAccessRequest(session)

	if r.Method != "POST" {
		return accessRequest, errors.Wrap(ErrInvalidRequest, "HTTP method is not POST")
	} else if err := r.ParseForm(); err != nil {
		return accessRequest, errors.Wrap(ErrInvalidRequest, err.Error())
	}

	accessRequest.Form = r.PostForm
	if session == nil {
		return accessRequest, errors.New("Session must not be nil")
	}

	accessRequest.SetRequestedScopes(removeEmpty(strings.Split(r.PostForm.Get("scope"), " ")))
	accessRequest.GrantTypes = removeEmpty(strings.Split(r.PostForm.Get("grant_type"), " "))
	if len(accessRequest.GrantTypes) < 1 {
		return accessRequest, errors.Wrap(ErrInvalidRequest, "No grant type given")
	}

	// Decode client_id and client_secret which should be in "application/x-www-form-urlencoded" format.
	var clientID, clientSecret string
	if id, secret, ok := r.BasicAuth(); !ok {
		return accessRequest, errors.Wrap(ErrInvalidRequest, "HTTP authorization header missing or invalid")
	} else if clientID, err = url.QueryUnescape(id); err != nil {
		return accessRequest, errors.Wrap(ErrInvalidRequest, `The client id in the HTTP authorization header could not be decoded from "application/x-www-form-urlencoded"`)
	} else if clientSecret, err = url.QueryUnescape(secret); err != nil {
		return accessRequest, errors.Wrap(ErrInvalidRequest, `The client secret in the HTTP authorization header could not be decoded from "application/x-www-form-urlencoded"`)
	}

	client, err := f.Store.GetClient(ctx, clientID)
	if err != nil {
		return accessRequest, errors.Wrap(ErrInvalidClient, err.Error())
	}

	if !client.IsPublic() {
		// Enforce client authentication
		if err := f.Hasher.Compare(client.GetHashedSecret(), []byte(clientSecret)); err != nil {
			return accessRequest, errors.Wrap(ErrInvalidClient, err.Error())
		}
	}
	accessRequest.Client = client

	var found bool = false
	for _, loader := range f.TokenEndpointHandlers {
		if err := loader.HandleTokenEndpointRequest(ctx, accessRequest); err == nil {
			found = true
		} else if errors.Cause(err) == ErrUnknownRequest {
			// do nothing
		} else if err != nil {
			return accessRequest, err
		}
	}

	if !found {
		return nil, errors.WithStack(ErrInvalidRequest)
	}
	return accessRequest, nil
}
