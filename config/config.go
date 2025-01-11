package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

type ConfigHomeAssistant struct {
	URL   string `mapstructure:"url"`
	Token string `mapstructure:"token"`
}

type Config struct {
	HomeAssistant ConfigHomeAssistant `mapstructure:"homeassistant"`
}

func Load() (*Config, error) {
	viper.AutomaticEnv()

	viper.SetConfigName("config.yml")
	viper.SetConfigType("yaml")

	// (Cross platform) default config configDirPath (~/.config/go-automate or %APPDATA%\go-automate)
	configDirPath := ""
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		configDirPath = os.Getenv("XDG_CONFIG_HOME") + "/go-automate"
	} else if os.Getenv("APPDATA") != "" {
		configDirPath = os.Getenv("APPDATA") + "/go-automate"
	} else if os.Getenv("HOME") != "" {
		configDirPath = os.Getenv("HOME") + "/.config/go-automate"
	} else {
		return nil, fmt.Errorf("Could not determine config path")
	}

	// Create the config directory if it doesn't exist
	os.MkdirAll(configDirPath, 0755)
	os.WriteFile(configDirPath+"/config.yml", []byte{}, 0644)
	viper.AddConfigPath(configDirPath)

	// Set default values
	viper.SetDefault("homeassistant.url", "http://homeassisant.local:8123")
	viper.SetDefault("homeassistant.token", "")

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &cfg, nil
}

func (cfg *Config) Save() error {
	viper.Set("homeassistant.url", cfg.HomeAssistant.URL)
	viper.Set("homeassistant.token", cfg.HomeAssistant.Token)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}
	return nil
}

func (cfg *Config) Setup() (*Config, error) {
	if cfg.HomeAssistant.URL != "" && cfg.HomeAssistant.Token != "" {
		log.Info("------ Setup ------")

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("What is your Home Assistant URL?").
					Value(&cfg.HomeAssistant.URL).
					Validate(func(s string) error {
						if s == "" {
							return errors.New("Please enter a valid URL")
						}
						if s[:7] != "http://" && s[:8] != "https://" {
							return errors.New("Please enter a valid URL")
						}
						return nil
					}),

				huh.NewInput().
					Title("What is your Home Assistant Token (Long-Lived Access Token)?").
					Description("You can create a Long-Lived Access Token in your Home Assistant profile.").
					Value(&cfg.HomeAssistant.Token).
					Validate(func(s string) error {
						if s == "" {
							return errors.New("Please enter a valid token")
						}
						return nil
					}),
			),
		)
		if err := form.Run(); err != nil {
			return nil, err
		}

		if err := cfg.Save(); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}
