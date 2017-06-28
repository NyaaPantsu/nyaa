package torrents

import (
	"errors"
	"net/http"

	"github.com/NyaaPantsu/nyaa/util/log"

	"github.com/NyaaPantsu/nyaa/models"
)

// Update : Update a torrent based on model
func Update(torrent *models.Torrent) (int, error) {
	if models.ORM.Model(torrent).UpdateColumn(torrent).Error != nil {
		return http.StatusInternalServerError, errors.New("Torrent was not updated")
	}

	// TODO Don't create a new client for each request
	if models.ElasticSearchClient != nil {
		err := torrent.AddToESIndex(models.ElasticSearchClient)
		if err == nil {
			log.Infof("Successfully updated torrent to ES index.")
		} else {
			log.Errorf("Unable to update torrent to ES index: %s", err)
		}
	}

	return http.StatusOK, nil
}

// UpdateUnscope : Update a torrent based on model
func UpdateUnscope(torrent *models.Torrent) (int, error) {
	if models.ORM.Unscoped().Model(torrent).UpdateColumn(torrent).Error != nil {
		return http.StatusInternalServerError, errors.New("Torrent was not updated")
	}

	// TODO Don't create a new client for each request
	if models.ElasticSearchClient != nil {
		err := torrent.AddToESIndex(models.ElasticSearchClient)
		if err == nil {
			log.Infof("Successfully updated torrent to ES index.")
		} else {
			log.Errorf("Unable to update torrent to ES index: %s", err)
		}
	}

	return http.StatusOK, nil
}
