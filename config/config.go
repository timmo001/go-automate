package config

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/viper"
	"golang.org/x/term"
)

type ConfigHomeAssistant struct {
	URL   string `mapstructure:"url"`
	Token string `mapstructure:"token"`
}

type Config struct {
	HomeAssistant ConfigHomeAssistant `mapstructure:"homeassistant"`
}

func configDirPath() (string, error) {
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		return os.Getenv("XDG_CONFIG_HOME") + "/go-automate", nil
	}
	if os.Getenv("APPDATA") != "" {
		return os.Getenv("APPDATA") + "/go-automate", nil
	}
	if os.Getenv("HOME") != "" {
		return os.Getenv("HOME") + "/.config/go-automate", nil
	}

	return "", fmt.Errorf("could not determine config path")
}

func configFilePath() (string, error) {
	dir, err := configDirPath()
	if err != nil {
		return "", err
	}

	return dir + "/config.yml", nil
}

func Load() (*Config, error) {
	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// (Cross platform) default config configDirPath (~/.config/go-automate or %APPDATA%\go-automate)
	configDirPath, err := configDirPath()
	if err != nil {
		return nil, err
	}

	// Create the config directory if it doesn't exist
	os.MkdirAll(configDirPath, 0755)
	// os.WriteFile(configDirPath+"/config.yml", []byte{}, 0644)
	viper.AddConfigPath(configDirPath)

	// Set default values
	viper.SetDefault("homeassistant.url", "http://homeassistant.local:8123")
	viper.SetDefault("homeassistant.token", "")

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
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

	path := viper.ConfigFileUsed()
	if path == "" {
		var err error
		path, err = configFilePath()
		if err != nil {
			return err
		}
	}

	viper.SetConfigPermissions(0600)
	if err := viper.WriteConfigAs(path); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}
	return nil
}

func (cfg *Config) Setup(interactive bool) (*Config, error) {
	if cfg.HomeAssistant.Token == "" {
		if !interactive {
			path := viper.ConfigFileUsed()
			if path == "" {
				path, _ = configFilePath()
			}

			return nil, fmt.Errorf(
				"home assistant token is not configured; configure %s or run go-automate in an interactive terminal once",
				path,
			)
		}

		slog.Info("------ Setup ------")

		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Printf("What is your Home Assistant URL? [%s]: ", cfg.HomeAssistant.URL)
			input, err := reader.ReadString('\n')
			if err != nil {
				return nil, fmt.Errorf("read Home Assistant URL: %w", err)
			}
			if value := strings.TrimSpace(input); value != "" {
				cfg.HomeAssistant.URL = value
			}
			if strings.HasPrefix(cfg.HomeAssistant.URL, "http://") || strings.HasPrefix(cfg.HomeAssistant.URL, "https://") {
				break
			}
			fmt.Fprintln(os.Stderr, "Please enter a valid URL")
		}

		fmt.Println("You can create a Long-Lived Access Token in your Home Assistant profile.")
		for {
			token, err := readHomeAssistantToken()
			if err != nil {
				return nil, fmt.Errorf("read Home Assistant token: %w", err)
			}
			cfg.HomeAssistant.Token = strings.TrimSpace(token)
			if cfg.HomeAssistant.Token != "" {
				break
			}
			fmt.Fprintln(os.Stderr, "Please enter a valid token")
		}

		if err := cfg.Save(); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func readHomeAssistantToken() (token string, err error) {
	fd := int(os.Stdin.Fd())
	state, err := term.MakeRaw(fd)
	if err != nil {
		return "", err
	}
	defer func() {
		err = errors.Join(err, term.Restore(fd, state))
	}()

	// In raw mode Ctrl+C is handled as input, so restoration runs on cancellation.
	terminal := term.NewTerminal(struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}, "")
	return terminal.ReadPassword("What is your Home Assistant Token (Long-Lived Access Token)? ")
}
