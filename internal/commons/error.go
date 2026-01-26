package commons

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func ErrorCheck(err error, message string, silent bool)  {

	if err != nil && !silent {
		red := color.New(color.FgRed).SprintFunc()
		bold := color.New(color.Bold).SprintFunc()
		fmt.Printf("%s : %s", red(message), red(bold(err.Error())))

		os.Exit(1)
		return
	}
}
