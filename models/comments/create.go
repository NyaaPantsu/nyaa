package comments

import (
	"time"

	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/utils/cache"
)

func Create(content string, torrent *models.Torrent, user *models.User) (*models.Comment, error) {
	comment := &models.Comment{TorrentID: torrent.ID, UserID: user.ID, Content: content, CreatedAt: time.Now()}
	err := models.ORM.Create(comment).Error
	if err != nil {
		return comment, err
	}
	NewCommentEvent(comment, torrent)

	comment.Torrent = torrent
	cache.C.Delete(torrent.Identifier())

	return comment, nil
}
