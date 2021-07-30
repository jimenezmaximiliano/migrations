package services

import (
	"fmt"
	"testing"

	"github.com/jimenezmaximiliano/migrations/mocks"
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
	service := NewFetcherService(dbRepository, fileRepository)

	migrations, err := service.GetMigrations(migrationsDir)

	if err != nil {
		test.Error(err)
	}

	if len(migrations.GetAll()) != 2 {
		test.Error("Expected only 1 migration in the collection")
	}

	migrationsToRun := migrations.GetMigrationsToRun()
	migrationToRun := migrationsToRun[0]

	if migrationToRun.GetAbsolutePath() != migrationPath2 {
		test.Errorf("Expected migration absolute path %s but got %s", migrationPath2, migrationToRun.GetAbsolutePath())
	}

	if migrationToRun.GetQuery() != migrationQuery2 {
		test.Errorf("Expected migration query %s but got %s", migrationQuery2, migrationToRun.GetQuery())
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

func TestGettingMigrationsFailsIfAMigrationPathCannotBeRead(test *testing.T) {

	dbRepository := &mocks.DBRepository{}
	dbRepository.On("GetAlreadyRunMigrationFilePaths", migrationsDir).
		Return(nil, nil)
	fileRepository := &mocks.FileRepository{}
	fileRepository.On("GetMigrationFilePaths", migrationsDir).
		Return([]string{migrationPath1}, nil)
	fileRepository.On("GetMigrationQuery", migrationPath1).
		Return("", fmt.Errorf("cannot read file"))
	service := NewFetcherService(dbRepository, fileRepository)

	_, err := service.GetMigrations(migrationsDir)

	if err == nil {
		test.Error(err)
	}
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
	service := NewFetcherService(dbRepository, fileRepository)

	_, err := service.GetMigrations(migrationsDir)

	if err == nil {
		test.Error(err)
	}
}
