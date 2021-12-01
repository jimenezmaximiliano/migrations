package models

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const validOrder = uint64(1627676757857350000)
const validName = "1627676757857350000_createTableGophers.sql"
const validPath = "/tmp/migrations/" + validName
const validQuery = "CREATE TABLE gophers;"

func TestMigrationDefaultValues(test *testing.T) {
	migration := migration{}

	path := migration.GetAbsolutePath()
	assert.Equal(test, "", path)

	name := migration.GetName()
	assert.Equal(test, "", name)

	query := migration.GetQuery()
	assert.Equal(test, "", query)

	status := migration.GetStatus()
	assert.Equal(test, StatusUnknown, status)

	order := migration.GetOrder()
	assert.Equal(test, uint64(0), order)
}

func TestMigrationConstruction(test *testing.T) {

	const status = StatusNotRun
	migration, err := NewMigration(validPath, validQuery, status)
	require.Nil(test, err)

	assert.Equal(test, validPath, migration.GetAbsolutePath())
	assert.Equal(test, validName, migration.GetName())
	assert.Equal(test, validQuery, migration.GetQuery())
	assert.Equal(test, status, migration.GetStatus())
	assert.Equal(test, validOrder, migration.GetOrder())
}

func TestMigrationConstructionFailsWithAnInvalidOrder(test *testing.T) {
	_, err := NewMigration("/tmp/maxi.sql", validQuery, StatusUnknown)

	assert.NotNil(test, err)
}

func TestMigrationConstructionFailsWithAnInvalidStatus(test *testing.T) {
	_, err := NewMigration(validPath, validQuery, -2)

	assert.NotNil(test, err)
}

func TestMigrationShouldBeRun(test *testing.T) {
	migration, err := NewMigration(validPath, validQuery, StatusNotRun)
	require.Nil(test, err)

	assert.True(test, migration.ShouldBeRun())

	notRunnableStatuses := []int8{StatusFailed, StatusSuccessful, StatusUnknown}
	for _, status := range notRunnableStatuses {
		migration, err = NewMigration(validPath, validQuery, status)
		require.Nil(test, err)

		assert.False(test, migration.ShouldBeRun())
	}
}

func TestStatusHelpers(test *testing.T) {
	migration, err := NewMigration(validPath, validQuery, StatusSuccessful)
	require.Nil(test, err)

	assert.True(test, migration.WasSuccessful())

	migration, err = NewMigration(validPath, validQuery, StatusFailed)
	assert.Nil(test, err)
	assert.True(test, migration.HasFailed())
}

func TestChangingTheMigrationsStatusToFailed(test *testing.T) {
	migration, err := NewMigration(validPath, validQuery, StatusNotRun)
	require.Nil(test, err)

	failedMigration := migration.NewAsFailed(errors.New("oops"))

	assert.Equal(test, migration.GetName(), failedMigration.GetName())
	assert.Equal(test, migration.GetQuery(), failedMigration.GetQuery())
	assert.Equal(test, migration.GetOrder(), failedMigration.GetOrder())
	assert.Equal(test, migration.GetAbsolutePath(), failedMigration.GetAbsolutePath())
	assert.True(test, failedMigration.HasFailed())
}

func TestChangingTheMigrationsStatusToSuccessful(test *testing.T) {
	migration, err := NewMigration(validPath, validQuery, StatusNotRun)
	require.Nil(test, err)

	successfulMigration := migration.NewAsSuccessful()

	assert.Equal(test, migration.GetName(), successfulMigration.GetName())
	assert.Equal(test, migration.GetQuery(), successfulMigration.GetQuery())
	assert.Equal(test, migration.GetOrder(), successfulMigration.GetOrder())
	assert.Equal(test, migration.GetAbsolutePath(), successfulMigration.GetAbsolutePath())
	assert.True(test, successfulMigration.WasSuccessful())
}

func TestShouldBeRunFirst(test *testing.T) {
	migration2020, _ := NewMigration("/2020_a.sql", "", StatusNotRun)
	migration2021, _ := NewMigration("/2021_b.sql", "", StatusNotRun)

	assert.True(test, migration2020.ShouldBeRunFirst(migration2021))
	assert.False(test, migration2021.ShouldBeRunFirst(migration2020))
}
