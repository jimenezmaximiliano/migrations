package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/jimenezmaximiliano/migrations/repositories"
	"github.com/jimenezmaximiliano/migrations/services"
)

// CreateMigration is a command that creates a migration file.
type CreateMigration struct {
	fileRepo repositories.FileRepository
	display  services.Display
	args     services.Arguments
}

// NewCreateMigrationCommand builds a CreateMigration.
func NewCreateMigrationCommand(
	fileRepo repositories.FileRepository,
	display services.Display,
	args services.Arguments,
) CreateMigration {
	return CreateMigration{
		fileRepo: fileRepo,
		display:  display,
		args:     args,
	}
}

// Command represents a command line command that can be run.
type Command interface {
	Run()
}

var _ Command = CreateMigration{}

// Run creates a migration file with a sample content using a sort of unique name based on a timestamp.
func (command CreateMigration) Run() {
	now := time.Now()
	fileName := fmt.Sprintf("%d_%s", now.UnixNano(), command.args.MigrationName)

	filePath := command.args.MigrationsPath + fileName
	extension := filePath[len(filePath)-4:]
	if extension != ".sql" {
		filePath += ".sql"
	}

	err := command.fileRepo.CreateMigration(filePath, "SELECT 1;")
	if err != nil {
		command.display.DisplayError(err)
		os.Exit(1)
	}

	command.display.DisplayInfo(fmt.Sprintf("migration file created at %s", filePath))
}
