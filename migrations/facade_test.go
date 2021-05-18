package migrations

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jimenezmaximiliano/migrations/migrations/migration"
)

func TestRunningMigrations(test *testing.T) {
	db, err := sql.Open("mysql", "user:password@/db")
	if err != nil {
		test.Fatalf("failed to open db: %s", err.Error())
	}

	_, _ = db.Exec("DROP TABLE gophers")
	_, _ = db.Exec("DROP TABLE migrations")

	result, err := RunMigrations(db, "../fixtures/create_and_insert")
	if err != nil {
		test.Errorf("failed to run migrations: %s", err.Error())
	}

	numberOfMigrations := len(result.GetAll())
	if numberOfMigrations != 2 {
		test.Errorf("Expected 2 migrations to be run but got %d", numberOfMigrations)
	}

	for _, currentMigration := range result.GetAll() {
		if currentMigration.GetStatus() != migration.StatusSuccessful {
			test.Errorf("Migration %s failed", currentMigration.GetName())
		}
	}
}
