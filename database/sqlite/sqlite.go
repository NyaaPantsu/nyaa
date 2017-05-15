package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// queryEvent is a queued event to be executed in a pipeline to ensure that sqlite access is done from 1 goroutine
type queryEvent struct {
	query        string
	param        []interface{}
	handleResult func(*sql.Rows, error)
}

func New(param string) (db *Database, err error) {
	db = new(Database)
	db.conn, err = sql.Open("sqlite3", param)
	if err == nil {
		db.query = make(chan *queryEvent, 128)
	} else {
		db = nil
	}
	return
}

type Database struct {
	conn  *sql.DB
	query chan *queryEvent
}
