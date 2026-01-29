package commands

import (
	// "bufio"
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	// "strings"

	"golang.org/x/term"
	"gopkg.in/ini.v1"

	"github.com/fatih/color"
	"github.com/obadiaspelembe/linode-objim-cli/internal/commons"
	"github.com/spf13/cobra"
)

func ConfigureCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "configure",
		Short:   "Configures the CLI with your Linode Object Storage credentials and preferences",
		Example: "linode-objim configure",
		Run: func(cmd *cobra.Command, args []string) {

			profile, _ := cmd.Flags().GetString("profile")

			green := color.New(color.FgGreen).SprintFunc()
			red := color.New(color.FgRed).SprintFunc()
			bold := color.New(color.Bold).SprintFunc()

			homeDir, err := os.UserHomeDir() 

			credentialsFile := filepath.Join(homeDir, commons.CONFIG_DIR_NAME, commons.CREDENTIALS_FILE)
 
			if err != nil {
				fmt.Printf("Failed to get user home directory: %v\n", err)
				
			}

			if err := os.MkdirAll(filepath.Join(homeDir, commons.CONFIG_DIR_NAME), 0700); err != nil {

				fmt.Fprintln(os.Stderr, red("\nCannot create directory: "+bold(err.Error())))
				os.Exit(1)

			}

			cfg, err := ini.Load(credentialsFile)
			if err != nil {
				cfg = ini.Empty()

			}

			regionMsg := "Region: "
			tokenMsg := "Token: "

			readFromCache := false

			sec, err := cfg.GetSection(profile)

			if err != nil {
				sec, err = cfg.NewSection(profile)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				readFromCache = true
				regionMsg = "Region["+sec.Key("region").String()+"]: "
				tokenMsg = "Token [******" + sec.Key("token").String()[len(sec.Key("token").String())-6:] + "]: "
			}

			reader := bufio.NewReader(os.Stdin)
			fmt.Print(regionMsg)
			var newRegion string 

			newRegion, _ = reader.ReadString('\n')
    		newRegion = strings.TrimSpace(newRegion)

			if newRegion == "" && readFromCache {
				newRegion = sec.Key("region").String()
			}

			fmt.Print(tokenMsg)

			bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				panic(err)
			}

			secret := string(bytePassword)

			if secret == "" && readFromCache {
				secret = sec.Key("token").String()
			}

			sec.Key("token").SetValue(secret)
			sec.Key("region").SetValue(newRegion)

			if err := cfg.SaveTo(credentialsFile); err != nil {
				fmt.Fprintln(os.Stderr, red("\nFailed to save configuration: "+bold(err.Error())))
				os.Exit(1)
			}

			fmt.Println(green("\nSuccess!"))

		},
	}

}
