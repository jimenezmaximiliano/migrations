package migrations

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jimenezmaximiliano/migrations/models"
)

func TestRunningMigrations(test *testing.T) {
	db, err := sql.Open("mysql", "user:password@/db")
	if err != nil {
		test.Fatalf("failed to open db: %s", err.Error())
	}

	_, _ = db.Exec("DROP TABLE gophers")
	_, _ = db.Exec("DROP TABLE migrations")

	result, err := RunMigrations(db, "./fixtures/create_and_insert")
	if err != nil {
		test.Errorf("failed to run migrations: %s", err.Error())
	}

	numberOfMigrations := len(result.GetAll())
	if numberOfMigrations != 2 {
		test.Errorf("Expected 2 migrations to be run but got %d", numberOfMigrations)
	}

	for _, currentMigration := range result.GetAll() {
		if currentMigration.GetStatus() != models.StatusSuccessful {
			test.Errorf("Migration %s failed", currentMigration.GetName())
		}
	}
}

func TestRunningAMigrationWithTwoQueries(test *testing.T) {
	db, err := sql.Open("mysql", "user:password@/db?multiStatements=true")
	if err != nil {
		test.Fatalf("failed to open db: %s", err.Error())
	}

	_, _ = db.Exec("DROP TABLE gophers")
	_, _ = db.Exec("DROP TABLE migrations")

	result, err := RunMigrations(db, "./fixtures/migration_with_two_queries")
	if err != nil {
		test.Errorf("failed to run migrations: %s", err.Error())
	}

	numberOfMigrations := len(result.GetAll())
	if numberOfMigrations != 2 {
		test.Errorf("Expected 2 migrations to be run but got %d", numberOfMigrations)
	}

	for _, currentMigration := range result.GetAll() {
		if currentMigration.GetStatus() != models.StatusSuccessful {
			test.Errorf("Migration %s failed", currentMigration.GetName())
		}
	}
}

func TestRunningMigrationsWhenAllMigrationsHaveAlreadyRun(test *testing.T) {
	db, err := sql.Open("mysql", "user:password@/db")
	if err != nil {
		test.Fatalf("failed to open db: %s", err.Error())
	}

	_, _ = db.Exec("DROP TABLE gophers")
	_, _ = db.Exec("DROP TABLE migrations")

	// First time
	_, err = RunMigrations(db, "./fixtures/create_and_insert")
	if err != nil {
		test.Errorf("failed to run migrations: %s", err.Error())
	}

	// Second time
	result, err := RunMigrations(db, "./fixtures/create_and_insert")
	if err != nil {
		test.Errorf("failed to run migrations: %s", err.Error())
	}

	numberOfMigrations := len(result.GetAll())
	if numberOfMigrations != 0 {
		test.Errorf("Expected 0 migrations to be run but got %d", numberOfMigrations)
	}
}

func TestRunningMigrationsStopsWhenAMigrationFails(test *testing.T) {
	db, err := sql.Open("mysql", "user:password@/db")
	if err != nil {
		test.Fatalf("failed to open db: %s", err.Error())
	}

	_, _ = db.Exec("DROP TABLE gophers")
	_, _ = db.Exec("DROP TABLE migrations")

	result, err := RunMigrations(db, "./fixtures/create_insert_error")
	if err != nil {
		test.Errorf("failed to run migrations: %s", err.Error())
	}

	all := result.GetAll()

	if len(all) != 4 {
		test.Errorf("exptected 4 migrations but got %d", len(all))
	}

	if all[0].GetStatus() != models.StatusSuccessful ||
		all[1].GetStatus() != models.StatusSuccessful ||
		all[2].GetStatus() != models.StatusFailed ||
		all[3].GetStatus() != models.StatusNotRun {
		test.Errorf("invalid status on migration results")
	}
}