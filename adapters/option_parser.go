package adapters

import (
	"flag"
)

// OptionParser parses command line flags.
type OptionParser interface {
	String(name string, value string, usage string) *string
	Parse()
}

// FlagOptionParser is an implementation of OptionParser using the package flag.
type FlagOptionParser struct{}

// Ensure FlagOptionParser implements OptionParser.
var _ OptionParser = FlagOptionParser{}

// String defines a string flag with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the flag.
func (adapter FlagOptionParser) String(name string, value string, usage string) *string {
	return flag.String(name, value, usage)
}

// Parse parses the command-line flags from os.Args[1:]. Must be called
// after all flags are defined and before flags are accessed by the program.
func (adapter FlagOptionParser) Parse() {
	flag.Parse()
}
