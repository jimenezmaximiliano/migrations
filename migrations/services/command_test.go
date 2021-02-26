package services

import (
	"testing"

	"github.com/jimenezmaximiliano/migrations/mocks"
	"github.com/stretchr/testify/mock"
)

func TestParsingArguments(test *testing.T) {
	path := "/tmp"
	parser := &mocks.OptionParser{}
	parser.On("String", "path", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(&path)
	parser.On("Parse")
	service := NewCommandService(parser)
	arguments := service.ParseArguments()

	if arguments.MigrationsPath != "/tmp/" {
		test.Fail()
	}
}

func TestParsingArgumentsWithEmptyPath(test *testing.T) {
	path := ""
	parser := &mocks.OptionParser{}
	parser.On("String", "path", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(&path)
	parser.On("Parse")
	service := NewCommandService(parser)
	arguments := service.ParseArguments()

	if arguments.MigrationsPath != "/" {
		test.Fail()
	}
}