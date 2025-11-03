package internal

import (
	"fmt"

	"github.com/spf13/cobra"
)


func MvCommand() *cobra.Command {
	return &cobra.Command{
		Use: "mv",
		Args:  cobra.ExactArgs(2),
		Short: "Move objects from/into an object storage bucket",
		Example: "linode-objim mv <source> <destination>. \n\n<source> and <destination> can be either local paths or object storage paths.\n\nObject storage paths should be in the format: obj:<bucket-name>/<object-key>",
		Run: func(cmd *cobra.Command, args [] string) {
			if len(args) >= 1 {
				var bucket = args[0]
				fmt.Println("Copying objects into bucket: ", bucket)
			}
	
		},
	}

}
