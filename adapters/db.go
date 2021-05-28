package adapters

import (
	"database/sql"
)

// DB is an adapter interface for database/sql.DB.
type DB interface {
	Ping() error
	Query(query string, args ...interface{}) (DBRows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// NewDBAdapter returns an implementation of DB.
func NewDBAdapter(db *sql.DB) DB {
	return dbAdapter{
		db: db,
	}
}

type dbAdapter struct {
	db *sql.DB
}

// Ping verifies a connection to the database is still alive,
// establishing a connection if necessary.
func (adapter dbAdapter) Ping() error {
	return adapter.db.Ping()
}

// Query executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (adapter dbAdapter) Query(query string, args ...interface{}) (DBRows, error) {
	return adapter.db.Query(query, args...)
}

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (adapter dbAdapter) Exec(query string, args ...interface{}) (sql.Result, error) {
	return adapter.db.Exec(query, args...)
}

// DBRows is the result of a query. Its cursor starts before the first row
// of the result set. Use Next to advance from row to row.
type DBRows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close() error
}

// Ensure sql.Rows implements DBRows.
var _ DBRows = &sql.Rows{}
