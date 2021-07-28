package services

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/jimenezmaximiliano/migrations/mocks"
)

func TestParsingArguments(test *testing.T) {
	path := "/tmp"
	name := "name"
	parser := &mocks.ArgumentParser{}
	parser.On("OptionString", "path", mock.AnythingOfType("string")).
		Return(&path)
	parser.On("OptionString", "name", mock.AnythingOfType("string")).
		Return(&name)
	parser.On("PositionalArguments").
		Return([]string{"command"})
	parser.On("Parse").Return(nil)
	service := NewCommandService(parser)
	arguments := service.ParseArguments()

	if arguments.MigrationsPath != "/tmp/" {
		test.Fail()
	}

	if arguments.MigrationName != "name" {
		test.Fail()
	}

	if arguments.Command != "command" {
		test.Fail()
	}
}

func TestParsingArgumentsWithEmptyPath(test *testing.T) {
	path := ""
	name := "name"
	parser := &mocks.ArgumentParser{}
	parser.On("OptionString", "path", mock.AnythingOfType("string")).
		Return(&path)
	parser.On("OptionString", "name", mock.AnythingOfType("string")).
		Return(&name)
	parser.On("PositionalArguments").
		Return([]string{"command"})
	parser.On("Parse").Return(nil)
	service := NewCommandService(parser)
	arguments := service.ParseArguments()

	if arguments.MigrationsPath != "" {
		test.Fail()
	}
}

func TestParsingArgumentsWithEmptyName(test *testing.T) {
	path := "/tmp"
	name := ""
	parser := &mocks.ArgumentParser{}
	parser.On("OptionString", "path", mock.AnythingOfType("string")).
		Return(&path)
	parser.On("OptionString", "name", mock.AnythingOfType("string")).
		Return(&name)
	parser.On("PositionalArguments").
		Return([]string{"command"})
	parser.On("Parse").Return(nil)
	service := NewCommandService(parser)
	arguments := service.ParseArguments()

	if arguments.MigrationName != "" {
		test.Fail()
	}
}

func TestParsingArgumentsWithEmptyCommand(test *testing.T) {
	path := "/tmp"
	name := "name"
	parser := &mocks.ArgumentParser{}
	parser.On("OptionString", "path", mock.AnythingOfType("string")).
		Return(&path)
	parser.On("OptionString", "name", mock.AnythingOfType("string")).
		Return(&name)
	parser.On("PositionalArguments").
		Return([]string{})
	parser.On("Parse").Return(nil)
	service := NewCommandService(parser)
	arguments := service.ParseArguments()

	if arguments.Command != "" {
		test.Fail()
	}
}