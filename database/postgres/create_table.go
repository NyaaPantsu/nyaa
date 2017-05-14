package postgres

import (
	"database/sql"
	"fmt"
	"strings"
)

// sql query
type sqlQuery struct {
	query  string
	params []interface{}
}

func (q *sqlQuery) Exec(conn *sql.DB) (err error) {
	_, err = conn.Exec(q.query, q.params...)
	return
}

func (q *sqlQuery) QueryRow(conn *sql.DB, visitor func(*sql.Row) error) (err error) {
	err = visitor(conn.QueryRow(q.query, q.params))
	return
}

func (q *sqlQuery) Query(conn *sql.DB, visitor func(*sql.Rows) error) (err error) {
	var rows *sql.Rows
	rows, err = conn.Query(q.query, q.params...)
	if err == sql.ErrNoRows {
		err = nil
	} else if err == nil {
		err = visitor(rows)
		rows.Close()
	}
	return
}

// make a createQuery that creates an index for column on table
func createIndex(table, column string) sqlQuery {
	return sqlQuery{
		query: fmt.Sprintf("CREATE INDEX IF NOT EXISTS ON %s %s ", table, column),
	}
}

// make a createQuery that creates a trigraph index on a table for multiple columns
func createTrigraph(table string, columns ...string) sqlQuery {
	return sqlQuery{
		query: fmt.Sprintf("CREATE INDEX IF NOT EXISTS ON %s USING gin(%s gin_trgm_ops)", table, strings.Join(columns, ", ")),
	}
}

// defines table creation info
type createTable struct {
	name       string       // Table's name
	columns    tableColumns // Table's columns
	preCreate  []sqlQuery   // Queries to run before table is ensrued
	postCreate []sqlQuery   // Queries to run after table is ensured
}

func (t createTable) String() string {
	return fmt.Sprintf("CREATE TABLE %s IF NOT EXISTS ( %s )", t.name, t.columns)
}

func (t createTable) Exec(conn *sql.DB) (err error) {
	// pre queries
	for idx := range t.preCreate {
		err = t.preCreate[idx].Exec(conn)
		if err != nil {
			return
		}
	}
	// table definition
	_, err = conn.Exec(t.String())
	if err != nil {
		return
	}
	// post queries
	for idx := range t.postCreate {
		err = t.postCreate[idx].Exec(conn)
		if err != nil {
			return
		}
	}
	return
}

// tableColumns is a list of columns for a table to be created
type tableColumns []string

func (def tableColumns) String() string {
	return strings.Join(def, ", ")
}
