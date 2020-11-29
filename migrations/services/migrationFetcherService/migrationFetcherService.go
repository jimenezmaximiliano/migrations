package migrationFetcherService

import (
	"fmt"
	"github.com/jimenezmaximiliano/very-simple-migrations/migrations/migration"
	"github.com/jimenezmaximiliano/very-simple-migrations/migrations/repositories/dbRepository"
	"github.com/jimenezmaximiliano/very-simple-migrations/migrations/repositories/fileRepository"
)

type MigrationFetcherService interface {
	getMigrations(migrationsDirectoryAbsolutePath string) ([]migration.Migration, error)
}

type migrationFetcherService struct {
	dbRepository dbRepository.DbRepository
	fileRepository fileRepository.FileRepository
}

func New(dbRepository dbRepository.DbRepository, fileRepository fileRepository.FileRepository) MigrationFetcherService {
	return migrationFetcherService{
		dbRepository: dbRepository,
		fileRepository: fileRepository,
	}
}

func (service migrationFetcherService) getMigrations(migrationsDirectoryAbsolutePath string) ([]migration.Migration, error) {
	migrationFilePathsFromFiles, err := service.fileRepository.GetMigrationFilePaths(migrationsDirectoryAbsolutePath)

	if err != nil {
		return []migration.Migration{}, fmt.Errorf("verySimpleMigrations.fetcherService.getMigrationsFromFiles (absolutePath: %s) %w", migrationsDirectoryAbsolutePath, err)
	}

	runMigrationFilePaths, err := service.dbRepository.GetAlreadyRunMigrationFilePaths(migrationsDirectoryAbsolutePath)

	if err != nil {
		return []migration.Migration{}, fmt.Errorf("verySimpleMigrations.fetcherService.getAlreadyRunMigrationFilePaths (absolutePath: %s) %w", migrationsDirectoryAbsolutePath, err)
	}
}