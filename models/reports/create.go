package reports

import (
	"time"

	"github.com/NyaaPantsu/nyaa/models"
	"errors"
)

func Create(desc string, message string, torrent *models.Torrent, user *models.User) (*models.TorrentReport, error) {
	report := &models.TorrentReport{
		Description: desc,
		Message:     message,
		TorrentID:   torrent.ID,
		UserID:      user.ID,
		CreatedAt:   time.Now(),
	}
	if models.ORM.Create(report).Error != nil {
		return report, errors.New("torrent_report_not_created")
	}
	return report, nil
}
