package migrations

import (
	"database/sql"

	"github.com/jimenezmaximiliano/migrations/migrations/adapters"
	"github.com/jimenezmaximiliano/migrations/migrations/migration"
	"github.com/jimenezmaximiliano/migrations/migrations/repositories"
	"github.com/jimenezmaximiliano/migrations/migrations/services"
)

func RunMigrations(db *sql.DB, migrationsDirectoryAbsolutePath string) (migration.MigrationCollection, error) {
	fileSystem := adapters.FileSystemAdapter{}
	dbRepository := repositories.NewDbRepository(db)
	fileRepository := repositories.NewFileRepository(fileSystem, fileSystem)
	migrationFetcher := services.NewMigrationFetcherService(dbRepository, fileRepository)
	migrationRunner := services.NewMigrationRunnerService(migrationFetcher, dbRepository, migrationsDirectoryAbsolutePath)

	return migrationRunner.RunMigrations()
}
