package services

import (
	"github.com/jimenezmaximiliano/migrations/adapters"
	"github.com/jimenezmaximiliano/migrations/helpers"
)

// Arguments represents the command line arguments for the migrations commands.
type Arguments struct {
	MigrationsPath string
	MigrationName string
	Command string
}

// Command parses command line arguments.
type Command interface {
	ParseArguments() Arguments
}

type commandService struct {
	argumentParser adapters.ArgumentParser
}

// Ensure commandService implements Command.
var _ Command = commandService{}

// NewCommandService returns an implementation of Command.
func NewCommandService(argumentParser adapters.ArgumentParser) Command {
	return commandService{
		argumentParser: argumentParser,
	}
}

// ParseArguments parses command line arguments.
func (service commandService) ParseArguments() Arguments {
	pathOption := service.argumentParser.OptionString("path", "")
	nameOption := service.argumentParser.OptionString("name", "")
	err := service.argumentParser.Parse()
	if err != nil {
		return Arguments{}
	}

	// Migration's directory path
	parsedPathOption := *pathOption
	if parsedPathOption != "" {
		parsedPathOption = helpers.AddTrailingSlashToPathIfNeeded(parsedPathOption)
	}

	// Migration name
	parsedNameOption := *nameOption

	positionalArguments := service.argumentParser.PositionalArguments()

	var command string
	if len(positionalArguments) > 0 {
		command = positionalArguments[0]
	}

	return Arguments{
		MigrationsPath: parsedPathOption,
		MigrationName:  parsedNameOption,
		Command: command,
	}
}
