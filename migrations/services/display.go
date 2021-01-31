package services

import (
	"github.com/jimenezmaximiliano/migrations/migrations/migration"
)

type Display interface {
	DisplayRunMigrations(migrations migration.MigrationCollection)
	DisplaySetupError(err error)
}

type Print func(format string, a ...interface{}) (n int, err error)

func NewDisplayService(print Print) Display {
	return displayService{
		print: print,
	}
}

type displayService struct {
	print Print
}

const (
	messageFormat        = "\n[%s] %s"
	informationalMessage = " INFO "
	successfulMigration  = "  OK  "
	failedMigration      = "  KO  "
)

func (service displayService) DisplayRunMigrations(migrations migration.MigrationCollection) {

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

func (service displayService) DisplaySetupError(err error) {
	service.print("\nFailed to setup migrations:\n%v\n\n", err)
}

func (service displayService) info(message string) {
	service.print(messageFormat, informationalMessage, message)
}

func (service displayService) success(message string) {
	service.print(messageFormat, successfulMigration, message)
}

func (service displayService) failure(message string) {
	service.print(messageFormat, failedMigration, message)
}
