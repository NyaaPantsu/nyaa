package models

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"time"

	"github.com/NyaaPantsu/nyaa/utils/fosite/client"
	"github.com/NyaaPantsu/nyaa/utils/log"
	"github.com/ory/fosite"
	"github.com/pkg/errors"
)
const (
	TableOpenID  = "oidc"
	TableAccess  = "access"
	TableRefresh = "refresh"
	TableCode    = "code"
)

const oauth_prefix = "hydra_oauth2_"

type OauthAbstract struct {
	Signature    string    `gorm:"column:signature;primary_key;not null"`
	RequestID    string    `gorm:"column:request_id;not null"`
	RequestedAt  time.Time `gorm:"column:requested_at;not null;type:timestamp;default:CURRENT_TIMESTAMP"`
	ClientID     string    `gorm:"column:client_id;not null"`
	Scope        string    `gorm:"column:scope;not null"`
	GrantedScope string    `gorm:"column:granted_scope;not null"`
	FormData     string    `gorm:"column:form_data;not null"`
	SessionData  []byte    `gorm:"column:session_data;not null"`
}

type OpenID struct {
	OauthAbstract
}

type Access struct {
	OauthAbstract
}

type Code struct {
	OauthAbstract
}

type Refresh struct {
	OauthAbstract
}

func (o OpenID) TableName() string {
	return oauth_prefix + TableOpenID
}

func (a Access) TableName() string {
	return oauth_prefix + TableAccess
}

func (r Refresh) TableName() string {
	return oauth_prefix + TableRefresh
}

func (c Code) TableName() string {
	return oauth_prefix + TableCode
}

func (s *OauthAbstract) ToRequest(session fosite.Session, cm Manager) (*fosite.Request, error) {
	if session != nil {
		if err := json.Unmarshal(s.SessionData, session); err != nil {
			return nil, errors.WithStack(err)
		}
	} else {
		log.Debugf("Got an empty session in toRequest")
	}

	c, err := cm.GetClient(context.Background(), s.ClientID)
	if err != nil {
		return nil, err
	}

	val, err := url.ParseQuery(s.FormData)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	r := &fosite.Request{
		ID:            s.RequestID,
		RequestedAt:   s.RequestedAt,
		Client:        c,
		Scopes:        fosite.Arguments(strings.Split(s.Scope, "|")),
		GrantedScopes: fosite.Arguments(strings.Split(s.GrantedScope, "|")),
		Form:          val,
		Session:       session,
	}

	return r, nil
}

type Manager interface {
	Storage

	Authenticate(id string, secret []byte) (*client.Client, error)
}

type Storage interface {
	fosite.Storage

	CreateClient(c *client.Client) error

	UpdateClient(c *client.Client) error

	DeleteClient(id string) error

	GetClients() (map[string]client.Client, error)

	GetConcreteClient(id string) (*client.Client, error)
}
