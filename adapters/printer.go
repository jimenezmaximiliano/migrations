package adapters

import (
	"fmt"
	"io"
)

// Printer formats and prints messages to the given writer.
type Printer interface {
	Print(writer io.Writer, format string, a ...interface{}) error
}

// PrinterAdapter is an implementation of Printer.
type PrinterAdapter struct{}

var _ Printer = PrinterAdapter{}

// Print outputs a string given a format.
func (adapter PrinterAdapter) Print(writer io.Writer, format string, a ...interface{}) error {
	_, err := fmt.Fprintf(writer, format, a...)

	return err
}

type NilPrinterAdapter struct{}

// Print outputs a string given a format.
func (adapter NilPrinterAdapter) Print(writer io.Writer, format string, a ...interface{}) error {
	return nil
}
