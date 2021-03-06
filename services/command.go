package services

import (
	"github.com/jimenezmaximiliano/migrations/adapters"
	"github.com/jimenezmaximiliano/migrations/helpers"
)

// Arguments represents the command line arguments for the migrations command.
type Arguments struct {
	MigrationsPath string
}

// Command parses command line arguments.
type Command interface {
	ParseArguments() Arguments
}

type commandService struct {
	optionParser adapters.OptionParser
}

// Ensure commandService implements Command.
var _ Command = commandService{}

// NewCommandService returns an implementation of Command.
func NewCommandService(optionParser adapters.OptionParser) Command {
	return commandService{
		optionParser: optionParser,
	}
}

// ParseArguments parses command line arguments.
func (service commandService) ParseArguments() Arguments {
	path := service.optionParser.String("path", "", "")
	service.optionParser.Parse()
	resultPath := *path

	if resultPath != "" {
		resultPath = helpers.AddTrailingSlashToPathIfNeeded(resultPath)
	}

	return Arguments{
		MigrationsPath: resultPath,
	}
}
