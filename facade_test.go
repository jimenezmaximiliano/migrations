package migrations_test

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jimenezmaximiliano/migrations"
	"github.com/jimenezmaximiliano/migrations/models"
)

func TestRunningMigrations(test *testing.T) {
	db, err := sql.Open("mysql", "user:password@/db")
	require.Nil(test, err)

	_, err = db.Exec("DROP TABLE IF EXISTS gophers")
	require.Nil(test, err)

	_, err = db.Exec("DROP TABLE IF EXISTS migrations")
	require.Nil(test, err)

	result, err := migrations.RunMigrations(db, "./fixtures/create_and_insert")
	require.Nil(test, err)
	require.Len(test, result.GetAll(), 2)

	for _, currentMigration := range result.GetAll() {
		assert.Equal(test, models.StatusSuccessful, currentMigration.GetStatus())
	}
}

func TestRunningAMigrationWithTwoQueries(test *testing.T) {
	db, err := sql.Open("mysql", "user:password@/db?multiStatements=true")
	require.Nil(test, err)

	_, err = db.Exec("DROP TABLE IF EXISTS gophers")
	require.Nil(test, err)

	_, err = db.Exec("DROP TABLE IF EXISTS migrations")
	require.Nil(test, err)

	result, err := migrations.RunMigrations(db, "./fixtures/migration_with_two_queries")
	require.Nil(test, err)

	require.Len(test, result.GetAll(), 2)

	for _, currentMigration := range result.GetAll() {
		assert.Equal(test, models.StatusSuccessful, currentMigration.GetStatus())
	}
}

func TestRunningMigrationsWhenAllMigrationsHaveAlreadyRun(test *testing.T) {
	db, err := sql.Open("mysql", "user:password@/db")
	require.Nil(test, err)

	_, err = db.Exec("DROP TABLE IF EXISTS gophers")
	require.Nil(test, err)

	_, err = db.Exec("DROP TABLE IF EXISTS migrations")
	require.Nil(test, err)

	// First time
	_, err = migrations.RunMigrations(db, "./fixtures/create_and_insert")
	require.Nil(test, err)

	// Second time
	result, err := migrations.RunMigrations(db, "./fixtures/create_and_insert")
	require.Nil(test, err)

	require.Len(test, result.GetAll(), 0)
}

func TestRunningMigrationsStopsWhenAMigrationFails(test *testing.T) {
	db, err := sql.Open("mysql", "user:password@/db")
	require.Nil(test, err)

	_, err = db.Exec("DROP TABLE IF EXISTS gophers")
	require.Nil(test, err)

	_, err = db.Exec("DROP TABLE IF EXISTS migrations")
	require.Nil(test, err)

	result, err := migrations.RunMigrations(db, "./fixtures/create_insert_error")
	require.Nil(test, err)
	require.Len(test, result.GetAll(), 4)

	assert.Equal(test, models.StatusSuccessful, result.GetAll()[0].GetStatus())
	assert.Equal(test, models.StatusSuccessful, result.GetAll()[1].GetStatus())
	assert.Equal(test, models.StatusFailed, result.GetAll()[2].GetStatus())
	assert.Equal(test, models.StatusNotRun, result.GetAll()[3].GetStatus())
}
