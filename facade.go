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
	arguments := services.Arguments{
		MigrationsPath: migrationsDirectoryAbsolutePath,
	}
	fileRepository := repositories.NewFileRepository(adapters.IOUtilAdapter{})
	migrationRunner := getMigrationRunner(DB, fileRepository, arguments)

	return migrationRunner.RunMigrations()
}

// SetupDB is a function that handles the configuration for the DB connection.
type SetupDB func() (*sql.DB, error)

// RunMigrationsCommand runs migrations as a command (it will output the results to stdout).
func RunMigrationsCommand(setupDB SetupDB) {
	displayService := getDisplayService()
	arguments, argumentsAreValid := getArgumentService(displayService).ParseAndValidate()
	if !argumentsAreValid {
		os.Exit(1)
	}

	DB, err := setupDB()
	if err != nil {
		displayService.DisplayErrorWithMessage(err, "failed to setup the DB")
		os.Exit(1)
	}

	fileRepository := repositories.NewFileRepository(adapters.IOUtilAdapter{})
	migrationRunner := getMigrationRunner(DB, fileRepository, arguments)

	switch arguments.Command {
	case "migrate":
		result, err := migrationRunner.RunMigrations()
		if err != nil {
			displayService.DisplayErrorWithMessage(err, "something went wrong while running migrations")
			os.Exit(1)
		}
		displayService.DisplayRunMigrations(result)
	case "create":
		commands.NewCreateMigrationCommand(fileRepository, displayService, arguments).Run()
	}

	os.Exit(0)
}

func getMigrationRunner(DB *sql.DB, fileRepository repositories.FileRepository, arguments services.Arguments) services.Runner {
	dbAdapter := adapters.NewDBAdapter(DB)
	dbRepository := repositories.NewDBRepository(dbAdapter)
	migrationFetcher := services.NewFetcherService(dbRepository, fileRepository)

	return services.NewRunnerService(migrationFetcher, dbRepository, arguments.MigrationsPath)
}

func getDisplayService() services.DisplayService {
	printerAdapter := adapters.PrinterAdapter{}

	return services.NewDisplayService(printerAdapter)
}

func getArgumentService(displayService services.Display) services.CommandArgumentService {
	return services.NewCommandArgumentService(displayService, adapters.NewArgumentParser())
}
