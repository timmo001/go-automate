package notify

import (
	"fmt"
	"os/exec"
)

// Based on options of notify-send command
type Notify struct {
	Summary string
	// Optional parameters
	Body *string
	// Optional options
	Urgency    *string
	ExpireTime *string
	AppName    *string
	Icon       *string
	Category   *string
	Transient  *bool
	Hint       *string
	Actions    *[]string
}

func SendNotification(
	data *Notify,
) error {
	// Run notify-send command with data
	cmd := exec.Command("notify-send", data.Summary)
	if data.Body != nil {
		cmd.Args = append(cmd.Args, *data.Body)
	}

	if data.Urgency != nil {
		cmd.Args = append(cmd.Args, fmt.Sprintf("--urgency=%s", *data.Urgency))
	}
	if data.ExpireTime != nil {
		cmd.Args = append(cmd.Args, fmt.Sprintf("--expire-time=%s", *data.ExpireTime))
	}
	if data.AppName != nil {
		cmd.Args = append(cmd.Args, fmt.Sprintf("--app-name=%s", *data.AppName))
	}
	if data.Icon != nil {
		cmd.Args = append(cmd.Args, fmt.Sprintf("--icon=%s", *data.Icon))
	}
	if data.Category != nil {
		cmd.Args = append(cmd.Args, fmt.Sprintf("--category=%s", *data.Category))
	}
	if data.Transient != nil && *data.Transient {
		cmd.Args = append(cmd.Args, "--transient")
	}
	if data.Hint != nil {
		cmd.Args = append(cmd.Args, fmt.Sprintf("--hint=%s", *data.Hint))
	}
	if data.Actions != nil {
		for _, action := range *data.Actions {
			cmd.Args = append(cmd.Args, fmt.Sprintf("--action=%s", action))
		}
	}

	return cmd.Run()
}
