package main

import (
	"context"
	"os"

	"github.com/charmbracelet/log"

	"github.com/timmo001/go-automate/config"
	"github.com/urfave/cli/v3"
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

									return ha.
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
