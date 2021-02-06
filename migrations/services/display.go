package services

import (
	"github.com/jimenezmaximiliano/migrations/migrations/migration"
)

// Display handles the output of the migrations command.
type Display interface {
	DisplayRunMigrations(migrations migration.MigrationCollection)
	DisplaySetupError(err error)
	DisplayGeneralError(err error)
}

// Print outputs a string given a format.
type Print func(format string, a ...interface{}) (n int, err error)

// NewDisplayService returns an implementation of Display.
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

// DisplayRunMigrations outputs the results of run migrations.
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

// DisplaySetupError outputs an error that occur during the setup process (before running migrations).
func (service displayService) DisplaySetupError(err error) {
	service.print("\nFailed to setup migrations:\n%v\n\n", err)
}

// DisplayGeneralError outputs an error that occur while running a migration.
func (service displayService) DisplayGeneralError(err error) {
	service.print("\nAn error occur while running migrations:\n%v\n\n", err)
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
