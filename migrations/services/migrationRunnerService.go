package services

import (
	"fmt"

	"github.com/jimenezmaximiliano/migrations/migrations/migration"
	"github.com/jimenezmaximiliano/migrations/migrations/repositories"
)

// MigrationRunnerService is an interface that handles running migrations
type MigrationRunnerService interface {
	RunMigrations() (migration.MigrationCollection, error)
}

type migrationRunnerService struct {
	migrationFetcherService         MigrationFetcherService
	dbRepository                    repositories.DbRepository
	migrationsDirectoryAbsolutePath string
}

func NewMigrationRunnerService(
	migrationFetcherService MigrationFetcherService,
	dbRepository repositories.DbRepository,
	migrationsDirectoryAbsolutePath string) MigrationRunnerService {

	return migrationRunnerService{
		migrationFetcherService:         migrationFetcherService,
		dbRepository:                    dbRepository,
		migrationsDirectoryAbsolutePath: migrationsDirectoryAbsolutePath,
	}
}

func (service migrationRunnerService) RunMigrations() (migration.MigrationCollection, error) {

	err := service.dbRepository.CreateMigrationsTableIfNeeded()
	if err != nil {
		return migration.MigrationCollection{}, err
	}

	allMigrations, err := service.migrationFetcherService.GetMigrations(service.migrationsDirectoryAbsolutePath)
	migrationsToRun := allMigrations.GetMigrationsToRun()

	if len(migrationsToRun) == 0 {
		return migration.MigrationCollection{}, nil
	}

	return service.runMigrations(migrationsToRun)
}

func (service migrationRunnerService) runMigrations(migrationsToRun []migration.Migration) (migration.MigrationCollection, error) {

	result := migration.MigrationCollection{}

	for _, migration := range migrationsToRun {
		err := service.dbRepository.RunMigrationQuery(migration.GetQuery())

		if err != nil {
			result.Add(migration.NewAsFailed())

			return result, err
		}

		result.Add(migration.NewAsSuccessful())
		err = service.dbRepository.RegisterRunMigration(migration.GetName())

		if err != nil {
			return result, fmt.Errorf("migrations.runnerService (absolutePath: %s) %w", migration.GetAbsolutePath(), err)
		}
	}

	return result, nil
}
