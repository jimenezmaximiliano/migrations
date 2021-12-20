package services_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/jimenezmaximiliano/migrations/mocks"
	"github.com/jimenezmaximiliano/migrations/services"
)

func TestParsingValidArgs(test *testing.T) {
	display := &mocks.Display{}
	parser := &mocks.ArgumentParser{}

	os.Args = []string{"", "migrate", "-path=/tmp"}
	path := "/tmp"
	name := ""

	parser.On("OptionString", "path", mock.AnythingOfType("string")).
		Return(&path)
	parser.On("OptionString", "name", mock.AnythingOfType("string")).
		Return(&name)
	parser.On("PositionalArguments").
		Return([]string{"migrate"})
	parser.On("ParseArguments", mock.AnythingOfType("[]string")).Return(nil)

	service := services.NewCommandArgumentService(display, parser)

	args, ok := service.ParseAndValidate()

	assert.True(test, ok)
	assert.Equal(test, "/tmp/", args.MigrationsPath)
	assert.Equal(test, "", args.MigrationName)
	assert.Equal(test, "migrate", args.Command)
}

func TestParsingValidArgsInDifferentOrder(test *testing.T) {
	display := &mocks.Display{}
	parser := &mocks.ArgumentParser{}

	os.Args = []string{"", "--path=/tmp", "migrate"}
	path := "/tmp"
	name := ""

	parser.On("OptionString", "path", mock.AnythingOfType("string")).
		Return(&path)
	parser.On("OptionString", "name", mock.AnythingOfType("string")).
		Return(&name)
	parser.On("PositionalArguments").
		Return([]string{"migrate"})
	parser.On("ParseArguments", mock.AnythingOfType("[]string")).Return(nil)

	service := services.NewCommandArgumentService(display, parser)

	args, ok := service.ParseAndValidate()

	assert.True(test, ok)
	assert.Equal(test, "/tmp/", args.MigrationsPath)
	assert.Equal(test, "", args.MigrationName)
	assert.Equal(test, "migrate", args.Command)
}

func TestMigrateIsTheDefaultCommand(test *testing.T) {
	display := &mocks.Display{}
	parser := &mocks.ArgumentParser{}

	os.Args = []string{"", "--path=/tmp"}
	path := "/tmp"
	name := ""

	parser.On("OptionString", "path", mock.AnythingOfType("string")).
		Return(&path)
	parser.On("OptionString", "name", mock.AnythingOfType("string")).
		Return(&name)
	parser.On("PositionalArguments").
		Return([]string{})
	parser.On("ParseArguments", mock.AnythingOfType("[]string")).Return(nil)

	service := services.NewCommandArgumentService(display, parser)

	args, ok := service.ParseAndValidate()

	assert.True(test, ok)
	assert.Equal(test, "/tmp/", args.MigrationsPath)
	assert.Equal(test, "", args.MigrationName)
	assert.Equal(test, "migrate", args.Command)
}

func TestWrongCommand(test *testing.T) {
	display := &mocks.Display{}
	parser := &mocks.ArgumentParser{}

	os.Args = []string{"", "oops"}
	path := "/tmp"
	name := ""

	parser.On("OptionString", "path", mock.AnythingOfType("string")).
		Return(&path)
	parser.On("OptionString", "name", mock.AnythingOfType("string")).
		Return(&name)
	parser.On("PositionalArguments").
		Return([]string{"oops"})
	parser.On("ParseArguments", mock.AnythingOfType("[]string")).Return(nil)

	display.On("DisplayError", mock.MatchedBy(func(err error) bool { return true })).Return(nil)
	display.On("DisplayHelp").Return(nil)

	service := services.NewCommandArgumentService(display, parser)

	_, ok := service.ParseAndValidate()

	assert.False(test, ok)
}

func TestMissingOption(test *testing.T) {
	display := &mocks.Display{}
	parser := &mocks.ArgumentParser{}

	os.Args = []string{"", "migrate"}
	path := ""
	name := ""

	parser.On("OptionString", "path", mock.AnythingOfType("string")).
		Return(&path)
	parser.On("OptionString", "name", mock.AnythingOfType("string")).
		Return(&name)
	parser.On("PositionalArguments").
		Return([]string{"migrate"})
	parser.On("ParseArguments", mock.AnythingOfType("[]string")).Return(nil)

	display.On("DisplayError", mock.MatchedBy(func(err error) bool { return true })).Return(nil)
	display.On("DisplayHelp").Return(nil)

	service := services.NewCommandArgumentService(display, parser)

	_, ok := service.ParseAndValidate()

	assert.False(test, ok)
}
