package services

import (
	"fmt"
	"testing"

	"github.com/jimenezmaximiliano/migrations/mocks"
)

func TestGettingMigrations(test *testing.T) {
	const migrationsDir = "/tmp/"
	const migrationPath = "/tmp/1.sql"
	const migrationQuery = "SELECT 1"
	dbRepository := &mocks.DBRepository{}
	dbRepository.On("GetAlreadyRunMigrationFilePaths", migrationsDir).
		Return(nil, nil)
	fileRepository := &mocks.FileRepository{}
	fileRepository.On("GetMigrationFilePaths", migrationsDir).
		Return([]string{migrationPath}, nil)
	fileRepository.On("GetMigrationQuery", migrationPath).
		Return(migrationQuery, nil)
	service := NewFetcherService(dbRepository, fileRepository)

	migrations, err := service.GetMigrations(migrationsDir)

	if err != nil {
		test.Error(err)
	}

	if len(migrations.GetAll()) != 1 {
		test.Error("Expected only 1 migration in the collection")
	}

	migrationsToRun := migrations.GetMigrationsToRun()
	resultMigration := migrationsToRun[0]

	if resultMigration.GetAbsolutePath() != migrationPath {
		test.Errorf("Expected migration absolute path %s but got %s", migrationPath, resultMigration.GetAbsolutePath())
	}

	if resultMigration.GetQuery() != migrationQuery {
		test.Errorf("Expected migration query %s but got %s", migrationQuery, resultMigration.GetQuery())
	}
}

func TestGettingMigrationsFailsIfItCannotReadFromTheDB(test *testing.T) {
	const migrationsDir = "/tmp/"
	dbRepository := &mocks.DBRepository{}
	dbRepository.On("GetAlreadyRunMigrationFilePaths", migrationsDir).
		Return(nil, fmt.Errorf("db error"))
	fileRepository := &mocks.FileRepository{}
	fileRepository.On("GetMigrationFilePaths", migrationsDir).
		Return(nil, nil)
	service := NewFetcherService(dbRepository, fileRepository)

	_, err := service.GetMigrations(migrationsDir)

	if err == nil {
		test.Error(err)
	}
}

func TestGettingMigrationsFailsIfItCannotReadFromTheFileSystem(test *testing.T) {
	const migrationsDir = "/tmp/"
	dbRepository := &mocks.DBRepository{}
	fileRepository := &mocks.FileRepository{}
	fileRepository.On("GetMigrationFilePaths", migrationsDir).
		Return(nil, fmt.Errorf("some fs error"))
	service := NewFetcherService(dbRepository, fileRepository)

	_, err := service.GetMigrations(migrationsDir)

	if err == nil {
		test.Error(err)
	}
}
