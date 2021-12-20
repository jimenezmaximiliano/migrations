package adapters_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jimenezmaximiliano/migrations/adapters"
)

func TestParsingPositionalArgument(test *testing.T) {
	test.Parallel()

	parser := adapters.NewArgumentParser()

	err := parser.ParseArguments([]string{"myArg1", "myArg2"})
	require.Nil(test, err)

	positionalArgs := parser.PositionalArguments()

	require.Len(test, positionalArgs, 2)
	assert.Equal(test, "myArg1", positionalArgs[0])
	assert.Equal(test, "myArg2", positionalArgs[1])
}

func TestParsingPositionalArgumentWithDeprecatedParseMethod(test *testing.T) {
	test.Parallel()

	parser := adapters.NewArgumentParser()

	os.Args = []string{"", "myArg1", "myArg2"}
	err := parser.Parse()
	require.Nil(test, err)

	positionalArgs := parser.PositionalArguments()

	require.Len(test, positionalArgs, 2)
	assert.Equal(test, "myArg1", positionalArgs[0])
	assert.Equal(test, "myArg2", positionalArgs[1])
}

func TestParsingOptions(test *testing.T) {
	test.Parallel()

	parser := adapters.NewArgumentParser()

	opt1 := parser.OptionString("opt1", "")
	opt2 := parser.OptionString("opt2", "")

	err := parser.ParseArguments([]string{"--opt1=1", "--opt2=2"})
	require.Nil(test, err)

	assert.Equal(test, "1", *opt1)
	assert.Equal(test, "2", *opt2)
}
