package printer

import (
	"encoding/json"
	"fmt"
	"io"
)

// Printer knows how to print results or errors
type Printer interface {
	// Print prints the given object
	Print(object interface{})
}

// printer is the simple implementation of Printer interface
type printer struct {
	out io.Writer
}

// NewStandardOutputPrinter is the constructor of printer
func NewStandardOutputPrinter(out io.Writer) Printer {
	return &printer{
		out: out,
	}
}

// Print implements Printer
func (p *printer) Print(object interface{}) {
	p.printPrettyJSON(object)
}

// printPrettyJSON pretty prints the given object as JSON
func (p *printer) printPrettyJSON(object interface{}) {
	encoder := json.NewEncoder(p.out)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(object); err != nil {
		fmt.Fprintln(p.out, "There was an error but it cannot be marshalled.")
	}
}
