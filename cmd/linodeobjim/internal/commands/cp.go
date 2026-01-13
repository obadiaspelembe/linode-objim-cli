package commands

import (
	"fmt"
 
	"strings"

	"github.com/obadiaspelembe/linode-objim-cli/cmd/linodeobjim/internal/commands/api"
	"github.com/obadiaspelembe/linode-objim-cli/cmd/linodeobjim/internal/commons"

	"github.com/spf13/cobra"
)


func CpCommand() *cobra.Command {
	return &cobra.Command{
		Use: "cp",
		Args:  cobra.ExactArgs(2),
		Short: "cp objects from/into an object storage bucket",
		Example: "linode-objim cp <source> <destination>. \n\n<source> and <destination> can be either local paths or object storage paths.\n\nObject storage paths should be in the format: obj:<bucket-name>/<object-key>",
		Run: func(cmd *cobra.Command, args [] string) {
				if len(args) == 2 {
				var source = args[0]
				var destination = args[1]
				fmt.Println("Copying objects from: ", source, " to: ", destination)

				config, err := commons.LoadConfig("default")
				if err != nil {
					fmt.Println("Error loading config:", err)
					return
				}

				if strings.HasPrefix(source, "obj:") {
					// Handle object storage source

					parts := strings.Split(source[4:], "/")

					fmt.Println("Object Storage Source Parts: ", parts)
					if len(parts) >= 2 {
						bucket := parts[0]
						objectKey := strings.ReplaceAll(source[4:], bucket, "")

						fmt.Println(objectKey)

						SignedURL, err := api.GetLinodeSignedURL(config.Token, config.Region, bucket, objectKey, "text/plain")

						if err != nil {
							fmt.Println("Error getting signed URL:", err)
							return
						}

						err = api.SaveToLocalFile(SignedURL, "."+ objectKey)

						if err != nil {
							fmt.Println(err.Error())
						}
						fmt.Println("Signed URL: ", objectKey)
					}
				} else {
					// Handle local file source
					fmt.Println("Source is a local file path")
				}

			}
	
		},
	}

}
