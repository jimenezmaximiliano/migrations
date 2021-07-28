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
	Parse() error
}

// flagArgumentParser is an implementation of ArgumentParser using the package flag.
type flagArgumentParser struct{
	flagSet *flag.FlagSet
}

// Ensure flagArgumentParser implements ArgumentParser.
var _ ArgumentParser = flagArgumentParser{}

func NewArgumentParser() ArgumentParser {
	flagSet := flag.NewFlagSet("flags", flag.ContinueOnError)
	flagSet.SetOutput(ioutil.Discard)

	return flagArgumentParser{
		flagSet: flagSet,
	}
}

// OptionString defines a string flag with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the flag.
func (adapter flagArgumentParser) OptionString(name string, value string) *string {
	return adapter.flagSet.String(name, value, "")
}

// Parse parses the command-line flags from os.Args[1:]. Must be called
// after all flags are defined and before flags are accessed by the program.
func (adapter flagArgumentParser) Parse() error {
	return adapter.flagSet.Parse(os.Args[1:])
}

// PositionalArguments returns the non-flag command-line arguments.
func (adapter flagArgumentParser) PositionalArguments() []string {
	return adapter.flagSet.Args()
}
