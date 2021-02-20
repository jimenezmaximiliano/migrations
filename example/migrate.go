package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jimenezmaximiliano/migrations/migrations"
)

func main() {
	migrations.RunMigrationsCommand(func() (*sql.DB, error) {
		return sql.Open("mysql", "user:password@/db")
	})
}
