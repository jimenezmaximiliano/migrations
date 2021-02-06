package services

import (
	"fmt"

	"github.com/jimenezmaximiliano/migrations/migrations/migration"
	"github.com/jimenezmaximiliano/migrations/migrations/repositories"
)

// Runner handles running migrations.
type Runner interface {
	RunMigrations() (migration.MigrationCollection, error)
}

type runnerService struct {
	migrationFetcherService         Fetcher
	dbRepository                    repositories.DbRepository
	migrationsDirectoryAbsolutePath string
}

// NewRunnerService returns an implementation of Runner
func NewRunnerService(
	migrationFetcherService Fetcher,
	dbRepository repositories.DbRepository,
	migrationsDirectoryAbsolutePath string) Runner {

	return runnerService{
		migrationFetcherService:         migrationFetcherService,
		dbRepository:                    dbRepository,
		migrationsDirectoryAbsolutePath: migrationsDirectoryAbsolutePath,
	}
}

// RunMigrations runs a collection of migrations checking first if they have been run already
func (service runnerService) RunMigrations() (migration.MigrationCollection, error) {
	err := service.dbRepository.CreateMigrationsTableIfNeeded()
	if err != nil {
		return migration.MigrationCollection{}, err
	}

	allMigrations, err := service.migrationFetcherService.GetMigrations(service.migrationsDirectoryAbsolutePath)
	if err != nil {
		return migration.MigrationCollection{}, err
	}

	migrationsToRun := allMigrations.GetMigrationsToRun()

	if len(migrationsToRun) == 0 {
		return migration.MigrationCollection{}, nil
	}

	return service.runMigrations(migrationsToRun)
}

func (service runnerService) runMigrations(migrationsToRun []migration.Migration) (migration.MigrationCollection, error) {

	result := migration.MigrationCollection{}

	for _, migration := range migrationsToRun {
		err := service.dbRepository.RunMigrationQuery(migration.GetQuery())
		if err != nil {
			result.Add(migration.NewAsFailed())

			return result, nil
		}
		
		result.Add(migration.NewAsSuccessful())
		err = service.dbRepository.RegisterRunMigration(migration.GetName())
		if err != nil {
			return result, fmt.Errorf("could not register the migration as run (absolutePath: %s) %w", migration.GetAbsolutePath(), err)
		}
	}

	return result, nil
}
