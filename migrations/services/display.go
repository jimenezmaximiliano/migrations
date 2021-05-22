package services

import (
	"fmt"
	"os"

	"github.com/jimenezmaximiliano/migrations/migrations/adapters"
	"github.com/jimenezmaximiliano/migrations/migrations/migration"
)

// Display handles the output of the migrations command.
type Display interface {
	DisplayRunMigrations(migrations migration.Collection)
	DisplaySetupError(err error)
	DisplayGeneralError(err error)
	DisplayHelp()
}

type displayService struct {
	printer adapters.Printer
}

// Ensure displayService implements Display.
var _ Display = displayService{}

// NewDisplayService returns an implementation of Display.
func NewDisplayService(printer adapters.Printer) Display {
	return displayService{
		printer: printer,
	}
}

const (
	messageFormat        = "\n[%s] %s"
	informationalMessage = " INFO "
	successfulMigration  = "  OK  "
	failedMigration      = " FAIL "
)

// DisplayRunMigrations outputs the results of run migrations.
func (service displayService) DisplayRunMigrations(migrations migration.Collection) {
	service.info("Run migrations")
	if migrations.IsEmpty() {
		service.info("No migrations to run")
		service.info("Done")
		service.printer.Print(os.Stdout, "\n\n")
		return
	}
	migrationProcessHasFailed := false
	for _, migration := range migrations.GetAll() {
		if migration.WasSuccessful() {
			service.success(migration.GetName())
			continue
		}

		if migration.HasFailed() {
			service.failure(fmt.Sprintf("Migration %s failed with error [%s]", migration.GetAbsolutePath(), migration.GetError()))
			migrationProcessHasFailed = true
			continue
		}

		service.info(fmt.Sprintf("Not run: %s", migration.GetName()))
	}

	if migrationProcessHasFailed {
		service.failure("The migration process has failed")
	}

	service.info("Done")
	service.printer.Print(os.Stdout, "\n\n")
}

// DisplaySetupError outputs an error that occur during the setup process (before running migrations).
func (service displayService) DisplaySetupError(err error) {
	service.printer.Print(os.Stderr, "\nFailed to setup migrations:\n%v\n\n", err)
}

// DisplayGeneralError outputs an error that occur while running a migration.
func (service displayService) DisplayGeneralError(err error) {
	service.printer.Print(os.Stderr, "\nAn error occur while running migrations:\n%v\n\n", err)
}

func (service displayService) info(message string) {
	service.printer.Print(os.Stdout, messageFormat, informationalMessage, message)
}

func (service displayService) success(message string) {
	service.printer.Print(os.Stdout, messageFormat, successfulMigration, message)
}

func (service displayService) failure(message string) {
	service.printer.Print(os.Stderr, messageFormat, failedMigration, message)
}

func (service displayService) DisplayHelp() {
	service.printer.Print(os.Stdout, "Documentation: https://github.com/jimenezmaximiliano/migrations\n\n")
	service.printer.Print(os.Stdout, "Usage:\n")
	service.printer.Print(os.Stdout, "\tmigrate -path=/path/to/migrations/directory/\n\n")
}
