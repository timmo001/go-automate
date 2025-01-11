package main

import (
	"sync"

	"github.com/charmbracelet/log"
	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"

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

func setupHotkeys(cfg *config.Config) error {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()

		err := listenHotkey(hotkey.KeyS, hotkey.ModCtrl, hotkey.ModShift)
		if err != nil {
			log.Errorf("error: %v", err)
		}
	}()
	go func() {
		defer wg.Done()

		err := listenHotkey(hotkey.KeyA, hotkey.ModCtrl, hotkey.ModShift)
		if err != nil {
			log.Errorf("error: %v", err)
		}
	}()
	wg.Wait()

	return nil
}

func listenHotkey(key hotkey.Key, mods ...hotkey.Modifier) (err error) {
	ms := []hotkey.Modifier{}
	ms = append(ms, mods...)
	hk := hotkey.New(ms, key)

	err = hk.Register()
	if err != nil {
		return
	}

	// Blocks until the hokey is triggered.
	<-hk.Keydown()
	log.Printf("hotkey: %v is down\n", hk)
	<-hk.Keyup()
	log.Printf("hotkey: %v is up\n", hk)
	hk.Unregister()
	return
}
