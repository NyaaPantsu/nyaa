package torrents

import (
	"time"

	"github.com/NyaaPantsu/nyaa/models/tag"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/utils/log"
	"github.com/NyaaPantsu/nyaa/utils/validator/torrent"
)

// Create a new torrent based on the uploadform request struct
func Create(user *models.User, uploadForm *torrentValidator.TorrentRequest) (*models.Torrent, error) {
	torrent := models.Torrent{
		Name:        uploadForm.Name,
		Category:    uploadForm.CategoryID,
		SubCategory: uploadForm.SubCategoryID,
		Status:      uploadForm.Status,
		Hidden:      uploadForm.Hidden,
		Hash:        uploadForm.Infohash,
		Date:        time.Now(),
		Filesize:    uploadForm.Filesize,
		Languages:   uploadForm.Languages,
		Description: uploadForm.Description,
		WebsiteLink: uploadForm.WebsiteLink,
		UploaderID:  user.ID}
	torrent.EncodeLanguages() // Convert languages array in language string
	torrent.ParseTrackers(uploadForm.Trackers)
	for _, tagForm := range uploadForm.Tags {
		tag := &models.Tag{
			Tag:       tagForm.Tag,
			Type:      tagForm.Type,
			Accepted:  true,
			TorrentID: torrent.ID,
			UserID:    0, // 0 so we don't increase pantsu points for every tag for the actual user (would be too much increase)
			Weight:    config.Get().Torrents.Tags.MaxWeight + 1,
		}
		if tags.FilterOrCreate(tag, &torrent, user) { // We create a tag (filter doesn't apply since new torrent), only callbackOnType is called
			torrent.Tags = append(torrent.Tags, *tag) // Finally we append it to the torrent
		}
	}

	err := models.ORM.Create(&torrent).Error
	log.Infof("Torrent ID %d created!\n", torrent.ID)
	if err != nil {
		log.CheckErrorWithMessage(err, "ERROR_TORRENT_CREATE: Cannot create a torrent")
	}
	if config.Get().Search.EnableElasticSearch && models.ElasticSearchClient != nil {
		err := torrent.AddToESIndex(models.ElasticSearchClient)
		if err == nil {
			log.Infof("Successfully added torrent to ES index.")
		} else {
			log.Errorf("Unable to add torrent to ES index: %s", err)
		}
	} else {
		log.Error("Unable to create elasticsearch client")
	}
	NewTorrentEvent(user, &torrent)
	if len(uploadForm.FileList) > 0 {
		for _, uploadedFile := range uploadForm.FileList {
			file := models.File{TorrentID: torrent.ID, Filesize: uploadedFile.Filesize}
			err := file.SetPath(uploadedFile.Path)
			if err != nil {
				return &torrent, err
			}
			models.ORM.Create(&file)
		}
	}

	torrent.Update(false)
	user.IncreasePantsu()

	return &torrent, nil
}
