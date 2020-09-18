package main

import (
	"database/sql"
	migrate "github.com/jimenezmaximiliano/very-simple-migrations"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"path/filepath"
)

func main() {
	db, err := getDb()
	fail(err)

	migrationsPath, err := filepath.Abs("migrations/")
	fail(err)

	result, err := migrate.Run(db, migrationsPath)

	log.Println(result)
}

func fail(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func getDb() (*sql.DB, error) {
	dbPath, err := filepath.Abs("db.sqlite3")

	if err != nil {
		return nil, err
	}

	return sql.Open("sqlite3", dbPath)
}
