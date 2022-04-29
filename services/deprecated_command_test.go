package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/jimenezmaximiliano/migrations/mocks"
	"github.com/jimenezmaximiliano/migrations/services"
)

func TestParsingArgumentsWithEmptyPath(test *testing.T) {
	path := ""
	name := "name"
	parser := &mocks.ArgumentParser{}
	defer parser.AssertExpectations(test)
	parser.On("OptionString", "path", mock.AnythingOfType("string")).
		Return(&path)
	parser.On("OptionString", "name", mock.AnythingOfType("string")).
		Return(&name)
	parser.On("PositionalArguments").
		Return([]string{"command"})
	parser.On("ParseArguments", mock.AnythingOfType("[]string")).Return(nil)

	service := services.NewCommandService(parser)
	arguments := service.ParseArguments()

	assert.Equal(test, "", arguments.MigrationsPath)
}

func TestParsingArgumentsWithEmptyName(test *testing.T) {
	path := "/tmp"
	name := ""
	parser := &mocks.ArgumentParser{}
	defer parser.AssertExpectations(test)
	parser.On("OptionString", "path", mock.AnythingOfType("string")).
		Return(&path)
	parser.On("OptionString", "name", mock.AnythingOfType("string")).
		Return(&name)
	parser.On("PositionalArguments").
		Return([]string{"command"})
	parser.On("ParseArguments", mock.AnythingOfType("[]string")).Return(nil)

	service := services.NewCommandService(parser)
	arguments := service.ParseArguments()

	assert.Equal(test, "", arguments.MigrationName)
}

func TestParsingArgumentsWithEmptyCommand(test *testing.T) {
	path := "/tmp"
	name := "name"
	parser := &mocks.ArgumentParser{}
	parser.AssertExpectations(test)
	parser.On("OptionString", "path", mock.AnythingOfType("string")).
		Return(&path)
	parser.On("OptionString", "name", mock.AnythingOfType("string")).
		Return(&name)
	parser.On("PositionalArguments").
		Return([]string{})
	parser.On("Parse").Return(nil)
	service := services.NewCommandService(parser)
	parser.On("ParseArguments", mock.AnythingOfType("[]string")).Return(nil)

	arguments := service.ParseArguments()

	assert.Equal(test, "migrate", arguments.Command)
}
