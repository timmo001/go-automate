# Agent Guide for Go Automate

## Build & Test Commands
- Build: `go build`
- Install: `go install`
- Run: `go run main.go [command]`
- Format: `go fmt ./...`
- Update deps: `go mod tidy`
- No tests currently exist in codebase

## Code Style & Conventions
- **Imports**: Standard library first, blank line, then third-party packages grouped alphabetically
- **Naming**: PascalCase for exported types/functions, camelCase for unexported; use descriptive names (e.g., `ConfigHomeAssistant`, `cmdHACallService`)
- **Error handling**: Return errors up the stack; use `log.Fatalf()` for fatal errors only; wrap errors with `fmt.Errorf()` and `%w` verb
- **Types**: Use struct tags for JSON/mapstructure serialization; pointers for optional fields
- **CLI structure**: Use urfave/cli/v3 with hierarchical commands; provide aliases for common commands
- **Configuration**: Load via viper from `~/.config/go-automate/config.yml`; use XDG standards for cross-platform paths
- **Home Assistant**: WebSocket API only (not REST); connection flow: Connect → Auth → Request; generate random IDs per request
- **Logging**: Use charmbracelet/log package; `log.Info()` for user-facing messages, `log.Debug()` for debug info, `log.Error()`/`log.Fatal()` for errors
