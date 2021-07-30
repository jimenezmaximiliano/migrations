package models

import (
	"testing"
)

const anotherPath = "/tmp/3_another.sql"

func TestAddingAnItem(test *testing.T) {
	collection := Collection{}
	migration, err := NewMigration(validPath, validQuery, StatusNotRun)
	if err != nil {
		test.Fatal(err)
	}
	err = collection.Add(migration)
	if err != nil {
		test.Fatal(err)
	}

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

	err = collection.Add(migration)
	if err != nil {
		test.Fatal(err)
	}

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

	err = collection.Add(migration)
	if err != nil {
		test.Fatal(err)
	}

	migration2, err := NewMigration(anotherPath, validQuery, StatusNotRun)
	if err != nil {
		test.Fatal(err)
	}

	err = collection.Add(migration2)
	if err != nil {
		test.Fatal(err)
	}

	migrations := collection.GetAll()

	if len(migrations) != 2 {
		test.Errorf("Expected 2 migrations but got %d", len(migrations))
	}
}

func TestAddingTwoMigrationsWithTheSameOrderFails(test *testing.T) {
	collection := Collection{}
	migration, err := NewMigration("/tmp/1_a.sql", validQuery, StatusNotRun)
	if err != nil {
		test.Fatal(err)
	}

	err = collection.Add(migration)
	if err != nil {
		test.Fatal(err)
	}

	migration2, err := NewMigration("/tmp/1_b.sql", validQuery, StatusNotRun)
	if err != nil {
		test.Fatal(err)
	}

	err = collection.Add(migration2)
	if err == nil {
		test.Fail()
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
	migration, err := NewMigration("/tmp/1_obladi.sql", validQuery, StatusNotRun)
	if err != nil {
		test.Fatal(err)
	}

	err = collection.Add(migration)
	if err != nil {
		test.Fatal(err)
	}

	migration2, err := NewMigration("/tmp/2_oblada.sql", validQuery, StatusNotRun)
	if err != nil {
		test.Fatal(err)
	}

	err = collection.Add(migration2)
	if err != nil {
		test.Fatal(err)
	}

	migration3, err := NewMigration("/tmp/3_desmond.sql", validQuery, StatusNotRun)
	if err != nil {
		test.Fatal(err)
	}

	err = collection.Add(migration3)
	if err != nil {
		test.Fatal(err)
	}

	migrations := collection.GetAll()

	if migrations[0].GetAbsolutePath() != "/tmp/1_obladi.sql" ||
		migrations[1].GetAbsolutePath() != "/tmp/2_oblada.sql" ||
		migrations[2].GetAbsolutePath() != "/tmp/3_desmond.sql" {
		test.Error()
	}
}

func TestGetMigrationsToRunSortsMigrations(test *testing.T) {
	collection := Collection{}
	migration, err := NewMigration("/tmp/1_obladi.sql", validQuery, StatusNotRun)
	if err != nil {
		test.Fatal(err)
	}

	err = collection.Add(migration)
	if err != nil {
		test.Fatal(err)
	}

	migration2, err := NewMigration("/tmp/2_oblada.sql", validQuery, StatusNotRun)
	if err != nil {
		test.Fatal(err)
	}

	err = collection.Add(migration2)
	if err != nil {
		test.Fatal(err)
	}

	migration3, err := NewMigration("/tmp/3_desmond.sql", validQuery, StatusNotRun)
	if err != nil {
		test.Fatal(err)
	}

	err = collection.Add(migration3)
	if err != nil {
		test.Fatal(err)
	}

	migrations := collection.GetMigrationsToRun()

	if migrations[0].GetAbsolutePath() != "/tmp/1_obladi.sql" ||
		migrations[1].GetAbsolutePath() != "/tmp/2_oblada.sql" ||
		migrations[2].GetAbsolutePath() != "/tmp/3_desmond.sql" {
		test.Error()
	}
}

func TestGettingMigrationsToRun(test *testing.T) {
	collection := Collection{}
	migration, err := NewMigration("/tmp/1_a.sql", validQuery, StatusNotRun)
	if err != nil {
		test.Fatal(err)
	}

	err = collection.Add(migration)
	if err != nil {
		test.Fatal(err)
	}

	migration2, err := NewMigration("/tmp/2_b.sql", validQuery, StatusFailed)
	if err != nil {
		test.Fatal(err)
	}

	err = collection.Add(migration2)
	if err != nil {
		test.Fatal(err)
	}

	migration3, err := NewMigration("/tmp/3_c.sql", validQuery, StatusSuccessful)
	if err != nil {
		test.Fatal(err)
	}

	err = collection.Add(migration3)
	if err != nil {
		test.Fatal(err)
	}

	migration4, err := NewMigration("/tmp/4_d.sql", validQuery, StatusUnknown)
	if err != nil {
		test.Fatal(err)
	}

	err = collection.Add(migration4)
	if err != nil {
		test.Fatal(err)
	}

	migrations := collection.GetMigrationsToRun()

	if len(migrations) != 1 || migrations[0].GetAbsolutePath() != "/tmp/1_a.sql" {
		test.Error()
	}
}
