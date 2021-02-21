package adapters

import (
	"database/sql"
	"os"
)

// SQLResult is used to mock sql.Result
type SQLResult interface {
	sql.Result
}

// File is used to mock os.FileInfo
type File interface {
	os.FileInfo
}
