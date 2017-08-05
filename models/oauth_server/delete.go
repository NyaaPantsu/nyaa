package oauth_server

import (
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/ory/fosite"
	"github.com/pkg/errors"
)

func DeleteBySession(signature string, table string) error {
	d := SelectModel(table, models.OauthAbstract{})
	err := models.ORM.Where("signature = ?", signature).Delete(d).Error

	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func DeleteByID(id string, table string) error {
	err := models.ORM.Where("request_id", id).Delete(SelectModel(table, models.OauthAbstract{})).Error
	if err != nil {
		return errors.Wrap(fosite.ErrNotFound, "")
	}
	return nil
}
