package postgres

import (
	"database/sql"

	"github.com/ewhal/nyaa/util/log"
	_ "github.com/lib/pq"
)

func New(param string) (db *Database, err error) {
	db = new(Database)
	db.conn, err = sql.Open("postgres", param)
	if err != nil {
		db = nil
	}
	return
}

type Database struct {
	conn     *sql.DB
	prepared map[string]*sql.Stmt
}

func (db *Database) getPrepared(name string) *sql.Stmt {
	return db.prepared[name]
}

func (db *Database) Query(q string, param ...interface{}) (*sql.Rows, error) {
	return db.conn.Query(q, param)
}

func (db *Database) Init() (err error) {

	// ensure tables
	for idx := range tables {
		log.Debugf("ensure table %s", tables[idx].name)
		err = tables[idx].Exec(db.conn)
		if err != nil {
			log.Errorf("Failed to ensure table %s: %s", tables[idx].name, err.Error())
			return
		}
	}

	// generate prepared statements
	for k := range statements {
		var stmt *sql.Stmt
		stmt, err = db.conn.Prepare(statements[k])
		if err != nil {
			log.Errorf("failed to build prepared statement %s: %s", k, err.Error())
			return
		}
		db.prepared[k] = stmt
	}
	return
}

// execute prepared statement with arguments, visit result and autoclose rows after done visiting
func (db *Database) queryWithPrepared(name string, visit func(*sql.Rows) error, params ...interface{}) (err error) {
	var rows *sql.Rows
	rows, err = db.getPrepared(name).Query(params)
	if err == sql.ErrNoRows {
		err = nil
	} else if err == nil {
		err = visit(rows)
		rows.Close()
	}
	return
}

// execute prepared statement with arguments, visit single row
func (db *Database) queryRowWithPrepared(name string, visit func(*sql.Row) error, params ...interface{}) (err error) {
	err = visit(db.getPrepared(name).QueryRow(params))
	return
}

// execute a query by name and return how many rows were affected
func (db *Database) execQuery(name string, p ...interface{}) (affected uint32, err error) {
	var result sql.Result
	result, err = db.getPrepared(name).Exec(p)
	if err == nil {
		var d int64
		d, err = result.RowsAffected()
		if err == nil {
			affected = uint32(d)
		}
	}
	return
}
