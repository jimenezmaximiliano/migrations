package migrations

import (
	"database/sql"
	"fmt"
	"github.com/jimenezmaximiliano/very-simple-migrations/migrations/migration"
	"github.com/jimenezmaximiliano/very-simple-migrations/migrations/repositories/dbRepository"
	"github.com/jimenezmaximiliano/very-simple-migrations/migrations/repositories/fileRepository"
	"io/ioutil"
	"strings"
)

func Run(db *sql.DB, migrationsAbsolutePath string) ([]migration.Migration, error) {
	dbRepository := DbRepository.New(db)

	err := dbRepository.CreateMigrationsTableIfNeeded()
	if err != nil {
		return []migration.Migration{}, err
	}

	migrationsToRun, err := getMigrations(dbRepository, migrationsAbsolutePath)
	if err != nil {
		return []migration.Migration{}, err
	}

	if len(migrationsToRun) == 0  {
		return []migration.Migration{}, nil
	}

	return runMigrations(dbRepository, migrationsToRun)
}

func getMigrations(dbRepository DbRepository.DbRepository, migrationsAbsolutePath string) ([]migration.Migration, error) {
	migrationsFromFiles, err := FileRepository.GetMigrations(migrationsAbsolutePath)

	if err != nil {
		return []migration.Migration{}, fmt.Errorf("verySimpleMigrations.getMigrations.fromFiles (%s)\n%w", migrationsAbsolutePath, err)
	}

	alreadyRunMigrationFilePaths, err := dbRepository.GetAlreadyRunMigrationFilePaths(migrationsAbsolutePath)

	if err != nil {
		return []migration.Migration{}, fmt.Errorf("verySimpleMigrations.getAlreadyRunMigrations \n%w", err)
	}

	return filterMigrationFilePaths(migrationFilePaths, alreadyRunMigrationFilePaths), nil
}

func filterMigrationFilePaths(all []RunMigration.RunMigration, alreadyRun []string) []RunMigration.RunMigration {
	var migrationsToRun []string

	for _, migration := range all {
		if inSlice(alreadyRun, migration) {
			continue
		}

		migrationsToRun = append(migrationsToRun, migration)
	}

	return migrationsToRun
}

func inSlice(slice []string, element string) bool {
	for _, currentElement := range slice {
		if currentElement == element {
			return true
		}
	}

	return false
}

func runMigrations(dbRepository DbRepository.DbRepository, migrations []RunMigration.RunMigration, result *migrationsResult.MigrationsResult) error {
	for _, migration := range migrations {
		err := runMigration(dbRepository, migration, result)

		if err != nil {
			return err
		}
	}

	return nil
}

func runMigration(dbRepository DbRepository.DbRepository, migrationAbsoluteFilePath string, result *migrationsResult.MigrationsResult) error {
	query, err := ioutil.ReadFile(migrationAbsoluteFilePath)

	if err != nil {
		result.FailedMigrations = append(result.FailedMigrations, migrationAbsoluteFilePath)
		return fmt.Errorf("verySimpleMigrations.readMigrationFile (%s) \n%w", migrationAbsoluteFilePath, err)
	}

	err = dbRepository.RunMigrationQuery(string(query))

	if err != nil {
		result.FailedMigrations = append(result.FailedMigrations, migrationAbsoluteFilePath)
		return fmt.Errorf("verySimpleMigrations.runMigration (%s) \n%w", migrationAbsoluteFilePath, err)
	}

	migrationFileNameParts := strings.Split(migrationAbsoluteFilePath, "/")
	migrationFileName := migrationFileNameParts[len(migrationFileNameParts) - 1]

	err = dbRepository.RegisterRunMigration(migrationFileName)

	if err != nil {
		return fmt.Errorf("verySimpleMigrations.insertMigrationIntoMigrationsTableAfterRunningIt (%s) \n%w", migrationAbsoluteFilePath, err)
	}

	result.SuccessfulMigrations = append(result.SuccessfulMigrations, migrationAbsoluteFilePath)

	return nil
}