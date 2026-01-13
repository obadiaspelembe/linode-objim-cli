package commands

import (
	"fmt"

	"strings"

	"github.com/fatih/color"
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

				green := color.New(color.FgGreen).SprintFunc()
				red := color.New(color.FgRed).SprintFunc()
				// cyan := color.New(color.FgCyan).SprintFunc()
				// bold := color.New(color.Bold).SprintFunc()
				if len(args) == 2 {
				var source = args[0]
				// var destination = args[1]

				config, err := commons.LoadConfig("default")
				if err != nil {
					fmt.Println(red("Error loading config: " + err.Error()))
					return
				}

				if strings.HasPrefix(source, commons.BUCKET_PREFIX) {
					

					parts := strings.Split(source[len(commons.BUCKET_PREFIX):], "/") 

					if len(parts) >= 2 {
						bucket := parts[0]
						objectKey := strings.ReplaceAll(source[len(commons.BUCKET_PREFIX):], bucket, "")
						SignedURL, err := api.GetLinodeSignedURL(config.Token, config.Region, bucket, objectKey, "text/plain")

						if err != nil {
							fmt.Println(red("Error getting signed URL " + err.Error()))
							return
						}

						err = api.SaveToLocalFile(SignedURL, "."+ objectKey)

						if err != nil {
							fmt.Println(red(err.Error()))
						}
					}

					fmt.Println(green("Success!"))
				} else {
					fmt.Println("Source is a local file path")
				}

			}
	
		},
	}

}
