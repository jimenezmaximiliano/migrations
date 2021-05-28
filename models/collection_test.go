package models

import (
	"testing"
)

const anotherPath = "/tmp/another.sql"

func TestAddingAnItem(test *testing.T) {
	collection := Collection{}
	migration, err := NewMigration(validPath, validQuery, StatusNotRun)
	if err != nil {
		test.Fatal(err)
	}
	collection.Add(migration)

	if collection.IsEmpty() {
		test.Fail()
	}
}

func TestFindingAMigrationPath(test *testing.T) {
	collection := Collection{}
	migration, err := NewMigration(validPath, validQuery, StatusNotRun)
	if err != nil {
		test.Fatal(err)
	}
	collection.Add(migration)

	if !collection.ContainsMigrationPath(validPath) {
		test.Errorf("expected %s to be in the collection but wasn't found", validPath)
	}

	if collection.ContainsMigrationPath(anotherPath) {
		test.Errorf("expected %s to not be in the collection but was found", anotherPath)
	}
}

func TestGettingEveryMigrationWithTwoMigrations(test *testing.T) {
	collection := Collection{}
	migration, err := NewMigration(validPath, validQuery, StatusNotRun)
	if err != nil {
		test.Fatal(err)
	}
	collection.Add(migration)
	migration2, err := NewMigration(anotherPath, validQuery, StatusNotRun)
	if err != nil {
		test.Fatal(err)
	}
	collection.Add(migration2)
	migrations := collection.GetAll()

	if len(migrations) != 2 {
		test.Errorf("Expected 2 migrations but got %d", len(migrations))
	}
}

func TestGettingEveryMigrationOnAnEmptyCollection(test *testing.T) {
	collection := Collection{}
	migrations := collection.GetAll()

	if len(migrations) != 0 {
		test.Errorf("Expected 0 migrations but got %d", len(migrations))
	}
}

func TestGetAllSortsMigrations(test *testing.T) {
	collection := Collection{}
	migration, err := NewMigration("/tmp/obladi.sql", validQuery, StatusNotRun)
	if err != nil {
		test.Fatal(err)
	}
	collection.Add(migration)
	migration2, err := NewMigration("/tmp/oblada.sql", validQuery, StatusNotRun)
	if err != nil {
		test.Fatal(err)
	}
	collection.Add(migration2)
	migration3, err := NewMigration("/tmp/1.sql", validQuery, StatusNotRun)
	if err != nil {
		test.Fatal(err)
	}
	collection.Add(migration3)
	migrations := collection.GetAll()

	if migrations[0].GetAbsolutePath() != "/tmp/1.sql" ||
		migrations[1].GetAbsolutePath() != "/tmp/oblada.sql" ||
		migrations[2].GetAbsolutePath() != "/tmp/obladi.sql" {
		test.Error()
	}
}

func TestGetMigrationsToRunSortsMigrations(test *testing.T) {
	collection := Collection{}
	migration, err := NewMigration("/tmp/obladi.sql", validQuery, StatusNotRun)
	if err != nil {
		test.Fatal(err)
	}
	collection.Add(migration)
	migration2, err := NewMigration("/tmp/oblada.sql", validQuery, StatusNotRun)
	if err != nil {
		test.Fatal(err)
	}
	collection.Add(migration2)
	migration3, err := NewMigration("/tmp/1.sql", validQuery, StatusNotRun)
	if err != nil {
		test.Fatal(err)
	}
	collection.Add(migration3)
	migrations := collection.GetMigrationsToRun()

	if migrations[0].GetAbsolutePath() != "/tmp/1.sql" ||
		migrations[1].GetAbsolutePath() != "/tmp/oblada.sql" ||
		migrations[2].GetAbsolutePath() != "/tmp/obladi.sql" {
		test.Error()
	}
}

func TestGettingMigrationsToRun(test *testing.T) {
	collection := Collection{}
	migration, err := NewMigration("/tmp/1.sql", validQuery, StatusNotRun)
	if err != nil {
		test.Fatal(err)
	}
	collection.Add(migration)
	migration2, err := NewMigration("/tmp/2.sql", validQuery, StatusFailed)
	if err != nil {
		test.Fatal(err)
	}
	collection.Add(migration2)
	migration3, err := NewMigration("/tmp/3.sql", validQuery, StatusSuccessful)
	if err != nil {
		test.Fatal(err)
	}
	collection.Add(migration3)
	migration4, err := NewMigration("/tmp/3.sql", validQuery, StatusUnknown)
	if err != nil {
		test.Fatal(err)
	}
	collection.Add(migration4)
	migrations := collection.GetMigrationsToRun()

	if len(migrations) != 1 || migrations[0].GetAbsolutePath() != "/tmp/1.sql" {
		test.Error()
	}
}
