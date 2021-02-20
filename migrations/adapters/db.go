package adapters

import (
	"database/sql"
	stdSQL "database/sql"
)

type DB interface {
	Ping() error
	Query(query string, args ...interface{}) (DBRows, error)
	Exec(query string, args ...interface{}) (stdSQL.Result, error)
}

func NewDBAdapter(db *sql.DB) DB {
	return dbAdapter{
		db: db,
	}
}

type dbAdapter struct {
	db *sql.DB
}

func (adapter dbAdapter) Ping() error {
	return adapter.db.Ping()
}

func (adapter dbAdapter) Query(query string, args ...interface{}) (DBRows, error) {
	return adapter.db.Query(query, args...)
}

func (adapter dbAdapter) Exec(query string, args ...interface{}) (stdSQL.Result, error) {
	return adapter.db.Exec(query, args...)
}

type DBRows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close() error
}

// Ensure sql.Rows implements DBRows
var _ DBRows = &sql.Rows{}
