package oauth_server

import (
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/ory/fosite"
	"github.com/pkg/errors"
)

type modelAbstract interface {
	TableName() string
	ToRequest(session fosite.Session, cm models.Manager) (*fosite.Request, error)
}

func SelectModel(table string, data models.OauthAbstract) modelAbstract {

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

func FindBySignature(signature string, table string) (modelAbstract, error) {
	d := SelectModel(table, models.OauthAbstract{})
	err := models.ORM.Where("signature = ?", signature).Find(d).Error
	if err != nil {
		return nil, errors.Wrap(fosite.ErrNotFound, "")
	}
	return d, nil
}
