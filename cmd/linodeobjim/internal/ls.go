package internal

import (
	"fmt"

	"github.com/spf13/cobra"
)


func LsCommand() *cobra.Command {
	return &cobra.Command{
		Use: "ls",
		Args:  cobra.ExactArgs(1),
		Short: "List objects in a object storage bucket",
		Example: "linode-objim ls <bucket-name>",
		Run: func(cmd *cobra.Command, args [] string) {
			if len(args) >= 1 {
				var bucket = args[0]
				fmt.Println("Listing objects in bucket: ", bucket)
			}
	
		},
	}

}
