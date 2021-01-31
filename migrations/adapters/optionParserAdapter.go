package adapters

import "flag"

type OptionParser interface {
	String(name string, value string, usage string) *string
	Parse()
}

type FlagOptionParser struct{}

func (adapter FlagOptionParser) String(name string, value string, usage string) *string {
	return flag.String(name, value, usage)
}

func (adapter FlagOptionParser) Parse() {
	flag.Parse()
}
