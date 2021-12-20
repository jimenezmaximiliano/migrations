package services

import (
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/jimenezmaximiliano/migrations/adapters"
	"github.com/jimenezmaximiliano/migrations/helpers"
)

// Arguments represents the command line arguments for the migrations commands.
type Arguments struct {
	MigrationsPath string
	MigrationName  string
	Command        string
}

// CommandArgument is the API to handle command arguments.
type CommandArgument interface {
	ParseAndValidate() (Arguments, bool)
}

// CommandArgumentService is an implementation of CommandArgument.
type CommandArgumentService struct {
	displayService Display
	parser         adapters.ArgumentParser
}

var _ CommandArgument = CommandArgumentService{}

// NewCommandArgumentService creates a new CommandArgumentService.
func NewCommandArgumentService(displayService Display, parser adapters.ArgumentParser) CommandArgumentService {
	return CommandArgumentService{
		displayService: displayService,
		parser:         parser,
	}
}

// ParseAndValidate parses command line arguments and validates them. In case the validation fails, it'll return true
// as the second returned value.
func (service CommandArgumentService) ParseAndValidate() (Arguments, bool) {
	args := service.parse()

	if args.Command != "migrate" {
		service.displayService.DisplayError(errors.Errorf("invalid 'command' argument: [%s]", args.Command))
		service.displayService.DisplayHelp()
		return args, false
	}

	if args.MigrationsPath == "" {
		service.displayService.DisplayError(errors.New("missing 'path' option for command 'migrate'"))
		service.displayService.DisplayHelp()
		return args, false
	}

	return args, true
}

func (service CommandArgumentService) parse() Arguments {
	var parsedArgs Arguments
	rawArgs := getRearrangedArguments()

	pathOption := service.parser.OptionString("path", "")
	nameOption := service.parser.OptionString("name", "")

	err := service.parser.ParseArguments(rawArgs)
	if err != nil {
		// Let it continue so we can apply default values.
	}

	// Migration's directory path
	if pathOption != nil {
		parsedArgs.MigrationsPath = *pathOption
	}
	if parsedArgs.MigrationsPath != "" {
		parsedArgs.MigrationsPath = helpers.AddTrailingSlashToPathIfNeeded(parsedArgs.MigrationsPath)
	}

	// Migration name
	if nameOption != nil {
		parsedArgs.MigrationName = *nameOption
	}

	// Command
	positionalArguments := service.parser.PositionalArguments()
	if len(positionalArguments) > 0 {
		parsedArgs.Command = positionalArguments[0]
	}
	// Default to migrate for retro compatibility.
	if parsedArgs.Command == "" {
		parsedArgs.Command = "migrate"
	}

	return parsedArgs
}

func getRearrangedArguments() []string {
	args := os.Args[1:]
	rearrangedArgs := make([]string, len(args))
	options := make([]string, 0)
	nonOptions := make([]string, 0)

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			options = append(options, arg)
			continue
		}

		nonOptions = append(nonOptions, arg)
	}

	var index int
	for _, option := range options {
		rearrangedArgs[index] = option
		index++
	}

	for _, nonOption := range nonOptions {
		rearrangedArgs[index] = nonOption
		index++
	}

	return rearrangedArgs
}
