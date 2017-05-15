package postgres

import (
	"github.com/ewhal/nyaa/common"
	"github.com/ewhal/nyaa/model"
)

func (db *Database) InsertComment(comment *model.Comment) (err error) {
	_, err = db.getPrepared(queryInsertComment).Exec(comment.ID, comment.TorrentID, comment.Content, comment.CreatedAt)
	return
}

func (db *Database) GetCommentsWhere(param *common.CommentParam) (comments []model.Comment, err error) {

	return
}

func (db *Database) DeleteCommentsWhere(param *common.CommentParam) (deleted uint32, err error) {

	return
}
