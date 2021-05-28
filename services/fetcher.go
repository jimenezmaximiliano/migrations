package services

import (
	"github.com/jimenezmaximiliano/migrations/models"
	"github.com/jimenezmaximiliano/migrations/repositories"
)

// Fetcher fetches Migrations from a given directory.
type Fetcher interface {
	GetMigrations(migrationsDirectoryAbsolutePath string) (models.Collection, error)
}

type fetcherService struct {
	dbRepository   repositories.DBRepository
	fileRepository repositories.FileRepository
}

// Ensure fetcherService implements Fetcher.
var _ Fetcher = fetcherService{}

// NewFetcherService returns an implementation of MigrationFetcherService.
func NewFetcherService(dbRepository repositories.DBRepository, fileRepository repositories.FileRepository) Fetcher {
	return fetcherService{
		dbRepository:   dbRepository,
		fileRepository: fileRepository,
	}
}

// GetMigrations returns a collection of Migrations from a given directory.
func (service fetcherService) GetMigrations(migrationsDirectoryAbsolutePath string) (models.Collection, error) {
	migrationFilePathsFromFiles, runMigrationFilePaths, err := service.
		readMigrationPathsFromTheFileSystemAndTheDB(migrationsDirectoryAbsolutePath)
	if err != nil {
		return models.Collection{}, err
	}

	collection, err := service.parseRunMigrationsFromDB(runMigrationFilePaths)
	if err != nil {
		return models.Collection{}, err
	}

	collection, err = service.parseMigrationsFromFiles(migrationFilePathsFromFiles, collection)
	if err != nil {
		return models.Collection{}, err
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

func (service fetcherService) parseRunMigrationsFromDB(filePaths []string) (models.Collection, error) {
	collection := models.Collection{}
	for _, filePath := range filePaths {
		migrationQuery, err := service.fileRepository.GetMigrationQuery(filePath)
		if err != nil {
			return collection, err
		}

		migration, err := models.NewMigration(filePath, migrationQuery, models.StatusSuccessful)
		if err != nil {
			return collection, err
		}

		collection.Add(migration)
	}

	return collection, nil
}

func (service fetcherService) parseMigrationsFromFiles(filePaths []string, collection models.Collection) (models.Collection, error) {
	for _, migrationFilePath := range filePaths {
		if collection.ContainsMigrationPath(migrationFilePath) {
			continue
		}
		migrationQuery, err := service.fileRepository.GetMigrationQuery(migrationFilePath)
		if err != nil {
			return collection, err
		}

		migration, err := models.NewMigration(migrationFilePath, migrationQuery, models.StatusNotRun)
		if err != nil {
			return collection, err
		}

		collection.Add(migration)
	}

	return collection, nil
}
