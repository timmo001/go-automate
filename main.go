package main

import (
	"github.com/charmbracelet/log"

	"github.com/timmo001/go-automate/config"
	"github.com/timmo001/go-automate/homeassistant"
)

func main() {
	log.Info("------ Go Automate ------")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Infof("Loaded config: %v", cfg)

	cfg, err = cfg.Setup()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	homeassistant.Config = &cfg.HomeAssistant

  RunCommand()

	log.Info("------ Exiting ------")
}

