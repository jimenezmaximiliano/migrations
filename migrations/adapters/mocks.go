package adapters

import "database/sql"

// SQLResult is used to mocke sql.Result locally
type SQLResult interface {
	sql.Result
}
