package services

import (
	"github.com/jimenezmaximiliano/migrations/migrations/adapters"
	"github.com/jimenezmaximiliano/migrations/migrations/helpers"
)

type Arguments struct {
	MigrationsPath string
}

type Command interface {
	ParseArguments() Arguments
}

func NewCommandService(optionParser adapters.OptionParser) Command {
	return commandService{
		optionParser: optionParser,
	}
}

type commandService struct {
	optionParser adapters.OptionParser
}

func (service commandService) ParseArguments() Arguments {
	path := service.optionParser.String("path", "", "")
	service.optionParser.Parse()
	pathWithTrailingSlash := helpers.AddTrailingSlashToPathIfNeeded(*path)

	return Arguments{
		MigrationsPath: pathWithTrailingSlash,
	}
}
