package storage

import (
	"context"

	"encoding/json"
	"strings"

	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/oauth_server"
	"github.com/NyaaPantsu/nyaa/models/users"
	"github.com/NyaaPantsu/nyaa/utils/log"
	"github.com/ory/fosite"
	"github.com/pkg/errors"
)

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

func (s *FositeSQLStore) createSession(signature string, requester fosite.Requester, table string) error {
	data, err := sqlSchemaFromRequest(signature, requester)
	if err != nil {
		return err
	}

	model := oauth_server.SelectModel(table, *data)
	err = models.ORM.Create(model).Error

	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *FositeSQLStore) findSessionBySignature(signature string, session fosite.Session, table string) (fosite.Requester, error) {
	d, err := oauth_server.FindBySignature(signature, table)
	if err != nil {
		return nil, err
	}

	return d.ToRequest(session, s.Manager)
}

func (s *FositeSQLStore) deleteSession(signature string, table string) error {
	return oauth_server.DeleteBySession(signature, table)
}

func (s *FositeSQLStore) CreateOpenIDConnectSession(_ context.Context, signature string, requester fosite.Requester) error {
	return s.createSession(signature, requester, models.TableOpenID)
}

func (s *FositeSQLStore) GetOpenIDConnectSession(_ context.Context, signature string, requester fosite.Requester) (fosite.Requester, error) {
	return s.findSessionBySignature(signature, requester.GetSession(), models.TableOpenID)
}

func (s *FositeSQLStore) DeleteOpenIDConnectSession(_ context.Context, signature string) error {
	return s.deleteSession(signature, models.TableOpenID)
}

func (s *FositeSQLStore) CreateAuthorizeCodeSession(_ context.Context, signature string, requester fosite.Requester) error {
	return s.createSession(signature, requester, models.TableCode)
}

func (s *FositeSQLStore) GetAuthorizeCodeSession(_ context.Context, signature string, session fosite.Session) (fosite.Requester, error) {
	return s.findSessionBySignature(signature, session, models.TableCode)
}

func (s *FositeSQLStore) DeleteAuthorizeCodeSession(_ context.Context, signature string) error {
	return s.deleteSession(signature, models.TableCode)
}

func (s *FositeSQLStore) CreateAccessTokenSession(_ context.Context, signature string, requester fosite.Requester) error {
	return s.createSession(signature, requester, models.TableAccess)
}

func (s *FositeSQLStore) GetAccessTokenSession(_ context.Context, signature string, session fosite.Session) (fosite.Requester, error) {
	return s.findSessionBySignature(signature, session, models.TableAccess)
}

func (s *FositeSQLStore) DeleteAccessTokenSession(_ context.Context, signature string) error {
	return s.deleteSession(signature, models.TableAccess)
}

func (s *FositeSQLStore) CreateRefreshTokenSession(_ context.Context, signature string, requester fosite.Requester) error {
	return s.createSession(signature, requester, models.TableRefresh)
}

func (s *FositeSQLStore) GetRefreshTokenSession(_ context.Context, signature string, session fosite.Session) (fosite.Requester, error) {
	return s.findSessionBySignature(signature, session, models.TableRefresh)
}

func (s *FositeSQLStore) DeleteRefreshTokenSession(_ context.Context, signature string) error {
	return s.deleteSession(signature, models.TableRefresh)
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
	return s.revokeSession(id, models.TableRefresh)
}

func (s *FositeSQLStore) RevokeAccessToken(ctx context.Context, id string) error {
	return s.revokeSession(id, models.TableAccess)
}

func (s *FositeSQLStore) revokeSession(id string, table string) error {
	return oauth_server.DeleteByID(id, table)
}

func (s *FositeSQLStore) Authenticate(_ context.Context, name string, secret string) error {
	_, _, err := users.Exists(name, secret)
	if err != nil {
		return fosite.ErrNotFound
	}
	return nil
}
