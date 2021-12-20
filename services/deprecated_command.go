package services

import (
	"github.com/jimenezmaximiliano/migrations/adapters"
)

// Deprecated: use CommandArgument instead.
// Command parses command line arguments.
type Command interface {
	ParseArguments() Arguments
}

type CommandService struct {
	argumentService CommandArgument
}

// Ensure CommandService implements Command.
var _ Command = CommandService{}

// Deprecated: use NewCommandArgumentService instead.
// NewCommandService returns an implementation of Command.
func NewCommandService(argumentParser adapters.ArgumentParser) Command {
	return CommandService{
		argumentService: NewCommandArgumentService(NewDisplayService(adapters.NilPrinterAdapter{}), argumentParser),
	}
}

// Deprecated: use CommandArgument instead.
// ParseArguments parses command line arguments.
func (service CommandService) ParseArguments() Arguments {
	args, ok := service.argumentService.ParseAndValidate()
	if !ok {
		return Arguments{}
	}

	return args
}
