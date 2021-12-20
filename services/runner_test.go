package services_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jimenezmaximiliano/migrations/mocks"
	"github.com/jimenezmaximiliano/migrations/models"
	"github.com/jimenezmaximiliano/migrations/services"
)

func TestRunningMigrationsFailsIfTheDBConnectionDoesNotWork(test *testing.T) {
	test.Parallel()

	fetcher := &mocks.Fetcher{}
	db := &mocks.DBRepository{}
	db.On("Ping").Return(fmt.Errorf("db connection error"))
	service := services.NewRunnerService(fetcher, db, "/tmp")

	_, err := service.RunMigrations()

	assert.NotNil(test, err)
}

func TestRunningMigrationsFailsIfAMigrationsTableCannotBeCreated(test *testing.T) {
	test.Parallel()

	fetcher := &mocks.Fetcher{}
	db := &mocks.DBRepository{}
	db.On("Ping").Return(nil)
	db.On("CreateMigrationsTableIfNeeded").Return(fmt.Errorf("cannot create table"))
	service := services.NewRunnerService(fetcher, db, "/tmp")

	_, err := service.RunMigrations()

	assert.NotNil(test, err)
}

func TestRunningMigrationsFailsIfItCannotFetchMigrationsFromFilesOrTheDB(test *testing.T) {
	test.Parallel()

	db := &mocks.DBRepository{}
	db.On("Ping").Return(nil)
	db.On("CreateMigrationsTableIfNeeded").Return(nil)
	fetcher := &mocks.Fetcher{}
	fetcher.On("GetMigrations", "/tmp/").Return(models.Collection{}, fmt.Errorf("cannot fetch migrations"))
	service := services.NewRunnerService(fetcher, db, "/tmp")

	_, err := service.RunMigrations()

	assert.NotNil(test, err)
}

func TestRunningMigrationsDoesNotFailIfThereAreNoMigratiosToRun(test *testing.T) {
	test.Parallel()

	db := &mocks.DBRepository{}
	db.On("Ping").Return(nil)
	db.On("CreateMigrationsTableIfNeeded").Return(nil)
	fetcher := &mocks.Fetcher{}
	fetcher.On("GetMigrations", "/tmp/").Return(models.Collection{}, nil)
	service := services.NewRunnerService(fetcher, db, "/tmp")

	_, err := service.RunMigrations()

	assert.Nil(test, err)
}

func TestRunningAMigrationSuccessfully(test *testing.T) {
	test.Parallel()

	db := &mocks.DBRepository{}
	db.On("Ping").Return(nil)
	db.On("CreateMigrationsTableIfNeeded").Return(nil)
	db.On("RunMigrationQuery", "SELECT 1").Return(nil)
	db.On("RegisterRunMigration", "1_a.sql").Return(nil)
	fetcher := &mocks.Fetcher{}
	collection := models.Collection{}
	migration, _ := models.NewMigration("/tmp/1_a.sql", "SELECT 1", models.StatusNotRun)
	err := collection.Add(migration)
	require.Nil(test, err)

	fetcher.On("GetMigrations", "/tmp/").Return(collection, nil)
	service := services.NewRunnerService(fetcher, db, "/tmp")

	result, err := service.RunMigrations()

	assert.Nil(test, err)
	assert.Equal(test, "/tmp/1_a.sql", result.GetAll()[0].GetAbsolutePath())
	assert.True(test, result.GetAll()[0].WasSuccessful())
}

func TestRunningAMigrationThatFails(test *testing.T) {
	test.Parallel()

	db := &mocks.DBRepository{}
	db.On("Ping").Return(nil)
	db.On("CreateMigrationsTableIfNeeded").Return(nil)
	db.On("RunMigrationQuery", "SELECT 1").Return(fmt.Errorf("query failed"))
	fetcher := &mocks.Fetcher{}
	collection := models.Collection{}
	migration, _ := models.NewMigration("/tmp/1_a.sql", "SELECT 1", models.StatusNotRun)
	err := collection.Add(migration)
	require.Nil(test, err)

	fetcher.On("GetMigrations", "/tmp/").Return(collection, nil)
	service := services.NewRunnerService(fetcher, db, "/tmp")

	result, err := service.RunMigrations()

	assert.Nil(test, err)
	assert.Equal(test, "/tmp/1_a.sql", result.GetAll()[0].GetAbsolutePath())
	assert.True(test, result.GetAll()[0].HasFailed())
}

func TestRunningAMigrationSuccessfullyAndThenFailingToRegisterIt(test *testing.T) {
	test.Parallel()

	db := &mocks.DBRepository{}
	db.On("Ping").Return(nil)
	db.On("CreateMigrationsTableIfNeeded").Return(nil)
	db.On("RunMigrationQuery", "SELECT 1").Return(nil)
	db.On("RegisterRunMigration", "1_a.sql").Return(fmt.Errorf("failed to register run migration"))
	fetcher := &mocks.Fetcher{}
	collection := models.Collection{}
	migration, _ := models.NewMigration("/tmp/1_a.sql", "SELECT 1", models.StatusNotRun)
	err := collection.Add(migration)
	require.Nil(test, err)

	fetcher.On("GetMigrations", "/tmp/").Return(collection, nil)
	service := services.NewRunnerService(fetcher, db, "/tmp")

	result, err := service.RunMigrations()

	assert.NotNil(test, err)
	assert.Equal(test, "/tmp/1_a.sql", result.GetAll()[0].GetAbsolutePath())
	assert.True(test, result.GetAll()[0].WasSuccessful())
}
