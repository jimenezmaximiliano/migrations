package main

import (
	"fmt"
	"os"

	"github.com/jimenezmaximiliano/migrations/adapters"
	"github.com/jimenezmaximiliano/migrations/services"
)

// Deprecated: use the migrations binary to create new migrations instead.
func main() {
	commandService := services.NewCommandService(adapters.NewArgumentParser())
	arguments := commandService.ParseArguments()

	if arguments.Command == "migration:create" {
		if arguments.MigrationName == "" || arguments.MigrationsPath == "" {
			displayHelp()
			os.Exit(1)
		}

		createdMigrationPath, err := createMigration(arguments.MigrationName, arguments.MigrationsPath)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}

		fmt.Fprintf(os.Stdout, "Created migration on %s\n", createdMigrationPath)
		os.Exit(0)
	}

	displayHelp()
	os.Exit(0)
}

// Deprecated: migo is deprecated.
func displayHelp() {
	fmt.Fprint(os.Stdout, "Documentation: https://github.com/jimenezmaximiliano/migrations\n\n")
	fmt.Fprint(os.Stdout, "Usage:\n")
	fmt.Fprint(os.Stdout, "\tmigo -path=/path/to/migrations/directory/ -name=createTableGophers migration:create\n\n")
}
