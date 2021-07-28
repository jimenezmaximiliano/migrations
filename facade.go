package migrations

import (
	"database/sql"
	"os"

	"github.com/jimenezmaximiliano/migrations/adapters"
	"github.com/jimenezmaximiliano/migrations/models"
	"github.com/jimenezmaximiliano/migrations/repositories"
	"github.com/jimenezmaximiliano/migrations/services"
)

// RunMigrations runs the migrations using the given DB connection and migrations directory path.
// Returns a MigrationCollection, to be used programmatically.
func RunMigrations(DB *sql.DB, migrationsDirectoryAbsolutePath string) (models.Collection, error) {
	fileSystem := adapters.IOUtilAdapter{}
	dbAdapter := adapters.NewDBAdapter(DB)
	dbRepository := repositories.NewDBRepository(dbAdapter)
	fileRepository := repositories.NewFileRepository(fileSystem)
	migrationFetcher := services.NewFetcherService(dbRepository, fileRepository)
	migrationRunner := services.NewRunnerService(migrationFetcher, dbRepository, migrationsDirectoryAbsolutePath)

	return migrationRunner.RunMigrations()
}

// SetupDB is a function that handles the configuration for the DB connection.
type SetupDB func() (*sql.DB, error)

// RunMigrationsCommand runs migrations as a command (it will output the results to stdout).
func RunMigrationsCommand(setupDB SetupDB) {
	printerAdapter := adapters.PrinterAdapter{}
	displayService := services.NewDisplayService(printerAdapter)
	commandService := services.NewCommandService(adapters.NewArgumentParser())
	arguments := commandService.ParseArguments()

	if arguments.MigrationsPath == "" {
		displayService.DisplayHelp()
		os.Exit(0)
	}

	DB, err := setupDB()
	if err != nil {
		displayService.DisplaySetupError(err)
		os.Exit(1)
		return
	}

	fileSystem := adapters.IOUtilAdapter{}
	dbAdapter := adapters.NewDBAdapter(DB)
	dbRepository := repositories.NewDBRepository(dbAdapter)
	fileRepository := repositories.NewFileRepository(fileSystem)
	migrationFetcher := services.NewFetcherService(dbRepository, fileRepository)
	migrationRunner := services.NewRunnerService(migrationFetcher, dbRepository, arguments.MigrationsPath)
	result, err := migrationRunner.RunMigrations()
	if err != nil {
		displayService.DisplayGeneralError(err)
		os.Exit(1)
		return
	}

	displayService.DisplayRunMigrations(result)
	os.Exit(0)
}
