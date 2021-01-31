package services

import (
	"github.com/jimenezmaximiliano/migrations/migrations/migration"
)

type DisplayService interface {
	DisplayRunMigrations(migrations migration.MigrationCollection)
}

type Print func(format string, a ...interface{}) (n int, err error)

func NewDisplayService(print Print) DisplayService {
	return migrationsDisplayService{
		print: print,
	}
}

type migrationsDisplayService struct {
	print Print
}

const (
	messageFormat        = "\n[%s] %s"
	informationalMessage = " INFO "
	successfulMigration  = "  OK  "
	failedMigration      = "  KO  "
)

func (service migrationsDisplayService) DisplayRunMigrations(migrations migration.MigrationCollection) {

	service.info("Run migrations")

	if migrations.IsEmpty() {
		service.info("No migrations to run")
		service.info("Done")
		service.print("\n\n")
		return
	}

	migrationProcessHasFailed := false

	for _, migration := range migrations.GetAll() {
		if migration.WasSuccessful() {
			service.success(migration.GetName())
			continue
		}

		service.failure(migration.GetName())
		migrationProcessHasFailed = true
	}

	if migrationProcessHasFailed {
		service.failure("The migration process has failed")
	}

	service.info("Done")
	service.print("\n\n")
}

func (service migrationsDisplayService) info(message string) {
	service.print(messageFormat, informationalMessage, message)
}

func (service migrationsDisplayService) success(message string) {
	service.print(messageFormat, successfulMigration, message)
}

func (service migrationsDisplayService) failure(message string) {
	service.print(messageFormat, failedMigration, message)
}
