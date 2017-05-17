package postgres

import (
	"github.com/NyaaPantsu/nyaa/common"
	"github.com/NyaaPantsu/nyaa/model"
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
