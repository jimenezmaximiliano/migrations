package services

import (
	"github.com/jimenezmaximiliano/migrations/migrations/migration"
	"github.com/jimenezmaximiliano/migrations/migrations/repositories"
)

// Fetcher fetches Migrations from a given directory.
type Fetcher interface {
	GetMigrations(migrationsDirectoryAbsolutePath string) (migration.MigrationCollection, error)
}

type fetcherService struct {
	dbRepository   repositories.DbRepository
	fileRepository repositories.FileRepository
}

// NewFetcherService returns an implemention of MigrationFetcherService.
func NewFetcherService(dbRepository repositories.DbRepository, fileRepository repositories.FileRepository) Fetcher {
	return fetcherService{
		dbRepository:   dbRepository,
		fileRepository: fileRepository,
	}
}

// GetMigrations returns a collection of Migrations from a given directory.
func (service fetcherService) GetMigrations(migrationsDirectoryAbsolutePath string) (migration.MigrationCollection, error) {
	migrationFilePathsFromFiles, runMigrationFilePaths, err := service.
		readMigrationPathsFromTheFileSystemAndTheDB(migrationsDirectoryAbsolutePath)
	if err != nil {
		return migration.MigrationCollection{}, err
	}

	collection, err := service.parseRunMigrationsFromDB(runMigrationFilePaths)
	if err != nil {
		return migration.MigrationCollection{}, err
	}

	collection, err = service.parseMigrationsFromFiles(migrationFilePathsFromFiles, collection)
	if err != nil {
		return migration.MigrationCollection{}, err
	}

	return collection, nil
}

func (service fetcherService) readMigrationPathsFromTheFileSystemAndTheDB(
	migrationsDirectoryAbsolutePath string,
) (pathsFromFiles []string, pathsFromDB []string, err error) {
	pathsFromFiles, err = service.fileRepository.GetMigrationFilePaths(migrationsDirectoryAbsolutePath)
	if err != nil {
		return nil, nil, err
	}

	pathsFromDB, err = service.dbRepository.GetAlreadyRunMigrationFilePaths(migrationsDirectoryAbsolutePath)
	if err != nil {
		return pathsFromFiles, nil, err
	}

	return pathsFromFiles, pathsFromDB, nil
}

func (service fetcherService) parseRunMigrationsFromDB(filePaths []string) (migration.MigrationCollection, error) {
	collection := migration.MigrationCollection{}
	for _, filePath := range filePaths {
		migrationQuery, err := service.fileRepository.GetMigrationQuery(filePath)
		if err != nil {
			return collection, err
		}

		migration, err := migration.NewMigration(filePath, migrationQuery, migration.StatusSuccessful)
		if err != nil {
			return collection, err
		}

		collection.Add(migration)
	}

	return collection, nil
}

func (service fetcherService) parseMigrationsFromFiles(filePaths []string, collection migration.MigrationCollection) (migration.MigrationCollection, error) {
	for _, migrationFilePath := range filePaths {
		if collection.ContainsMigrationPath(migrationFilePath) {
			continue
		}
		migrationQuery, err := service.fileRepository.GetMigrationQuery(migrationFilePath)
		if err != nil {
			return collection, err
		}

		migration, err := migration.NewMigration(migrationFilePath, migrationQuery, migration.StatusNotRun)
		if err != nil {
			return collection, err
		}

		collection.Add(migration)
	}

	return collection, nil
}
