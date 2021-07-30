package services

import (
	"fmt"
	"testing"

	"github.com/jimenezmaximiliano/migrations/mocks"
	"github.com/jimenezmaximiliano/migrations/models"
)

func TestRunningMigrationsFailsIfTheDBConnectionDoesNotWork(test *testing.T) {
	fetcher := &mocks.Fetcher{}
	db := &mocks.DBRepository{}
	db.On("Ping").Return(fmt.Errorf("db connection error"))
	service := NewRunnerService(fetcher, db, "/tmp")

	_, err := service.RunMigrations()

	if err == nil {
		test.Fail()
	}
}

func TestRunningMigrationsFailsIfAMigrationsTableCannotBeCreated(test *testing.T) {
	fetcher := &mocks.Fetcher{}
	db := &mocks.DBRepository{}
	db.On("Ping").Return(nil)
	db.On("CreateMigrationsTableIfNeeded").Return(fmt.Errorf("cannot create table"))
	service := NewRunnerService(fetcher, db, "/tmp")

	_, err := service.RunMigrations()

	if err == nil {
		test.Fail()
	}
}

func TestRunningMigrationsFailsIfItCannotFetchMigrationsFromFilesOrTheDB(test *testing.T) {
	db := &mocks.DBRepository{}
	db.On("Ping").Return(nil)
	db.On("CreateMigrationsTableIfNeeded").Return(nil)
	fetcher := &mocks.Fetcher{}
	fetcher.On("GetMigrations", "/tmp/").Return(models.Collection{}, fmt.Errorf("cannot fetch migrations"))
	service := NewRunnerService(fetcher, db, "/tmp")

	_, err := service.RunMigrations()

	if err == nil {
		test.Fail()
	}
}

func TestRunningMigrationsDoesNotFailIfThereAreNoMigratiosToRun(test *testing.T) {
	db := &mocks.DBRepository{}
	db.On("Ping").Return(nil)
	db.On("CreateMigrationsTableIfNeeded").Return(nil)
	fetcher := &mocks.Fetcher{}
	fetcher.On("GetMigrations", "/tmp/").Return(models.Collection{}, nil)
	service := NewRunnerService(fetcher, db, "/tmp")

	_, err := service.RunMigrations()

	if err != nil {
		test.Fail()
	}
}

func TestRunningAMigrationSuccessfully(test *testing.T) {
	db := &mocks.DBRepository{}
	db.On("Ping").Return(nil)
	db.On("CreateMigrationsTableIfNeeded").Return(nil)
	db.On("RunMigrationQuery", "SELECT 1").Return(nil)
	db.On("RegisterRunMigration", "1.sql").Return(nil)
	fetcher := &mocks.Fetcher{}
	collection := models.Collection{}
	migration, _ := models.NewMigration("/tmp/1.sql", "SELECT 1", models.StatusNotRun)
	err := collection.Add(migration)
	if err != nil {
		test.Fatal(err)
	}
	fetcher.On("GetMigrations", "/tmp/").Return(collection, nil)
	service := NewRunnerService(fetcher, db, "/tmp")

	result, err := service.RunMigrations()
	if err != nil {
		test.Fail()
	}

	if result.GetAll()[0].GetAbsolutePath() != "/tmp/1.sql" || result.GetAll()[0].HasFailed() {
		test.Fail()
	}
}

func TestRunningAMigrationThatFails(test *testing.T) {
	db := &mocks.DBRepository{}
	db.On("Ping").Return(nil)
	db.On("CreateMigrationsTableIfNeeded").Return(nil)
	db.On("RunMigrationQuery", "SELECT 1").Return(fmt.Errorf("query failed"))
	fetcher := &mocks.Fetcher{}
	collection := models.Collection{}
	migration, _ := models.NewMigration("/tmp/1.sql", "SELECT 1", models.StatusNotRun)
	err := collection.Add(migration)
	if err != nil {
		test.Fatal(err)
	}
	fetcher.On("GetMigrations", "/tmp/").Return(collection, nil)
	service := NewRunnerService(fetcher, db, "/tmp")

	result, err := service.RunMigrations()

	if err != nil {
		test.Fail()
	}

	if result.GetAll()[0].GetAbsolutePath() != "/tmp/1.sql" || result.GetAll()[0].WasSuccessful() {
		test.Fail()
	}
}

func TestRunningAMigrationSuccessfullyAndThenFailingToRegisterIt(test *testing.T) {
	db := &mocks.DBRepository{}
	db.On("Ping").Return(nil)
	db.On("CreateMigrationsTableIfNeeded").Return(nil)
	db.On("RunMigrationQuery", "SELECT 1").Return(nil)
	db.On("RegisterRunMigration", "1.sql").Return(fmt.Errorf("failed to register run migration"))
	fetcher := &mocks.Fetcher{}
	collection := models.Collection{}
	migration, _ := models.NewMigration("/tmp/1.sql", "SELECT 1", models.StatusNotRun)
	err := collection.Add(migration)
	if err != nil {
		test.Fatal(err)
	}
	fetcher.On("GetMigrations", "/tmp/").Return(collection, nil)
	service := NewRunnerService(fetcher, db, "/tmp")

	result, err := service.RunMigrations()

	if err == nil {
		test.Fail()
	}

	if result.GetAll()[0].GetAbsolutePath() != "/tmp/1.sql" || result.GetAll()[0].HasFailed() {
		test.Fail()
	}
}