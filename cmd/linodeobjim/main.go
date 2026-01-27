package main

import (
	"fmt"
	"os"

	"github.com/obadiaspelembe/linode-objim-cli/internal/commands"
	"github.com/spf13/cobra"
)

const SHORT_DESCRIPTION = "linode-objl+ is an improved version of linode object storage managemement cli tool."

var version = "0.0.1-alpha"
var rootCommand = &cobra.Command{
	Use:     "linode-objl+",
	Version: version,
	Short:   SHORT_DESCRIPTION,
	Long:    SHORT_DESCRIPTION,
}

func Execute() { 
	
	cpCmd := commands.CpCommand() 
	rmCmd := commands.RmCommand()
	configureCmd := commands.ConfigureCommand()
	lsCmd := commands.LsCommand()
	var recursive bool

	cpCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Copy from/to bucket directories and their contents recursively") 
	rmCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Remove bucket directories and their contents recursively")
	lsCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "List bucket subdirectories recursively")

	rootCommand.AddCommand(lsCmd)
	rootCommand.AddCommand(cpCmd)
	rootCommand.AddCommand(rmCmd)
	rootCommand.AddCommand(configureCmd)

	rootCommand.CompletionOptions = cobra.CompletionOptions{DisableDefaultCmd: true} 

	// rootCommand.AddCommand(pushCommand) 
	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func main() {
	Execute()
}