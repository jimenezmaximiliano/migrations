package migrations

import (
	"database/sql"
	"os"

	"github.com/jimenezmaximiliano/migrations/adapters"
	"github.com/jimenezmaximiliano/migrations/commands"
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
	argumentService := services.NewCommandArgumentService(displayService, adapters.NewArgumentParser())

	arguments, argumentsAreValid := argumentService.ParseAndValidate()
	if !argumentsAreValid {
		os.Exit(1)
	}

	DB, err := setupDB()
	if err != nil {
		displayService.DisplayErrorWithMessage(err, "failed to setup the DB")
		os.Exit(1)
	}

	fileSystem := adapters.IOUtilAdapter{}
	dbAdapter := adapters.NewDBAdapter(DB)
	dbRepository := repositories.NewDBRepository(dbAdapter)
	fileRepository := repositories.NewFileRepository(fileSystem)
	migrationFetcher := services.NewFetcherService(dbRepository, fileRepository)
	migrationRunner := services.NewRunnerService(migrationFetcher, dbRepository, arguments.MigrationsPath)

	switch arguments.Command {
	case "migrate":
		result, err := migrationRunner.RunMigrations()
		if err != nil {
			displayService.DisplayErrorWithMessage(err, "something went wrong while running migrations")
			os.Exit(1)
		}
		displayService.DisplayRunMigrations(result)
		break
	case "create":
		commands.NewCreateMigrationCommand(fileRepository, displayService, arguments).Run()
	}

	os.Exit(0)
}
