package commons

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

type Config struct {
	Token  string
	Region string
}

func LoadConfig(profile string) (Config, error) {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("failed to get user home directory: %w", err)
	}

	cfg, err := ini.Load(filepath.Join(homeDir, ".linodeobjim", "credentials.ini"))
	if err != nil {
		return Config{}, fmt.Errorf("failed to load configuration file: %w", err)
	}

	sec := cfg.Section(profile)
	if sec == nil {
		return Config{}, fmt.Errorf("profile %q not found in configuration file", profile)
	}

	token := sec.Key("token").String()
	if token == "" {
		return Config{}, fmt.Errorf("token not found in profile %q", profile)
	}
	region := sec.Key("region").String()
	if region == "" {
		return Config{}, fmt.Errorf("region not found in profile %q", profile)
	}

	return Config{
		Token:  token,
		Region: region,
	}, nil
}
