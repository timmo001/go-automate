package main

import (
	"context"
	"crypto/rand"
	"fmt"
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
	homeassistant.Config = &cfg.HomeAssistant

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
						Name:    "assist_satellite",
						Aliases: []string{"as"},
						Commands: []*cli.Command{
							{
								Name:    "announce",
								Aliases: []string{"a"},
								Action: func(ctx context.Context, cmd *cli.Command) error {
									message := cmd.Args().Get(1)
									log.Infof("Announcing: %s", message)

									return cmdHACallService(
										cmd,
										"assist_satellite",
										"announce",
										"area_id",
										map[string]string{
											"message": message,
										},
										false,
									)
								},
							},
						},
					},
					{
						Name:     "input_boolean",
						Aliases:  []string{"ib"},
						Commands: createToggleServiceCommands("input_boolean"),
					},
					{
						Name:     "light",
						Aliases:  []string{"l"},
						Commands: createToggleServiceCommands("light"),
					},
					{
						Name:     "switch",
						Aliases:  []string{"s"},
						Commands: createToggleServiceCommands("switch"),
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

func createToggleServiceCommands(domain string) []*cli.Command {
	return []*cli.Command{
		{
			Name:    "turn-on",
			Aliases: []string{"on"},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				return cmdHACallService(cmd, domain, "turn_on", "entity_id", nil, false)
			},
		},
		{
			Name:    "turn-off",
			Aliases: []string{"off"},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				return cmdHACallService(cmd, domain, "turn_off", "entity_id", nil, false)
			},
		},
		{
			Name:    "toggle",
			Aliases: []string{"t"},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				return cmdHACallService(cmd, domain, "toggle", "entity_id", nil, false)
			},
		},
	}
}

func cmdHACallService(
	cmd *cli.Command,
	domain, service, targetType string,
	data interface{},
	returnResponse bool,
) error {
	args := cmd.Args()
	firstArg := args.Get(0)
	log.Infof("First arg: %s", firstArg)

  var target string
  if targetType == "entity_id" {
    target = fmt.Sprintf("%s.%s", domain, firstArg)
  } else {
    target = firstArg
  }

	conn := homeassistant.Connect()
	resp := conn.SendRequest(homeassistant.HomeAssistantCallServiceRequest{
		ID:             RandomID(),
		Type:           "call_service",
		Domain:         domain,
		Service:        service,
		ServiceData:    data,
		Target:         map[string]string{targetType: target},
		ReturnResponse: returnResponse,
	}, true)
	log.Infof("Call service response: %v", resp)

	return nil
}

func RandomID() int {
	reader := rand.Reader
	n, err := rand.Int(reader, big.NewInt(1000))
	if err != nil {
		log.Fatalf("error generating random ID: %v", err)
	}
	return int(n.Int64())
}
