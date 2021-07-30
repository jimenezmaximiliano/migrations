package models

import (
	"errors"
	"testing"
)

const validOrder = 1627676757857350000
const validName = "1627676757857350000_createTableGophers.sql"
const validPath = "/tmp/migrations/" + validName
const validQuery = "CREATE TABLE gophers;"

func TestMigrationDefaultValues(test *testing.T) {
	migration := migration{}

	path := migration.GetAbsolutePath()
	const expectedPath = ""
	if path != expectedPath {
		test.Errorf("expected path %s but got %s", expectedPath, path)
	}

	name := migration.GetName()
	const expectedName = ""
	if name != expectedName {
		test.Errorf("expected name %s but got %s", expectedName, name)
	}

	query := migration.GetQuery()
	const expectedQuery = ""
	if query != expectedQuery {
		test.Errorf("expected query %s but got %s", expectedQuery, query)
	}

	status := migration.GetStatus()
	const expectedStatus = StatusUnknown
	if status != expectedStatus {
		test.Errorf("expected status %d but got %d", expectedStatus, status)
	}

	order := migration.GetOrder()
	const expectedOrder = 0
	if status != expectedStatus {
		test.Errorf("expected order %d but got %d", expectedOrder, order)
	}
}

func TestMigrationConstruction(test *testing.T) {

	const status = StatusNotRun
	migration, err := NewMigration(validPath, validQuery, status)
	if err != nil {
		test.Fatal(err)
	}

	if migration.GetAbsolutePath() != validPath {
		test.Errorf("expected absolute path %s but got %s", validPath, migration.GetAbsolutePath())
	}

	if migration.GetName() != validName {
		test.Errorf("expected name %s but got %s", validName, migration.GetName())
	}

	if migration.GetQuery() != validQuery {
		test.Errorf("expected query %s but got %s", validQuery, migration.GetQuery())
	}

	if migration.GetStatus() != status {
		test.Errorf("expected status %d but got %d", status, migration.GetStatus())
	}

	if migration.GetOrder() != validOrder {
		test.Errorf("expected order %d but got %d", validOrder, migration.GetOrder())
	}
}

func TestMigrationConstructionFailsWithAnInvalidOrder(test *testing.T) {
	_, err := NewMigration("/tmp/maxi.sql", validQuery, StatusUnknown)
	if err == nil {
		test.Fatal(err)
	}
}

func TestMigrationConstructionFailsWithAnInvalidStatus(test *testing.T) {
	_, err := NewMigration(validPath, validQuery, -2)
	if err == nil {
		test.Fatal(err)
	}
}

func TestMigrationShouldBeRun(test *testing.T) {
	migration, err := NewMigration(validPath, validQuery, StatusNotRun)
	if err != nil {
		test.Fatal(err)
	}
	if migration.ShouldBeRun() == false {
		test.Fail()
	}

	notRunnableStatuses := []int8{StatusFailed, StatusSuccessful, StatusUnknown}
	for _, status := range notRunnableStatuses {
		migration, err = NewMigration(validPath, validQuery, status)
		if err != nil {
			test.Fatal(err)
		}
		if migration.ShouldBeRun() == true {
			test.Fail()
		}
	}
}

func TestStatusHelpers(test *testing.T) {
	migration, err := NewMigration(validPath, validQuery, StatusSuccessful)
	if err != nil {
		test.Fatal(err)
	}
	if !migration.WasSuccessful() {
		test.Fail()
	}

	migration, err = NewMigration(validPath, validQuery, StatusFailed)
	if err != nil {
		test.Fatal(err)
	}
	if !migration.HasFailed() {
		test.Fail()
	}
}

func TestChangingTheMigrationsStatusToFailed(test *testing.T) {
	migration, err := NewMigration(validPath, validQuery, StatusNotRun)
	if err != nil {
		test.Fatal(err)
	}

	failedMigration := migration.NewAsFailed(errors.New("oops"))
	if migration.GetName() != failedMigration.GetName() ||
		migration.GetQuery() != failedMigration.GetQuery() ||
		migration.GetOrder() != failedMigration.GetOrder() ||
		migration.GetAbsolutePath() != failedMigration.GetAbsolutePath() ||
		!failedMigration.HasFailed() {
		test.Fail()
	}
}

func TestChangingTheMigrationsStatusToSuccessful(test *testing.T) {
	migration, err := NewMigration(validPath, validQuery, StatusNotRun)
	if err != nil {
		test.Fatal(err)
	}

	successfulMigration := migration.NewAsSuccessful()
	if migration.GetName() != successfulMigration.GetName() ||
		migration.GetQuery() != successfulMigration.GetQuery() ||
		migration.GetOrder() != successfulMigration.GetOrder() ||
		migration.GetAbsolutePath() != successfulMigration.GetAbsolutePath() ||
		!successfulMigration.WasSuccessful() {
		test.Fail()
	}
}

func TestShouldBeRunFirst(test *testing.T) {
	migration2020, _ := NewMigration("/2020_a.sql", "", StatusNotRun)
	migration2021, _ := NewMigration("/2021_b.sql", "", StatusNotRun)

	if migration2020.ShouldBeRunFirst(migration2021) == false {
		test.Fail()
	}

	if migration2021.ShouldBeRunFirst(migration2020) == true {
		test.Fail()
	}
}
