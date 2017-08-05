package torrents

import (
	"encoding/json"
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

	var tagsReq models.Tags
	json.Unmarshal([]byte(uploadForm.Tags), &tagsReq)
	for _, tag := range tagsReq {
		tag.Accepted = true
		tag.TorrentID = torrent.ID
		tag.Weight = config.Get().Torrents.Tags.MaxWeight
		tags.New(&tag, &torrent)            // We create new tags
		torrent.Tags = append(torrent.Tags) // Finally we append it to the torrent
	}
	user.IncreasePantsu()

	return &torrent, nil
}
