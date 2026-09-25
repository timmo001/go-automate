package config

import (
	"os"
	"runtime"
	"testing"

	"github.com/spf13/viper"
)

func TestFirstRunConfigCanBeSavedAndLoaded(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	viper.Reset()
	t.Cleanup(viper.Reset)

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	cfg.HomeAssistant.Token = "example-token"
	if err := cfg.Save(); err != nil {
		t.Fatal(err)
	}
	path, err := configFilePath()
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS != "windows" && info.Mode().Perm() != 0600 {
		t.Fatalf("config permissions = %04o, want 0600", info.Mode().Perm())
	}

	viper.Reset()
	loaded, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if loaded.HomeAssistant != cfg.HomeAssistant {
		t.Fatal("saved configuration was not loaded")
	}
}
