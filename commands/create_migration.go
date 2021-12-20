package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/jimenezmaximiliano/migrations/repositories"
	"github.com/jimenezmaximiliano/migrations/services"
)

type CreateMigration struct {
	fileRepo repositories.FileRepository
	display  services.Display
	args     services.Arguments
}

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

type Command interface {
	Run()
}

var _ Command = CreateMigration{}

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
