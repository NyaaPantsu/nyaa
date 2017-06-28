package torrents


// DeleteTorrent : delete a torrent based on id
func DeleteTorrent(id uint) (*models.Torrent, int, error) {
	var torrent models.Torrent
	if models.ORM.First(&torrent, id).RecordNotFound() {
		return &torrent, http.StatusNotFound, errors.New("Torrent is not found")
	}
	if models.ORM.Delete(&torrent).Error != nil {
		return &torrent, http.StatusInternalServerError, errors.New("Torrent was not deleted")
	}

	if models.ElasticSearchClient != nil {
		err := torrent.DeleteFromESIndex(models.ElasticSearchClient)
		if err == nil {
			log.Infof("Successfully deleted torrent to ES index.")
		} else {
			log.Errorf("Unable to delete torrent to ES index: %s", err)
		}
	}
	return &torrent, http.StatusOK, nil
}

// DefinitelyDelete : deletes definitely a torrent based on id
func DefinitelyDelete(id uint) (*models.Torrent, int, error) {
	var torrent models.Torrent
	if models.ORM.Unscoped().Model(&torrent).First(&torrent, id).RecordNotFound() {
		return &torrent, http.StatusNotFound, errors.New("Torrent is not found")
	}
	if models.ORM.Unscoped().Model(&torrent).Delete(&torrent).Error != nil {
		return &torrent, http.StatusInternalServerError, errors.New("Torrent was not deleted")
	}

	if models.ElasticSearchClient != nil {
		err := torrent.DeleteFromESIndex(models.ElasticSearchClient)
		if err == nil {
			log.Infof("Successfully deleted torrent to ES index.")
		} else {
			log.Errorf("Unable to delete torrent to ES index: %s", err)
		}
	}
	return &torrent, http.StatusOK, nil
}
