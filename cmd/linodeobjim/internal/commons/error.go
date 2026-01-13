package commons

import (
	"fmt"

	"github.com/fatih/color"
)

func ErrorCheck(err error, message string) string {

	if err == nil {
		return ""
	}

	red := color.New(color.FgRed).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()
	return fmt.Sprintf("%s : %s", red(message) , red(bold(err.Error())))
}