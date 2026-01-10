package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	// "path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/obadiaspelembe/linode-objim-cli/cmd/linodeobjim/internal/commons"
	"github.com/obadiaspelembe/linode-objim-cli/cmd/linodeobjim/internal/linode"
	"github.com/spf13/cobra"
)

type Object struct {
	Name         string `json:"name"`
	Size         int    `json:"size"`
	LastModified string `json:"last_modified"`
	ETag         string `json:"etag"`
}

type ObjectListResponse struct {
	Data        []Object `json:"data"`
	NextMarker  string   `json:"next_marker"`
	IsTruncated bool     `json:"is_truncated"`
}

func LsCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "ls",
		Args:    cobra.ExactArgs(1),
		Short:   "List objects in a object storage bucket",
		Example: "linode-objim ls <bucket-name>",
		Run: func(cmd *cobra.Command, args []string) {

			green := color.New(color.FgGreen).SprintFunc()
			red := color.New(color.FgRed).SprintFunc()
			cyan := color.New(color.FgCyan).SprintFunc()
			bold := color.New(color.Bold).SprintFunc()
			// magenta := color.New(color.FgMagenta).SprintFunc()

			config, err := commons.LoadConfig("default")
			if err != nil {
				fmt.Println(red("Error loading config: " + bold(err.Error())))
				return
			}

			if len(args) >= 1 {
				url := "https://api.linode.com/v4/object-storage/buckets/" + config.Region + "/" + args[0] + "/object-list"

				apiClient := linode.NewAPI(config.Token, config.Region)

				req := apiClient.InitializeRequest(url)

				res, _ := http.DefaultClient.Do(req)

				defer res.Body.Close()
				body, _ := io.ReadAll(res.Body)

				var result ObjectListResponse
				if err := json.Unmarshal(body, &result); err != nil {
					panic(err)
				}

				fmt.Printf(green("total %-d\n"), len(result.Data))
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
						fmt.Printf("%-4s %8d %-30s %4s \n",
							sObjec.ETag, sObjec.Size, last, cyan(bold(filepath)))

					} else {

						fmt.Printf("%-4s %8d %-45s\n",
							sObjec.ETag, sObjec.Size, sObjec.Name)
					}
				}

			}
		},
	}
}
