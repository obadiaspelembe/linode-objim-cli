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
		Use:     "cp",
		Args:    cobra.ExactArgs(2),
		Short:   "cp objects from/into an object storage bucket",
		Example: "linode-objim cp <source> <destination>. \n\n<source> and <destination> can be either local paths or object storage paths.\n\nObject storage paths should be in the format: obj:<bucket-name>/<object-key>",
		Run: func(cmd *cobra.Command, args []string) {

			printer := commons.NewPrinter()
			rec, _ := cmd.Flags().GetBool("recursive")

			if len(args) == 2 {
				var source = args[0]

				config, err := commons.LoadConfig("default")

				commons.ErrorCheck(err, "Failed to load configuration")

				if strings.HasPrefix(source, commons.BUCKET_PREFIX) {

					parts := strings.Split(source[len(commons.BUCKET_PREFIX):], "/")

					if len(parts) >= 2 {
						bucket := parts[0]
						if rec {

							if parts[1] == "" {
								result := api.GetObjectList(bucket, config.Region, config.Token)

								for _, sObjec := range result.Data {
									status := commons.ProcessBucketObject(config, bucket, sObjec.Name)

									if status {
										printer.Success(fmt.Sprintf("%s\n", sObjec.Name))
									}
								} 
							} else {
								

								result := api.GetObjectList(bucket, config.Region, config.Token)

								for _, sObjec := range result.Data {

									isPartOf := strings.HasPrefix(sObjec.Name, source[len(parts[0]) + len(commons.BUCKET_PREFIX) + 1 :])

									if isPartOf {
										status := commons.ProcessBucketObject(config, bucket, sObjec.Name)
	
										if status {
											printer.Success(fmt.Sprintf("%s\n", sObjec.Name))
										}
									}
								}
							}
 

							return
						}

						objectKey := strings.ReplaceAll(source[len(commons.BUCKET_PREFIX):], bucket, "")
						commons.ProcessBucketObject(config, bucket, objectKey)
					}

					printer.Success("Success!")
				} else {
					fmt.Println("Source is a local file path")
				}

			}

		},
	}

}
