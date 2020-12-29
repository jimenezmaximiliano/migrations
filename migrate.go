package main

import (
	"database/sql"
	"flag"
	"fmt"

	m "github.com/jimenezmaximiliano/migrations/migrations"

	_ "github.com/go-sql-driver/mysql"
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

	migrations, err := m.RunMigrations(db, *migrationsPath)

	for _, migration := range migrations.GetAll() {
		if migration.WasSuccessful() {
			fmt.Printf("[OK] %s\n", migration.GetName())
			continue
		}
		fmt.Printf("[KO] %s\n", migration.GetName())
	}

	if err != nil {
		fmt.Println("Failed to run migrations")
		fmt.Printf("%v\n", err)
	}
}
