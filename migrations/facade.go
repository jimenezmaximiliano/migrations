package migrations

import (
	"database/sql"
	"fmt"

	"github.com/jimenezmaximiliano/migrations/migrations/adapters"
	"github.com/jimenezmaximiliano/migrations/migrations/migration"
	"github.com/jimenezmaximiliano/migrations/migrations/repositories"
	"github.com/jimenezmaximiliano/migrations/migrations/services"
)

func RunMigrations(db *sql.DB, migrationsDirectoryAbsolutePath string) (migration.MigrationCollection, error) {
	fileSystem := adapters.FileSystemAdapter{}
	dbRepository := repositories.NewDbRepository(db)
	fileRepository := repositories.NewFileRepository(fileSystem, fileSystem)
	migrationFetcher := services.NewFetcherService(dbRepository, fileRepository)
	migrationRunner := services.NewRunnerService(migrationFetcher, dbRepository, migrationsDirectoryAbsolutePath)

	return migrationRunner.RunMigrations()
}

func DisplayResults(runMigrations migration.MigrationCollection) {
	displayService := services.NewDisplayService(fmt.Printf)

	displayService.DisplayRunMigrations(runMigrations)
}
