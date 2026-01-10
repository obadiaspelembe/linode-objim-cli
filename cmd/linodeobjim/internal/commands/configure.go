package commands

import (
	// "bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	// "strings"

	"golang.org/x/term"
	"gopkg.in/ini.v1"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func ConfigureCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "configure",
		Short:   "Configures the CLI with your Linode Object Storage credentials and preferences",
		Example: "linode-objim configure",
		Run: func(cmd *cobra.Command, args []string) {

			green := color.New(color.FgGreen).SprintFunc()
			red := color.New(color.FgRed).SprintFunc()
			bold := color.New(color.Bold).SprintFunc()

			homeDir, err := os.UserHomeDir()

			if err != nil {
				fmt.Printf("Failed to get user home directory: %v\n", err)
				
			}

			if err := os.MkdirAll(filepath.Join(homeDir, ".linodeobjim"), 0700); err != nil {

				fmt.Fprintln(os.Stderr, red("\nCannot create directory: "+bold(err.Error())))
				os.Exit(1)

			}

			cfg, err := ini.Load(filepath.Join(homeDir, ".linodeobjim", "credentials.ini"))
			if err != nil {
				cfg = ini.Empty()

			}

			fmt.Print("Token: ")

			bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				panic(err)
			}

			secret := string(bytePassword)

			// Create a section and set keys
			sec := cfg.Section("default")

			if sec == nil {
				sec, err = cfg.NewSection("default")
				if err != nil {
					log.Fatal(err)
				}
			}

			sec.Key("token").SetValue(secret)

			if err := cfg.SaveTo(filepath.Join(homeDir, ".linodeobjim", "credentials.ini")); err != nil {
				fmt.Fprintln(os.Stderr, red("\nFailed to save configuration: "+bold(err.Error())))
				os.Exit(1)
			}

			fmt.Println(green("\nSuccess!"))

		},
	}

}
