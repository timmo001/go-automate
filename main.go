package main

import (
	"context"
	"crypto/rand"
	"math/big"
	"os"

	"github.com/charmbracelet/log"

	"github.com/timmo001/go-automate/config"
	"github.com/timmo001/go-automate/homeassistant"
	"github.com/urfave/cli/v3"
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

	cmd := &cli.Command{
		Name:  "Go Automate",
		Usage: "Run common tasks",
		Commands: []*cli.Command{
			{
				Name:    "home-assistant",
				Aliases: []string{"ha"},
				Usage:   "Interact with Home Assistant",
				Commands: []*cli.Command{
					{
						Name:    "input_boolean",
						Aliases: []string{"ib"},
						Commands: []*cli.Command{
							{
								Name:    "toggle",
								Aliases: []string{"t"},
								Action: func(ctx context.Context, cmd *cli.Command) error {
									args := cmd.Args()
									firstArg := args.First()
									log.Infof("First arg: %s", firstArg)

									homeassistant.Config = &cfg.HomeAssistant
									conn := homeassistant.Connect()
									resp := conn.SendRequest(homeassistant.HomeAssistantCallServiceRequest{
										ID:      randomID(),
										Type:    "call_service",
										Domain:  "input_boolean",
										Service: "toggle",
										Target:  map[string]string{"entity_id": "input_boolean." + firstArg},
									})
									log.Infof("Call service response: %v", resp)

									return nil
								},
							},
						},
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatalf("error running cmd: %v", err)
	}

	log.Info("------ Exiting ------")
}

func randomID() int {
	reader := rand.Reader
	n, err := rand.Int(reader, big.NewInt(1000))
	if err != nil {
		log.Fatalf("error generating random ID: %v", err)
	}
	return int(n.Int64())
}
