package main

import (
	"database/sql"
	"flag"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	m "github.com/jimenezmaximiliano/migrations/migrations"
)

func main() {
	connectionString := flag.String("connection", "", "")
	migrationsPath := flag.String("path", "", "")

	flag.Parse()

	db, err := sql.Open("mysql", *connectionString)

	if err != nil {
		fmt.Println("Cannot open DB")
		panic(err)
	}

	migrations, _ := m.RunMigrations(db, *migrationsPath)

	m.DisplayResults(migrations)
}
