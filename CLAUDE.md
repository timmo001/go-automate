# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go Automate is a CLI utility designed to trigger automations in Home Assistant via keyboard shortcuts and patched Linux apps. It provides a command-line interface to control Home Assistant entities and send desktop notifications.

## Build and Development Commands

```bash
# Build the project
go build

# Install to $GOPATH/bin
go install

# Run without installing
go run main.go [command]

# Build and run
go build && ./go-automate [command]

# Install dependencies
go mod download

# Update dependencies
go mod tidy

# Format code
go fmt ./...
```

## Architecture

### Package Structure

- `main.go` - Entry point; CLI command definitions using urfave/cli/v3
- `config/` - Configuration management using viper; loads from `~/.config/go-automate/config.yml`
- `homeassistant/` - Home Assistant WebSocket API client
- `notify/` - Desktop notification wrapper around `notify-send`

### Key Patterns

**CLI Structure**: Commands are hierarchical using urfave/cli/v3. The main command tree is defined in main.go with nested commands for different domains (e.g., `ha light turn-on`, `ha switch toggle`).

**Configuration**:
- Config file location follows XDG standards: `~/.config/go-automate/config.yml` (Linux), `%APPDATA%\go-automate` (Windows)
- First-run setup prompts for Home Assistant URL and Long-Lived Access Token using charmbracelet/huh forms
- Configuration is loaded at startup and stored in a global variable (`homeassistant.Config`)

**Home Assistant Integration**:
- Uses WebSocket API (not REST API) for real-time communication
- Connection flow: Connect → Read welcome message → Authenticate → Send requests
- Each service call generates a random ID for request/response matching
- Commands follow Home Assistant's domain/service pattern (e.g., `light.turn_on`, `switch.toggle`)

**Entity Commands**:
- The `createToggleServiceCommands()` helper generates turn-on/turn-off/toggle commands for domains like `light`, `switch`, and `input_boolean`
- Entity IDs are constructed as `{domain}.{first_arg}` (e.g., `light.kitchen` from `ha light turn-on kitchen`)

### Command Examples

```bash
# Control lights
go-automate ha light turn-on kitchen
go-automate ha light toggle bedroom

# Control switches
go-automate ha switch turn-off fan

# Control input booleans
go-automate ha input_boolean toggle automation_enabled

# Announce via assist satellite
go-automate ha assist_satellite announce "Hello World"

# Send desktop notification
go-automate notify "Summary" "Body text"
```

## Testing

No test files currently exist in the codebase.

## Dependencies

Key dependencies:
- `github.com/urfave/cli/v3` - CLI framework
- `github.com/spf13/viper` - Configuration management
- `github.com/gorilla/websocket` - WebSocket client for Home Assistant
- `github.com/charmbracelet/huh` - Interactive forms for setup
- `github.com/charmbracelet/log` - Structured logging

## Notes

- The application requires a running Home Assistant instance with WebSocket API enabled
- Desktop notifications require `notify-send` to be installed on the system
- Token authentication uses Home Assistant Long-Lived Access Tokens
- WebSocket connections are created per-command and not persisted
