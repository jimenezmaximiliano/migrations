package services

import (
	"fmt"

	"github.com/jimenezmaximiliano/migrations/migrations"
	"github.com/jimenezmaximiliano/migrations/migrations/repositories"
)

// MigrationRunnerService is an interface that handles running migrations
type MigrationRunnerService interface {
	RunMigrations() ([]migrations.Migration, error)
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

func (service migrationRunnerService) RunMigrations() ([]migrations.Migration, error) {

	err := service.dbRepository.CreateMigrationsTableIfNeeded()
	if err != nil {
		return []migrations.Migration{}, err
	}

	migrations, err := service.migrationFetcherService.GetMigrations(service.migrationsDirectoryAbsolutePath)
	migrationsToRun := migrations.GetMigrationsToRun()

	if len(migrationsToRun) == 0 {
		return migrationsToRun, nil
	}

	return service.runMigrations(migrationsToRun)
}

func (service migrationRunnerService) runMigrations(migrationsToRun []migrations.Migration) ([]migrations.Migration, error) {

	result := migrations.MigrationCollection{}

	for _, migration := range migrationsToRun {
		err := service.dbRepository.RunMigrationQuery(migration.GetQuery())

		if err != nil {
			result.Add(migration.NewAsFailed())

			return result.GetAll(), err
		}

		result.Add(migration.NewAsSuccessful())
		err = service.dbRepository.RegisterRunMigration(migration.GetName())

		if err != nil {
			return result.GetAll(), fmt.Errorf("migrations.runnerService (absolutePath: %s) %w", migration.GetAbsolutePath(), err)
		}
	}

	return result.GetAll(), nil
}
