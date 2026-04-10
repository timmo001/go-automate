package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/charmbracelet/log"

	"github.com/timmo001/go-automate/config"
	"github.com/timmo001/go-automate/homeassistant"
	"github.com/timmo001/go-automate/notify"
	"github.com/urfave/cli/v3"
	"golang.org/x/term"
)

// Version is the version of the application, set at build time via ldflags
var Version = "dev"

func main() {
	log.Info("------ Go Automate ------")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Debugf("Loaded config: %v", cfg)

	cfg, err = cfg.Setup(isInteractiveSession())
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	homeassistant.Config = &cfg.HomeAssistant

	cmd := &cli.Command{
		Name:    "Go Automate",
		Usage:   "Run common tasks",
		Version: Version,
		Commands: []*cli.Command{
			{
				Name:    "home-assistant",
				Aliases: []string{"ha"},
				Usage:   "Interact with Home Assistant",
				Commands: []*cli.Command{
					{
						Name:    "watch",
						Aliases: []string{"w"},
						Usage:   "Watch Home Assistant entities for state changes",
						Commands: []*cli.Command{
							{
								Name:      "entity",
								Aliases:   []string{"e"},
								ArgsUsage: "<entity_id>",
								Flags:     createEntityWatchFlags(true),
								Action: func(ctx context.Context, cmd *cli.Command) error {
									return cmdHAWatchEntity(ctx, cmd)
								},
							},
						},
					},
					{
						Name:    "bridge",
						Aliases: []string{"b"},
						Usage:   "Run and query the local Home Assistant bridge",
						Commands: []*cli.Command{
							{
								Name:  "serve",
								Usage: "Serve a shared Home Assistant websocket bridge",
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:  "socket",
										Usage: "Path to the Home Assistant bridge socket",
									},
								},
								Action: func(ctx context.Context, cmd *cli.Command) error {
									return cmdHABridgeServe(ctx, cmd)
								},
							},
							{
								Name:    "watch",
								Aliases: []string{"w"},
								Usage:   "Watch entities through the local Home Assistant bridge",
								Commands: []*cli.Command{
									{
										Name:      "entity",
										Aliases:   []string{"e"},
										ArgsUsage: "<entity_id>",
										Flags:     createEntityWatchFlags(false, &cli.StringFlag{Name: "socket", Usage: "Path to the Home Assistant bridge socket"}),
										Action: func(ctx context.Context, cmd *cli.Command) error {
											return cmdHABridgeWatchEntity(ctx, cmd)
										},
									},
								},
							},
						},
					},
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
			{
				Name:    "notify",
				Aliases: []string{"n"},
				Usage:   "Send a notification",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return cmdNotify(cmd)
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatalf("error running cmd: %v", err)
	}

	log.Info("------ Exiting ------")
}

func isInteractiveSession() bool {
	return term.IsTerminal(int(os.Stdin.Fd())) && term.IsTerminal(int(os.Stdout.Fd()))
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

func createEntityWatchFlags(includeDirectFlags bool, extraFlags ...cli.Flag) []cli.Flag {
	flags := []cli.Flag{
		&cli.BoolFlag{
			Name:  "waybar",
			Usage: "Output JSON lines for Waybar",
		},
		&cli.StringFlag{
			Name:  "icon",
			Usage: "Icon to render for on/off states in Waybar mode",
		},
		&cli.StringFlag{
			Name:  "text-on",
			Usage: "Text to render when the state is on in Waybar mode",
		},
		&cli.StringFlag{
			Name:  "text-off",
			Usage: "Text to render when the state is not on in Waybar mode",
		},
		&cli.StringFlag{
			Name:  "tooltip-on",
			Usage: "Tooltip when the state is on in Waybar mode",
		},
		&cli.StringFlag{
			Name:  "tooltip-off",
			Usage: "Tooltip when the state is not on in Waybar mode",
		},
		&cli.StringFlag{
			Name:  "class-on",
			Usage: "CSS class when the state is on in Waybar mode",
		},
		&cli.StringFlag{
			Name:  "class-off",
			Usage: "CSS class when the state is not on in Waybar mode",
		},
		&cli.BoolFlag{
			Name:  "hide-off",
			Usage: "Hide the Waybar module when the state is not on",
		},
	}

	if includeDirectFlags {
		flags = append(flags,
			&cli.BoolFlag{
				Name:  "direct",
				Usage: "Bypass the local Home Assistant bridge and connect to Home Assistant directly",
			},
			&cli.StringFlag{
				Name:  "bridge-socket",
				Usage: "Path to the Home Assistant bridge socket",
			},
		)
	}

	flags = append(flags, extraFlags...)

	return flags
}

func cmdHACallService(
	cmd *cli.Command,
	domain, service, targetType string,
	data any,
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
	resp, err := conn.SendRequest(homeassistant.HomeAssistantCallServiceRequest{
		ID:             homeassistant.RandomID(),
		Type:           "call_service",
		Domain:         domain,
		Service:        service,
		ServiceData:    data,
		Target:         map[string]string{targetType: target},
		ReturnResponse: returnResponse,
	}, true)
	if err != nil {
		return err
	}
	log.Infof("Call service response: %v", resp)

	return nil
}

func cmdNotify(
	cmd *cli.Command,
) error {
	args := cmd.Args()
	summary := args.Get(0)
	body := args.Get(1)

	return notify.SendNotification(&notify.Notify{
		Summary: summary,
		Body:    &body,
	})
}

func cmdHAWatchEntity(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args()
	entityID := args.Get(0)
	if entityID == "" {
		return fmt.Errorf("entity_id is required")
	}

	options := entityWatchOutputOptions{
		Waybar:     cmd.Bool("waybar"),
		Icon:       cmd.String("icon"),
		TextOn:     cmd.String("text-on"),
		TextOff:    cmd.String("text-off"),
		TooltipOn:  cmd.String("tooltip-on"),
		TooltipOff: cmd.String("tooltip-off"),
		ClassOn:    cmd.String("class-on"),
		ClassOff:   cmd.String("class-off"),
		HideOff:    cmd.Bool("hide-off"),
	}

	socketPath, err := resolveBridgeSocketPath(cmd.String("bridge-socket"))
	if err != nil {
		return err
	}

	if !cmd.Bool("direct") {
		if err := watchEntityViaBridge(ctx, socketPath, entityID, options); err == nil {
			return nil
		} else {
			log.Debugf("Could not use Home Assistant bridge at %s, falling back to direct websocket: %v", socketPath, err)
		}
	}

	return watchEntityDirect(entityID, options)
}

func cmdHABridgeServe(ctx context.Context, cmd *cli.Command) error {
	socketPath, err := resolveBridgeSocketPath(cmd.String("socket"))
	if err != nil {
		return err
	}

	serveCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	bridge, err := homeassistant.NewBridge(homeassistant.Config, socketPath)
	if err != nil {
		return err
	}

	return bridge.Serve(serveCtx)
}

func cmdHABridgeWatchEntity(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args()
	entityID := args.Get(0)
	if entityID == "" {
		return fmt.Errorf("entity_id is required")
	}

	options := entityWatchOutputOptions{
		Waybar:     cmd.Bool("waybar"),
		Icon:       cmd.String("icon"),
		TextOn:     cmd.String("text-on"),
		TextOff:    cmd.String("text-off"),
		TooltipOn:  cmd.String("tooltip-on"),
		TooltipOff: cmd.String("tooltip-off"),
		ClassOn:    cmd.String("class-on"),
		ClassOff:   cmd.String("class-off"),
		HideOff:    cmd.Bool("hide-off"),
	}

	socketPath, err := resolveBridgeSocketPath(cmd.String("socket"))
	if err != nil {
		return err
	}

	return watchEntityViaBridge(ctx, socketPath, entityID, options)
}

func resolveBridgeSocketPath(socketPath string) (string, error) {
	if socketPath != "" {
		return socketPath, nil
	}

	return homeassistant.DefaultBridgeSocketPath()
}

func watchEntityViaBridge(
	ctx context.Context,
	socketPath string,
	entityID string,
	options entityWatchOutputOptions,
) error {
	return homeassistant.BridgeWatchEntity(ctx, socketPath, entityID, func(state *homeassistant.HomeAssistantState) error {
		printEntityState(state, options)
		return nil
	})
}

func watchEntityDirect(entityID string, options entityWatchOutputOptions) error {
	conn := homeassistant.Connect()
	defer conn.Close()

	initialState, err := conn.GetState(entityID)
	if err != nil {
		return err
	}
	if initialState != nil {
		printEntityState(initialState, options)
	}

	resp, err := conn.SubscribeEvents("state_changed")
	if err != nil {
		return err
	}
	if !resp.Success {
		return fmt.Errorf("subscribe failed: %s", resp.Error.Message)
	}

	for {
		event, err := conn.ReadEvent()
		if err != nil {
			return err
		}
		if event.Type != "event" || event.Event.EventType != "state_changed" {
			continue
		}

		if event.Event.Data.EntityID != entityID || event.Event.Data.NewState == nil {
			continue
		}

		printEntityState(event.Event.Data.NewState, options)
	}
}

type entityWatchOutputOptions struct {
	Waybar     bool
	Icon       string
	TextOn     string
	TextOff    string
	TooltipOn  string
	TooltipOff string
	ClassOn    string
	ClassOff   string
	HideOff    bool
}

func appendWaybarText(baseText string, label string) string {
	if label == "" {
		return baseText
	}
	if baseText == "" {
		return label
	}

	return fmt.Sprintf("%s %s", baseText, label)
}

func printEntityState(state *homeassistant.HomeAssistantState, options entityWatchOutputOptions) {
	if options.Waybar {
		text := state.State
		tooltip := state.State
		className := state.State

		if state.State == "on" {
			if options.Icon != "" {
				text = options.Icon
			}
			text = appendWaybarText(text, options.TextOn)
			if options.TooltipOn != "" {
				tooltip = options.TooltipOn
			}
			if options.ClassOn != "" {
				className = options.ClassOn
			}
		} else {
			if options.HideOff {
				text = ""
			} else if options.Icon != "" {
				text = options.Icon
			} else if options.Icon == "" {
				text = state.State
			}
			text = appendWaybarText(text, options.TextOff)
			if options.TooltipOff != "" {
				tooltip = options.TooltipOff
			}
			if options.ClassOff != "" {
				className = options.ClassOff
			}
			if options.HideOff {
				if className == "" {
					className = "hidden"
				} else {
					className += " hidden"
				}
			}
		}

		payload, err := json.Marshal(map[string]string{
			"text":    text,
			"tooltip": tooltip,
			"class":   className,
		})
		if err != nil {
			log.Fatalf("error marshalling waybar payload: %v", err)
		}

		fmt.Println(string(payload))
		return
	}

	fmt.Println(state.State)
}
