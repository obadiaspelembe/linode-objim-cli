package commands

import (
	"fmt"

	"strings"

	"github.com/obadiaspelembe/linode-objim-cli/internal/commands/api"
	"github.com/obadiaspelembe/linode-objim-cli/internal/commons"
	"github.com/spf13/cobra"
)

func LsCommand() *cobra.Command {
	return &cobra.Command{
		Use:           "ls",
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		Short:         "List objects in a object storage bucket",
		Example:       "linode-objim ls <bucket-name>",
		Run: func(cmd *cobra.Command, args []string) {

			rec, _ := cmd.Flags().GetBool("recursive")

			printer := commons.NewPrinter()

			config, err := commons.LoadConfig("default")

			commons.ErrorCheck(err, "Failed to load configuration", false)

			if len(args) >= 1 {
				// Call the API to get the object list
				result := api.GetObjectList(args[0], config.Region, config.Token, rec)

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

						if !rec {
							printer.InfoEx(
								fmt.Sprintf("%-16s %8s %-4s\n",
							" ", "DIR", sObjec.Name),
							"")
						} else {

							printer.InfoEx(
								fmt.Sprintf(" %2s \n", filepath),
								fmt.Sprintf("%-2s %8dB %-28s", commons.FormatTimeString(sObjec.LastModified), sObjec.Size, last))
						}

					} else {

						fmt.Printf("%-2s %8dB %-4s\n",
							commons.FormatTimeString(sObjec.LastModified), sObjec.Size, sObjec.Name)
					}
				}

			}
		},
	}
}
