package commons

import (
	"fmt"

	"github.com/fatih/color"
)

type Printer struct {
}

func NewPrinter() *Printer {
	return &Printer{}
}

func (printer *Printer) Success(content string) {
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("%s", green(content))
}

func (printer *Printer) SuccessEx(content string, exclude string) {
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("%s %s",exclude, green(content))
}

func (printer *Printer) Info(content string ) {
	cyan := color.New(color.FgCyan).SprintFunc()
	fmt.Printf("%s", cyan(content))
}

func (printer *Printer) InfoEx(content string, exclude string ) {
	cyan := color.New(color.FgCyan).SprintFunc()
	fmt.Printf("%s %s",exclude, cyan(content))
}

func (printer *Printer) Error(content string) {
	red := color.New(color.FgRed).SprintFunc()
	fmt.Printf("%s", red(content))
}



