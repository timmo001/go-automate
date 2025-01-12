package main

import (
	"sync"

	"github.com/charmbracelet/log"

	"github.com/timmo001/go-automate/config"
)

func main() {
	log.Info("------ Go Automate ------")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	cfg, err = cfg.Setup()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	mainthread.Init(func() {
		err := setupHotkeys(cfg)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	})

	log.Info("------ Exiting ------")
}
