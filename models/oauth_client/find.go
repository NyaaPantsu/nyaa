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

// FindAll : Find all comments based on conditions
func FindAll(limit int, offset int, conditions string, values ...interface{}) ([]models.OauthClient, int, error) {
	var clients []models.OauthClient
	var nbClients int
	models.ORM.Model(&clients).Where(conditions, values...).Count(&nbClients)
	err := models.ORM.Limit(limit).Offset(offset).Where(conditions, values...).Find(&clients).Error
	if err != nil {
		return clients, nbClients, err
	}
	return clients, nbClients, nil
}
