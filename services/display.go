package services

import (
	"fmt"
	"os"

	"github.com/jimenezmaximiliano/migrations/adapters"
	"github.com/jimenezmaximiliano/migrations/models"
)

// Display handles the output of the migrations command.
type Display interface {
	DisplayRunMigrations(migrations models.Collection)
	DisplayErrorWithMessage(err error, message string)
	DisplayError(err error)
	// Deprecated: use DisplayError instead
	DisplaySetupError(err error)
	// Deprecated: use DisplayError instead
	DisplayGeneralError(err error)
	DisplayHelp()
}

type DisplayService struct {
	printer adapters.Printer
}

// Ensure DisplayService implements Display.
var _ Display = DisplayService{}

// NewDisplayService returns an implementation of Display.
func NewDisplayService(printer adapters.Printer) DisplayService {
	return DisplayService{
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
func (service DisplayService) DisplayRunMigrations(migrations models.Collection) {
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
	_ = service.printer.Print(os.Stdout, "\n\n")
}

func (service DisplayService) DisplayError(err error) {
	_ = service.printer.Print(os.Stderr, "\n[ERROR] %s\n", err)
}

func (service DisplayService) DisplayErrorWithMessage(err error, message string) {
	_ = service.printer.Print(os.Stderr, "\n[ERROR] %s: %s\n", message, err)
}

func (service DisplayService) info(message string) {
	_ = service.printer.Print(os.Stdout, messageFormat, informationalMessage, message)
}

func (service DisplayService) success(message string) {
	_ = service.printer.Print(os.Stdout, messageFormat, successfulMigration, message)
}

func (service DisplayService) failure(message string) {
	_ = service.printer.Print(os.Stderr, messageFormat, failedMigration, message)
}

func (service DisplayService) DisplayHelp() {
	_ = service.printer.Print(os.Stdout, "\nUsage:\n")
	_ = service.printer.Print(os.Stdout, "\t[executable] [command] [-options]\n\n")
	_ = service.printer.Print(os.Stdout, "\tExamples:\n\n")
	_ = service.printer.Print(os.Stdout, "\tgo run main.go migrate -path=/path/to/migrations/directory/\n")
	_ = service.printer.Print(os.Stdout, "\t./myMigrationBinary migrate -path=/path/to/migrations/directory/\n\n")
	_ = service.printer.Print(os.Stdout, "\nDocumentation: https://github.com/jimenezmaximiliano/migrations\n\n")
}

// Deprecated: use DisplayError instead
// DisplaySetupError outputs an error that occur during the setup process (before running migrations).
func (service DisplayService) DisplaySetupError(err error) {
	service.DisplayErrorWithMessage(err, "failed to setup migrations")
}

// Deprecated: use DisplayError instead
// DisplayGeneralError outputs an error that occur while running a migration.
func (service DisplayService) DisplayGeneralError(err error) {
	service.DisplayErrorWithMessage(err, "an error occur while running migrations")
}
