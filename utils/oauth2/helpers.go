package oauth2

import (
	"crypto/rand"
	"crypto/rsa"
	"time"

	"github.com/NyaaPantsu/nyaa/utils/oauth2/storage"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/utils/oauth2/manager"
	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/token/jwt"
)

// Store This is a storage instance. We will add a client and a user to it so we can use these later on.
var Store *storage.FositeSQLStore

// OauthConfig for oauth2 provider
var OauthConfig *compose.Config

// Strat Because we are using oauth2 and open connect id, we use this little helper to combine the two in one
// variable.
var Strat compose.CommonStrategy

// Oauth2 is the exported oauth provider
var Oauth2 fosite.OAuth2Provider

func init() {
	Store = &storage.FositeSQLStore{
		&manager.SQLManager{&fosite.BCrypt{WorkFactor: 12}},
	}
	OauthConfig = new(compose.Config)
	Strat = compose.CommonStrategy{
		// alternatively you could use:
		//  OAuth2Strategy: compose.NewOAuth2JWTStrategy(MustRSAKey())
		CoreStrategy: compose.NewOAuth2HMACStrategy(OauthConfig, config.OAuthHash),

		// open id connect strategy
		OpenIDConnectTokenStrategy: compose.NewOpenIDConnectStrategy(MustRSAKey()),
	}
	Oauth2 = compose.Compose(
		OauthConfig,
		Store,
		Strat,
		nil,

		// enabled handlers
		compose.OAuth2AuthorizeExplicitFactory,
		compose.OAuth2AuthorizeImplicitFactory,
		compose.OAuth2ClientCredentialsGrantFactory,
		compose.OAuth2RefreshTokenGrantFactory,
		compose.OAuth2ResourceOwnerPasswordCredentialsFactory,

		compose.OAuth2TokenRevocationFactory,
		compose.OAuth2TokenIntrospectionFactory,

		// be aware that open id connect factories need to be added after oauth2 factories to work properly.
		compose.OpenIDConnectExplicitFactory,
		compose.OpenIDConnectImplicitFactory,
		compose.OpenIDConnectHybridFactory,
		compose.OpenIDConnectRefreshFactory,
	)
}

// A session is passed from the `/auth` to the `/token` endpoint. You probably want to store data like: "Who made the request",
// "What organization does that person belong to" and so on.
// For our use case, the session will meet the requirements imposed by JWT access tokens, HMAC access tokens and OpenID Connect
// ID Tokens plus a custom field

// NewSession is a helper function for creating a new session. This may look like a lot of code but since we are
// setting up multiple strategies it is a bit longer.
// Usually, you could do:
//
//  session = new(fosite.DefaultSession)
func NewSession(user string, clientURI string) *openid.DefaultSession {
	return &openid.DefaultSession{
		Claims: &jwt.IDTokenClaims{
			Issuer:    "https://nyaa.pantsu.cat/",
			Subject:   user,
			Audience:  clientURI,
			ExpiresAt: time.Now().Add(time.Hour * 6),
			IssuedAt:  time.Now(),
		},
		Headers: &jwt.Headers{
			Extra: make(map[string]interface{}),
		},
	}
}

// MustRSAKey must generate a RSA key or panic!
func MustRSAKey() *rsa.PrivateKey {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	return key
}
