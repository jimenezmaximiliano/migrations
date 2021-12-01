package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const anotherPath = "/tmp/3_another.sql"

func TestAddingAnItem(test *testing.T) {
	collection := Collection{}
	migration, err := NewMigration(validPath, validQuery, StatusNotRun)
	require.Nil(test, err)

	err = collection.Add(migration)
	require.Nil(test, err)

	assert.False(test, collection.IsEmpty())
}

func TestFindingAMigrationPath(test *testing.T) {
	collection := Collection{}
	migration, err := NewMigration(validPath, validQuery, StatusNotRun)
	require.Nil(test, err)

	err = collection.Add(migration)
	require.Nil(test, err)

	assert.True(test, collection.ContainsMigrationPath(validPath))
	assert.False(test, collection.ContainsMigrationPath(anotherPath))
}

func TestGettingEveryMigrationWithTwoMigrations(test *testing.T) {
	collection := Collection{}
	migration, err := NewMigration(validPath, validQuery, StatusNotRun)
	require.Nil(test, err)

	err = collection.Add(migration)
	require.Nil(test, err)

	migration2, err := NewMigration(anotherPath, validQuery, StatusNotRun)
	require.Nil(test, err)

	err = collection.Add(migration2)
	require.Nil(test, err)

	migrations := collection.GetAll()

	assert.Len(test, migrations, 2)
}

func TestAddingTwoMigrationsWithTheSameOrderFails(test *testing.T) {
	collection := Collection{}
	migration, err := NewMigration("/tmp/1_a.sql", validQuery, StatusNotRun)
	require.Nil(test, err)

	err = collection.Add(migration)
	require.Nil(test, err)

	migration2, err := NewMigration("/tmp/1_b.sql", validQuery, StatusNotRun)
	require.Nil(test, err)

	err = collection.Add(migration2)
	assert.NotNil(test, err)
}

func TestGettingEveryMigrationOnAnEmptyCollection(test *testing.T) {
	collection := Collection{}
	migrations := collection.GetAll()

	assert.Len(test, migrations, 0)
}

func TestGetAllSortsMigrations(test *testing.T) {
	collection := Collection{}
	migration, err := NewMigration("/tmp/1_obladi.sql", validQuery, StatusNotRun)
	require.Nil(test, err)

	err = collection.Add(migration)
	require.Nil(test, err)

	migration2, err := NewMigration("/tmp/2_oblada.sql", validQuery, StatusNotRun)
	require.Nil(test, err)

	err = collection.Add(migration2)
	require.Nil(test, err)

	migration3, err := NewMigration("/tmp/3_desmond.sql", validQuery, StatusNotRun)
	require.Nil(test, err)

	err = collection.Add(migration3)
	require.Nil(test, err)

	migrations := collection.GetAll()

	assert.Equal(test, "/tmp/1_obladi.sql", migrations[0].GetAbsolutePath())
	assert.Equal(test, "/tmp/2_oblada.sql", migrations[1].GetAbsolutePath())
	assert.Equal(test, "/tmp/3_desmond.sql", migrations[2].GetAbsolutePath())
}

func TestGetMigrationsToRunSortsMigrations(test *testing.T) {
	collection := Collection{}
	migration, err := NewMigration("/tmp/1_obladi.sql", validQuery, StatusNotRun)
	require.Nil(test, err)

	err = collection.Add(migration)
	require.Nil(test, err)

	migration2, err := NewMigration("/tmp/2_oblada.sql", validQuery, StatusNotRun)
	require.Nil(test, err)

	err = collection.Add(migration2)
	require.Nil(test, err)

	migration3, err := NewMigration("/tmp/3_desmond.sql", validQuery, StatusNotRun)
	require.Nil(test, err)

	err = collection.Add(migration3)
	require.Nil(test, err)

	migrations := collection.GetMigrationsToRun()

	assert.Equal(test, "/tmp/1_obladi.sql", migrations[0].GetAbsolutePath())
	assert.Equal(test, "/tmp/2_oblada.sql", migrations[1].GetAbsolutePath())
	assert.Equal(test, "/tmp/3_desmond.sql", migrations[2].GetAbsolutePath())
}

func TestGettingMigrationsToRun(test *testing.T) {
	collection := Collection{}
	migration, err := NewMigration("/tmp/1_a.sql", validQuery, StatusNotRun)
	require.Nil(test, err)

	err = collection.Add(migration)
	require.Nil(test, err)

	migration2, err := NewMigration("/tmp/2_b.sql", validQuery, StatusFailed)
	require.Nil(test, err)

	err = collection.Add(migration2)
	require.Nil(test, err)

	migration3, err := NewMigration("/tmp/3_c.sql", validQuery, StatusSuccessful)
	require.Nil(test, err)

	err = collection.Add(migration3)
	require.Nil(test, err)

	migration4, err := NewMigration("/tmp/4_d.sql", validQuery, StatusUnknown)
	require.Nil(test, err)

	err = collection.Add(migration4)
	require.Nil(test, err)

	migrations := collection.GetMigrationsToRun()

	assert.Len(test, migrations, 1)
	assert.Equal(test, "/tmp/1_a.sql", migrations[0].GetAbsolutePath())
}
