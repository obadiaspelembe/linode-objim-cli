package main

import (
	"fmt"
	"os"

	"github.com/obadiaspelembe/linode-objim-cli/cmd/linodeobjim/internal"
	"github.com/spf13/cobra"
)

const SHORT_DESCRIPTION = "linode-objim is an improved version of linode object storage managemement cli tool."

var version = "0.0.1-alpha"
var rootCommand = &cobra.Command{
	Use:     "linode-objim",
	Version: version,
	Short:   SHORT_DESCRIPTION,
	Long:    SHORT_DESCRIPTION,
}

func Execute() { 
	
	cpCmd := internal.CpCommand()
	mvCmd := internal.MvCommand()
	rmCmd := internal.RmCommand()
	var recursive bool

	cpCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Recursively copy objects")
	mvCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Recursively move objects")
	rmCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Recursively remove objects")

	rootCommand.AddCommand(internal.LsCommand())
	rootCommand.AddCommand(cpCmd)
	rootCommand.AddCommand(mvCmd)
	rootCommand.AddCommand(rmCmd)

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