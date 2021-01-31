package services

import (
	"fmt"

	"github.com/jimenezmaximiliano/migrations/migrations/migration"
	"github.com/jimenezmaximiliano/migrations/migrations/repositories"
)

// FetcherService returns Migrations from a given path
type FetcherService interface {
	GetMigrations(migrationsDirectoryAbsolutePath string) (migration.MigrationCollection, error)
}

type migrationFetcherService struct {
	dbRepository   repositories.DbRepository
	fileRepository repositories.FileRepository
}

// NewFetcherService returns an implemention of MigrationFetcherService
func NewFetcherService(dbRepository repositories.DbRepository, fileRepository repositories.FileRepository) FetcherService {
	return migrationFetcherService{
		dbRepository:   dbRepository,
		fileRepository: fileRepository,
	}
}

func (service migrationFetcherService) GetMigrations(migrationsDirectoryAbsolutePath string) (migration.MigrationCollection, error) {
	migrationFilePathsFromFiles, err := service.fileRepository.GetMigrationFilePaths(migrationsDirectoryAbsolutePath)

	if err != nil {
		return migration.MigrationCollection{}, fmt.Errorf("migrations.fetcherService (absolutePath: %s) %w", migrationsDirectoryAbsolutePath, err)
	}

	runMigrationFilePaths, err := service.dbRepository.GetAlreadyRunMigrationFilePaths(migrationsDirectoryAbsolutePath)

	if err != nil {
		return migration.MigrationCollection{}, fmt.Errorf("migrations.fetcherService (absolutePath: %s) %w", migrationsDirectoryAbsolutePath, err)
	}

	migrationCollection := migration.MigrationCollection{}

	for _, runMigrationFilePath := range runMigrationFilePaths {

		migrationQuery, err := service.fileRepository.GetMigrationQuery(runMigrationFilePath)

		if err != nil {
			return migrationCollection, fmt.Errorf("migrations.fetcherService (absolutePath: %s) %w", runMigrationFilePath, err)
		}

		migration, err := migration.NewMigration(runMigrationFilePath, migrationQuery, migration.StatusSuccessful)

		if err != nil {
			return migrationCollection, fmt.Errorf("migrations.fetcherService (absolutePath: %s) %w", runMigrationFilePath, err)
		}

		migrationCollection.Add(migration)
	}

	for _, migrationFilePath := range migrationFilePathsFromFiles {

		if migrationCollection.ContainsMigrationPath(migrationFilePath) {
			continue
		}

		migrationQuery, err := service.fileRepository.GetMigrationQuery(migrationFilePath)

		if err != nil {
			return migrationCollection, fmt.Errorf("migrations.fetcherService (absolutePath: %s) %w", migrationFilePath, err)
		}

		migration, err := migration.NewMigration(migrationFilePath, migrationQuery, migration.StatusNotRun)

		if err != nil {
			return migrationCollection, fmt.Errorf("migrations.fetcherService (absolutePath: %s) %w", migrationFilePath, err)
		}

		migrationCollection.Add(migration)
	}

	return migrationCollection, nil
}
