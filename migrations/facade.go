package migrations

import (
	"database/sql"
	"fmt"

	"github.com/jimenezmaximiliano/migrations/migrations/adapters"
	"github.com/jimenezmaximiliano/migrations/migrations/migration"
	"github.com/jimenezmaximiliano/migrations/migrations/repositories"
	"github.com/jimenezmaximiliano/migrations/migrations/services"
)

func RunMigrations(DB *sql.DB, migrationsDirectoryAbsolutePath string) (migration.MigrationCollection, error) {
	fileSystem := adapters.IOUtilAdapter{}
	dbRepository := repositories.NewDbRepository(DB)
	fileRepository := repositories.NewFileRepository(fileSystem, fileSystem)
	migrationFetcher := services.NewFetcherService(dbRepository, fileRepository)
	migrationRunner := services.NewRunnerService(migrationFetcher, dbRepository, migrationsDirectoryAbsolutePath)

	return migrationRunner.RunMigrations()
}

type SetupDB func() (*sql.DB, error)

func RunMigrationsCommand(setupDB SetupDB) {

	displayService := services.NewDisplayService(fmt.Printf)
	commandService := services.NewCommandService(adapters.FlagOptionParser{})
	arguments := commandService.ParseArguments()

	DB, err := setupDB()
	if err != nil {
		displayService.DisplaySetupError(err)
		return
	}

	fileSystem := adapters.IOUtilAdapter{}
	dbRepository := repositories.NewDbRepository(DB)
	fileRepository := repositories.NewFileRepository(fileSystem, fileSystem)
	migrationFetcher := services.NewFetcherService(dbRepository, fileRepository)
	migrationRunner := services.NewRunnerService(migrationFetcher, dbRepository, arguments.MigrationsPath)

	result, err := migrationRunner.RunMigrations()

	if err != nil {
		displayService.DisplayGeneralError(err)
	}

	displayService.DisplayRunMigrations(result)
}