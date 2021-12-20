package migrations

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jimenezmaximiliano/migrations/models"
)

func TestRunningMigrations(test *testing.T) {
	db, err := sql.Open("mysql", "user:password@/db")
	require.Nil(test, err)

	_, _ = db.Exec("DROP TABLE gophers")
	_, _ = db.Exec("DROP TABLE migrations")

	result, err := RunMigrations(db, "./fixtures/create_and_insert")
	require.Nil(test, err)
	require.Len(test, result.GetAll(), 2)

	for _, currentMigration := range result.GetAll() {
		assert.Equal(test, models.StatusSuccessful, currentMigration.GetStatus())
	}
}

func TestRunningAMigrationWithTwoQueries(test *testing.T) {
	db, err := sql.Open("mysql", "user:password@/db?multiStatements=true")
	require.Nil(test, err)

	_, _ = db.Exec("DROP TABLE gophers")
	_, _ = db.Exec("DROP TABLE migrations")

	result, err := RunMigrations(db, "./fixtures/migration_with_two_queries")
	require.Nil(test, err)

	require.Len(test, result.GetAll(), 2)

	for _, currentMigration := range result.GetAll() {
		assert.Equal(test, models.StatusSuccessful, currentMigration.GetStatus())
	}
}

func TestRunningMigrationsWhenAllMigrationsHaveAlreadyRun(test *testing.T) {
	db, err := sql.Open("mysql", "user:password@/db")
	require.Nil(test, err)

	_, _ = db.Exec("DROP TABLE gophers")
	_, _ = db.Exec("DROP TABLE migrations")

	// First time
	_, err = RunMigrations(db, "./fixtures/create_and_insert")
	require.Nil(test, err)

	// Second time
	result, err := RunMigrations(db, "./fixtures/create_and_insert")
	require.Nil(test, err)

	require.Len(test, result.GetAll(), 0)
}

func TestRunningMigrationsStopsWhenAMigrationFails(test *testing.T) {
	db, err := sql.Open("mysql", "user:password@/db")
	require.Nil(test, err)

	_, _ = db.Exec("DROP TABLE gophers")
	_, _ = db.Exec("DROP TABLE migrations")

	result, err := RunMigrations(db, "./fixtures/create_insert_error")
	require.Nil(test, err)
	require.Len(test, result.GetAll(), 4)

	assert.Equal(test, models.StatusSuccessful, result.GetAll()[0].GetStatus())
	assert.Equal(test, models.StatusSuccessful, result.GetAll()[1].GetStatus())
	assert.Equal(test, models.StatusFailed, result.GetAll()[2].GetStatus())
	assert.Equal(test, models.StatusNotRun, result.GetAll()[3].GetStatus())
}