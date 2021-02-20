package services

import (
	"fmt"

	"github.com/jimenezmaximiliano/migrations/migrations/helpers"
	"github.com/jimenezmaximiliano/migrations/migrations/migration"
	"github.com/jimenezmaximiliano/migrations/migrations/repositories"
)

// Runner handles running migrations.
type Runner interface {
	RunMigrations() (migration.Collection, error)
}

type runnerService struct {
	migrationFetcherService         fetcher
	dbRepository                    repositories.DBRepository
	migrationsDirectoryAbsolutePath string
}

// Ensure runnerService implements Runner
var _ Runner = runnerService{}

// Ensure runnerService implements Runner
var _ Runner = runnerService{}

// NewRunnerService returns an implementation of Runner
func NewRunnerService(
	migrationFetcherService fetcher,
	DBRepository repositories.DBRepository,
	migrationsDirectoryAbsolutePath string) Runner {

	return runnerService{
		migrationFetcherService:         migrationFetcherService,
		dbRepository:                    DBRepository,
		migrationsDirectoryAbsolutePath: helpers.AddTrailingSlashToPathIfNeeded(migrationsDirectoryAbsolutePath),
	}
}

// RunMigrations runs a collection of migrations checking first if they have been run already
func (service runnerService) RunMigrations() (migration.Collection, error) {
	err := service.dbRepository.Ping()
	if err != nil {
		return migration.Collection{}, err
	}

	err = service.dbRepository.CreateMigrationsTableIfNeeded()
	if err != nil {
		return migration.Collection{}, err
	}

	allMigrations, err := service.migrationFetcherService.GetMigrations(service.migrationsDirectoryAbsolutePath)
	if err != nil {
		return migration.Collection{}, err
	}

	migrationsToRun := allMigrations.GetMigrationsToRun()

	if len(migrationsToRun) == 0 {
		return migration.Collection{}, nil
	}

	return service.runMigrations(migrationsToRun)
}

func (service runnerService) runMigrations(migrationsToRun []migration.Migration) (migration.Collection, error) {

	result := migration.Collection{}

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
