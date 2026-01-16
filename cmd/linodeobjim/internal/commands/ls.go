package commands

import (
	"fmt"

	"strings"

	"github.com/obadiaspelembe/linode-objim-cli/cmd/linodeobjim/internal/commands/api"
	"github.com/obadiaspelembe/linode-objim-cli/cmd/linodeobjim/internal/commons"
	"github.com/spf13/cobra"
)

func LsCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "ls",
		Args: cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		Short:   "List objects in a object storage bucket",
		Example: "linode-objim ls <bucket-name>",
		Run: func(cmd *cobra.Command, args []string) {

			printer := commons.NewPrinter()

			config, err := commons.LoadConfig("default")

			commons.ErrorCheck(err, "Failed to load configuration")

			if len(args) >= 1 {
				// Call the API to get the object list
				result := api.GetObjectList(args[0], config.Region, config.Token)

				printer.Success(fmt.Sprintf("total %d\n", len(result.Data)))

				for _, sObjec := range result.Data {

					if strings.Contains(sObjec.Name, "/") {
						parts := strings.Split(sObjec.Name, "/")
						filepath := ""
						last := ""

						for i, part := range parts {

							if (len(parts) - 1) == i {
								last = part
								break
							}

							filepath = filepath + part + "/"

						}

						printer.InfoEx(
							fmt.Sprintf(" %4s \n", filepath),
							fmt.Sprintf("%-4s %8dB %-30s", sObjec.ETag, sObjec.Size, last))

					} else {
						fmt.Printf("%-4s %8dB %-45s\n",
							sObjec.ETag, sObjec.Size, sObjec.Name)
					}
				}

			}
		},
	}
}
