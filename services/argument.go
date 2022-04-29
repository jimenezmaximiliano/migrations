package services

import (
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/jimenezmaximiliano/migrations/adapters"
	"github.com/jimenezmaximiliano/migrations/helpers"
)

const (
	EnvVarMigrationsPath   string = "MIGRATIONS_PATH"
	EnvVarNewMigrationName string = "MIGRATIONS_NEW_MIGRATION_NAME"
	EnvVarCommand          string = "MIGRATIONS_COMMAND"
)

var ValidCommands = []string{"migrate", "create"}

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

	if !isCommandValid(args.Command) {
		service.displayService.DisplayError(errors.Errorf("invalid 'command' argument: [%s]", args.Command))
		service.displayService.DisplayHelp()
		return args, false
	}

	if args.MigrationsPath == "" {
		service.displayService.DisplayError(errors.Errorf("missing 'path' option for command '%s'", args.Command))
		service.displayService.DisplayHelp()
		return args, false
	}

	if args.Command == "create" && args.MigrationName == "" {
		service.displayService.DisplayError(errors.Errorf("missing 'name' option for command '%s'", args.Command))
		service.displayService.DisplayHelp()
		return args, false
	}

	return args, true
}

func (service CommandArgumentService) parse() Arguments {
	rawArgs := getRearrangedArguments()

	pathOption := service.parser.OptionString("path", "")
	nameOption := service.parser.OptionString("name", "")

	// Parse command line arguments.
	err := service.parser.ParseArguments(rawArgs)
	// nolint - Ignore this empty branch because it serves as documentation of the decision, and it's explicit.
	if err != nil {
		// Let it continue, so we can check environment variables and apply default values.
	}

	return Arguments{
		MigrationsPath: parseMigrationsDirectoryPath(pathOption),
		MigrationName:  parseNewMigrationName(nameOption),
		Command:        service.parseCommand(),
	}
}

func parseMigrationsDirectoryPath(pathOption *string) string {
	// Parse the path command option.
	if pathOption != nil && *pathOption != "" {
		return helpers.AddTrailingSlashToPathIfNeeded(*pathOption)
	}

	// Parse the path environment variable.
	pathEnvVar := os.Getenv(EnvVarMigrationsPath)
	if pathEnvVar != "" {
		return helpers.AddTrailingSlashToPathIfNeeded(pathEnvVar)
	}

	return ""
}

func parseNewMigrationName(nameOption *string) string {
	// Parse the name command option.
	if nameOption != nil && *nameOption != "" {
		return *nameOption
	}

	// Parse the name environment variable.
	return os.Getenv(EnvVarNewMigrationName)
}

func (service CommandArgumentService) parseCommand() string {
	// Parse the first argument.
	positionalArguments := service.parser.PositionalArguments()
	if len(positionalArguments) > 0 && positionalArguments[0] != "" {
		return positionalArguments[0]
	}

	commandEnvVar := os.Getenv(EnvVarCommand)
	if commandEnvVar != "" {
		return commandEnvVar
	}

	// Default to migrate for retro compatibility.
	return "migrate"
}

func getRearrangedArguments() []string {
	if len(os.Args) == 0 {
		return nil
	}

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

func isCommandValid(command string) bool {
	for _, validCommand := range ValidCommands {
		if validCommand == command {
			return true
		}
	}

	return false
}
