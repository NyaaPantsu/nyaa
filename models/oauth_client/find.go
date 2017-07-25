package oauth_client

import (
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/pkg/errors"
)

func FindByID(id string) (*models.OauthClient, error) {
	d := &models.OauthClient{}
	err := models.ORM.Where("id = ?", id).Find(d).Error
	if err != nil {
		return nil, errors.Wrap(ErrNotFound, "")
	}
	return d, nil
}
