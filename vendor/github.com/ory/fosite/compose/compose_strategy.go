package compose

import (
	"crypto/rsa"

	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/token/hmac"
	"github.com/ory/fosite/token/jwt"
)

type CommonStrategy struct {
	oauth2.CoreStrategy
	openid.OpenIDConnectTokenStrategy
}

func NewOAuth2HMACStrategy(config *Config, secret []byte) *oauth2.HMACSHAStrategy {
	return &oauth2.HMACSHAStrategy{
		Enigma: &hmac.HMACStrategy{
			GlobalSecret: secret,
		},
		AccessTokenLifespan:   config.GetAccessTokenLifespan(),
		AuthorizeCodeLifespan: config.GetAuthorizeCodeLifespan(),
	}
}

func NewOAuth2JWTStrategy(key *rsa.PrivateKey, strategy *oauth2.HMACSHAStrategy) *oauth2.RS256JWTStrategy {
	return &oauth2.RS256JWTStrategy{
		RS256JWTStrategy: &jwt.RS256JWTStrategy{
			PrivateKey: key,
		},
		HMACSHAStrategy: strategy,
	}
}

func NewOpenIDConnectStrategy(key *rsa.PrivateKey) *openid.DefaultStrategy {
	return &openid.DefaultStrategy{
		RS256JWTStrategy: &jwt.RS256JWTStrategy{
			PrivateKey: key,
		},
	}
}
