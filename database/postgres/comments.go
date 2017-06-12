package postgres

import (
	"github.com/NyaaPantsu/nyaa/common"
	"github.com/NyaaPantsu/nyaa/model"
)

// InsertComment : Insert a comment
func (db *Database) InsertComment(comment *model.Comment) (err error) {
	_, err = db.getPrepared(queryInsertComment).Exec(comment.ID, comment.TorrentID, comment.Content, comment.CreatedAt)
	return
}

// GetCommentsWhere : Get comments on condition
func (db *Database) GetCommentsWhere(param *common.CommentParam) (comments []model.Comment, err error) {

	return
}

// DeleteCommentsWhere : Delete comments on condition
func (db *Database) DeleteCommentsWhere(param *common.CommentParam) (deleted uint32, err error) {

	return
}
