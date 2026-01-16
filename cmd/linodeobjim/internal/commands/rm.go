package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)


func RmCommand() *cobra.Command {
	return &cobra.Command{
		Use: "rm",
		Args:  cobra.ExactArgs(1),
		Short: "Remove objects in a object storage bucket",
		Example: "linode-objim rm <bucket-name>/<object-key>",
		Run: func(cmd *cobra.Command, args [] string) {
			if len(args) >= 1 {
				var bucket = args[0]
				fmt.Println("Listing objects in bucket: ", bucket)
			}
	
		},
	}

}
