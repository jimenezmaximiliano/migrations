package services

import (
	"github.com/pkg/errors"

	"github.com/jimenezmaximiliano/migrations/helpers"
	"github.com/jimenezmaximiliano/migrations/models"
	"github.com/jimenezmaximiliano/migrations/repositories"
)

// Runner handles running migrations.
type Runner interface {
	RunMigrations() (models.Collection, error)
}

type runnerService struct {
	migrationFetcherService         Fetcher
	dbRepository                    repositories.DBRepository
	migrationsDirectoryAbsolutePath string
}

// Ensure runnerService implements Runner.
var _ Runner = runnerService{}

// Ensure runnerService implements Runner.
var _ Runner = runnerService{}

// NewRunnerService returns an implementation of Runner.
func NewRunnerService(
	migrationFetcherService Fetcher,
	DBRepository repositories.DBRepository,
	migrationsDirectoryAbsolutePath string) Runner {

	return runnerService{
		migrationFetcherService:         migrationFetcherService,
		dbRepository:                    DBRepository,
		migrationsDirectoryAbsolutePath: helpers.AddTrailingSlashToPathIfNeeded(migrationsDirectoryAbsolutePath),
	}
}

// RunMigrations runs a collection of migrations checking first if they have been run already.
func (service runnerService) RunMigrations() (models.Collection, error) {
	err := service.dbRepository.Ping()
	if err != nil {
		return models.Collection{}, errors.Wrap(err, "failed to connect to the DB")
	}

	err = service.dbRepository.CreateMigrationsTableIfNeeded()
	if err != nil {
		return models.Collection{}, err
	}

	allMigrations, err := service.migrationFetcherService.GetMigrations(service.migrationsDirectoryAbsolutePath)
	if err != nil {
		return models.Collection{}, err
	}

	migrationsToRun := allMigrations.GetMigrationsToRun()

	if len(migrationsToRun) == 0 {
		return models.Collection{}, nil
	}

	return service.runMigrations(migrationsToRun)
}

func (service runnerService) runMigrations(migrationsToRun []models.Migration) (models.Collection, error) {

	result := models.Collection{}
	failed := false

	for _, migration := range migrationsToRun {
		if failed {
			err := result.Add(migration.NewAsNotRun())
			if err != nil {
				return result, err
			}
			continue
		}

		err := service.dbRepository.RunMigrationQuery(migration.GetQuery())
		if err != nil {
			err = result.Add(migration.NewAsFailed(errors.WithStack(err)))
			if err != nil {
				return result, err
			}

			failed = true

			continue
		}

		err = result.Add(migration.NewAsSuccessful())
		if err != nil {
			return result, err
		}

		err = service.dbRepository.RegisterRunMigration(migration.GetName())
		if err != nil {
			return result, err
		}
	}

	return result, nil
}
