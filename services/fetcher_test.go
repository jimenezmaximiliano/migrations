package services_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jimenezmaximiliano/migrations/mocks"
	"github.com/jimenezmaximiliano/migrations/services"
)

const migrationsDir = "/tmp/"

const migrationPath1 = "/tmp/1_a.sql"
const migrationQuery1 = "SELECT 1"
const migrationPath2 = "/tmp/2_b.sql"
const migrationQuery2 = "SELECT 2"

func TestGettingMigrations(test *testing.T) {

	dbRepository := &mocks.DBRepository{}
	dbRepository.On("GetAlreadyRunMigrationFilePaths", migrationsDir).
		Return([]string{migrationPath1}, nil)
	fileRepository := &mocks.FileRepository{}
	fileRepository.On("GetMigrationFilePaths", migrationsDir).
		Return([]string{migrationPath1, migrationPath2}, nil)
	fileRepository.On("GetMigrationQuery", migrationPath1).
		Return(migrationQuery1, nil)
	fileRepository.On("GetMigrationQuery", migrationPath2).
		Return(migrationQuery2, nil)
	service := services.NewFetcherService(dbRepository, fileRepository)

	migrations, err := service.GetMigrations(migrationsDir)

	require.Nil(test, err)

	assert.Len(test, migrations.GetAll(), 2)

	migrationsToRun := migrations.GetMigrationsToRun()
	migrationToRun := migrationsToRun[0]

	assert.Equal(test, migrationPath2, migrationToRun.GetAbsolutePath())
	assert.Equal(test, migrationQuery2, migrationToRun.GetQuery())
}

func TestGettingMigrationsFailsIfItCannotReadFromTheDB(test *testing.T) {
	const migrationsDir = "/tmp/"
	dbRepository := &mocks.DBRepository{}
	dbRepository.On("GetAlreadyRunMigrationFilePaths", migrationsDir).
		Return(nil, fmt.Errorf("db error"))
	fileRepository := &mocks.FileRepository{}
	fileRepository.On("GetMigrationFilePaths", migrationsDir).
		Return(nil, nil)
	service := services.NewFetcherService(dbRepository, fileRepository)

	_, err := service.GetMigrations(migrationsDir)

	assert.NotNil(test, err)
}

func TestGettingMigrationsFailsIfItCannotReadFromTheFileSystem(test *testing.T) {
	const migrationsDir = "/tmp/"
	dbRepository := &mocks.DBRepository{}
	fileRepository := &mocks.FileRepository{}
	fileRepository.On("GetMigrationFilePaths", migrationsDir).
		Return(nil, fmt.Errorf("some fs error"))
	service := services.NewFetcherService(dbRepository, fileRepository)

	_, err := service.GetMigrations(migrationsDir)

	assert.NotNil(test, err)
}

func TestGettingMigrationsFailsIfAMigrationPathCannotBeRead(test *testing.T) {

	dbRepository := &mocks.DBRepository{}
	dbRepository.On("GetAlreadyRunMigrationFilePaths", migrationsDir).
		Return(nil, nil)
	fileRepository := &mocks.FileRepository{}
	fileRepository.On("GetMigrationFilePaths", migrationsDir).
		Return([]string{migrationPath1}, nil)
	fileRepository.On("GetMigrationQuery", migrationPath1).
		Return("", fmt.Errorf("cannot read file"))
	service := services.NewFetcherService(dbRepository, fileRepository)

	_, err := service.GetMigrations(migrationsDir)

	assert.NotNil(test, err)
}

func TestGettingMigrationsFailsIfAMigrationPathAlreadyRunCannotBeRead(test *testing.T) {

	dbRepository := &mocks.DBRepository{}
	dbRepository.On("GetAlreadyRunMigrationFilePaths", migrationsDir).
		Return([]string{migrationPath1}, nil)
	fileRepository := &mocks.FileRepository{}
	fileRepository.On("GetMigrationFilePaths", migrationsDir).
		Return([]string{migrationPath1}, nil)
	fileRepository.On("GetMigrationQuery", migrationPath1).
		Return("", fmt.Errorf("cannot read file"))
	service := services.NewFetcherService(dbRepository, fileRepository)

	_, err := service.GetMigrations(migrationsDir)

	assert.NotNil(test, err)
}
