package storage

import (
	"context"

	"encoding/json"
	"strings"

	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/users"
	"github.com/NyaaPantsu/nyaa/utils/log"
	"github.com/ory/fosite"
	"github.com/pkg/errors"
)

const (
	sqlTableOpenID  = "oidc"
	sqlTableAccess  = "access"
	sqlTableRefresh = "refresh"
	sqlTableCode    = "code"
)

type modelAbstract interface {
	TableName() string
	ToRequest(session fosite.Session, cm models.Manager) (*fosite.Request, error)
}

type FositeSQLStore struct {
	models.Manager
}

func sqlSchemaFromRequest(signature string, r fosite.Requester) (*models.OauthAbstract, error) {
	if r.GetSession() == nil {
		log.Debugf("Got an empty session in sqlSchemaFromRequest")
	}

	session, err := json.Marshal(r.GetSession())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &models.OauthAbstract{
		RequestID:    r.GetID(),
		Signature:    signature,
		RequestedAt:  r.GetRequestedAt(),
		ClientID:     r.GetClient().GetID(),
		Scope:        strings.Join([]string(r.GetRequestedScopes()), "|"),
		GrantedScope: strings.Join([]string(r.GetGrantedScopes()), "|"),
		FormData:     r.GetRequestForm().Encode(),
		SessionData:  session,
	}, nil

}

func selectModel(table string, data models.OauthAbstract) modelAbstract {

	switch table {
	case "oidc":
		return &models.OpenID{data}
	case "access":
		return &models.Access{data}
	case "refresh":
		return &models.Refresh{data}
	case "code":
		return &models.Code{data}
	}
	return nil
}

func (s *FositeSQLStore) createSession(signature string, requester fosite.Requester, table string) error {
	data, err := sqlSchemaFromRequest(signature, requester)
	if err != nil {
		return err
	}

	model := selectModel(table, *data)
	err = models.ORM.Create(model).Error

	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *FositeSQLStore) findSessionBySignature(signature string, session fosite.Session, table string) (fosite.Requester, error) {
	d := selectModel(table, models.OauthAbstract{})
	err := models.ORM.Where("signature = ?", signature).Find(d).Error
	if err != nil {
		return nil, errors.Wrap(fosite.ErrNotFound, "")
	}

	return d.ToRequest(session, s.Manager)
}

func (s *FositeSQLStore) deleteSession(signature string, table string) error {
	d := selectModel(table, models.OauthAbstract{})
	err := models.ORM.Where("signature = ?", signature).Delete(d).Error

	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *FositeSQLStore) CreateOpenIDConnectSession(_ context.Context, signature string, requester fosite.Requester) error {
	return s.createSession(signature, requester, sqlTableOpenID)
}

func (s *FositeSQLStore) GetOpenIDConnectSession(_ context.Context, signature string, requester fosite.Requester) (fosite.Requester, error) {
	return s.findSessionBySignature(signature, requester.GetSession(), sqlTableOpenID)
}

func (s *FositeSQLStore) DeleteOpenIDConnectSession(_ context.Context, signature string) error {
	return s.deleteSession(signature, sqlTableOpenID)
}

func (s *FositeSQLStore) CreateAuthorizeCodeSession(_ context.Context, signature string, requester fosite.Requester) error {
	return s.createSession(signature, requester, sqlTableCode)
}

func (s *FositeSQLStore) GetAuthorizeCodeSession(_ context.Context, signature string, session fosite.Session) (fosite.Requester, error) {
	return s.findSessionBySignature(signature, session, sqlTableCode)
}

func (s *FositeSQLStore) DeleteAuthorizeCodeSession(_ context.Context, signature string) error {
	return s.deleteSession(signature, sqlTableCode)
}

func (s *FositeSQLStore) CreateAccessTokenSession(_ context.Context, signature string, requester fosite.Requester) error {
	return s.createSession(signature, requester, sqlTableAccess)
}

func (s *FositeSQLStore) GetAccessTokenSession(_ context.Context, signature string, session fosite.Session) (fosite.Requester, error) {
	return s.findSessionBySignature(signature, session, sqlTableAccess)
}

func (s *FositeSQLStore) DeleteAccessTokenSession(_ context.Context, signature string) error {
	return s.deleteSession(signature, sqlTableAccess)
}

func (s *FositeSQLStore) CreateRefreshTokenSession(_ context.Context, signature string, requester fosite.Requester) error {
	return s.createSession(signature, requester, sqlTableRefresh)
}

func (s *FositeSQLStore) GetRefreshTokenSession(_ context.Context, signature string, session fosite.Session) (fosite.Requester, error) {
	return s.findSessionBySignature(signature, session, sqlTableRefresh)
}

func (s *FositeSQLStore) DeleteRefreshTokenSession(_ context.Context, signature string) error {
	return s.deleteSession(signature, sqlTableRefresh)
}

func (s *FositeSQLStore) CreateImplicitAccessTokenSession(ctx context.Context, signature string, requester fosite.Requester) error {
	return s.CreateAccessTokenSession(ctx, signature, requester)
}

func (s *FositeSQLStore) PersistAuthorizeCodeGrantSession(ctx context.Context, authorizeCode, accessSignature, refreshSignature string, request fosite.Requester) error {
	if err := s.DeleteAuthorizeCodeSession(ctx, authorizeCode); err != nil {
		return err
	} else if err := s.CreateAccessTokenSession(ctx, accessSignature, request); err != nil {
		return err
	}

	if refreshSignature == "" {
		return nil
	}

	if err := s.CreateRefreshTokenSession(ctx, refreshSignature, request); err != nil {
		return err
	}

	return nil
}

func (s *FositeSQLStore) PersistRefreshTokenGrantSession(ctx context.Context, originalRefreshSignature, accessSignature, refreshSignature string, request fosite.Requester) error {
	if err := s.DeleteRefreshTokenSession(ctx, originalRefreshSignature); err != nil {
		return err
	} else if err := s.CreateAccessTokenSession(ctx, accessSignature, request); err != nil {
		return err
	} else if err := s.CreateRefreshTokenSession(ctx, refreshSignature, request); err != nil {
		return err
	}

	return nil
}

func (s *FositeSQLStore) RevokeRefreshToken(ctx context.Context, id string) error {
	return s.revokeSession(id, sqlTableRefresh)
}

func (s *FositeSQLStore) RevokeAccessToken(ctx context.Context, id string) error {
	return s.revokeSession(id, sqlTableAccess)
}

func (s *FositeSQLStore) revokeSession(id string, table string) error {
	err := models.ORM.Where("request_id", id).Delete(selectModel(table, models.OauthAbstract{})).Error
	if err != nil {
		return errors.Wrap(fosite.ErrNotFound, "")
	}

	return nil
}

func (s *FositeSQLStore) Authenticate(_ context.Context, name string, secret string) error {
	_, _, err := users.Exists(name, secret)
	if err != nil {
		return fosite.ErrNotFound
	}
	return nil
}
