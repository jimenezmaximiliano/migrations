package adapters

import (
	"flag"
	"io/ioutil"
	"os"
)

// ArgumentParser parses command line flags.
type ArgumentParser interface {
	OptionString(name string, value string) *string
	PositionalArguments() []string
	ParseArguments(args []string) error
	Parse() error
}

// FlagArgumentParser is an implementation of ArgumentParser using the package flag.
type FlagArgumentParser struct {
	flagSet *flag.FlagSet
}

// Ensure FlagArgumentParser implements ArgumentParser.
var _ ArgumentParser = FlagArgumentParser{}

func NewArgumentParser() FlagArgumentParser {
	flagSet := flag.NewFlagSet("flags", flag.ContinueOnError)
	flagSet.SetOutput(ioutil.Discard)

	return FlagArgumentParser{
		flagSet: flagSet,
	}
}

// OptionString defines a string flag with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the flag.
func (adapter FlagArgumentParser) OptionString(name string, value string) *string {
	return adapter.flagSet.String(name, value, "")
}

// ParseArguments parses the command-line flags from os.Args[1:]. Must be called
// after all flags are defined and before flags are accessed by the program.
func (adapter FlagArgumentParser) ParseArguments(args []string) error {
	// Default os.Args[1:] for retro compatibility.
	if len(args) == 0 {
		args = os.Args[1:]
	}

	return adapter.flagSet.Parse(args)
}

// Deprecated: use ParseArguments instead.
// Parse is deprecated.
func (adapter FlagArgumentParser) Parse() error {
	return adapter.ParseArguments(nil)
}

// PositionalArguments returns the non-flag command-line arguments.
func (adapter FlagArgumentParser) PositionalArguments() []string {
	return adapter.flagSet.Args()
}
