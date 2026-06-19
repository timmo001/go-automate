# Agent Guide for Go Automate

## Documentation
- `docs/` is the source of truth for all user-facing documentation. It is an Astro/Starlight site; content lives in `docs/src/content/docs/` and is served at <https://go-automate.timmo.dev>.
- Any change to behaviour, commands, flags, configuration, the bridge protocol, or the TUI must update the relevant page under `docs/src/content/docs/` in the same change.
- Keep the root `README.md` short; it should link to `docs/` rather than duplicate content. Do not let other Markdown files become a second source of truth.
- Reference pages map to features: `using/cli.md`, `using/tui.md`, `using/home-assistant.md`, `using/watching.md`, `using/notifications.md`, `reference/commands.md`, `reference/bridge.md`, plus `install.mdx`, `configuration.mdx`, and `running.mdx`.
- Verify doc edits build with `pnpm build` from the `docs/` directory.

## Build & Test Commands
- Toolchain: `mise install` provisions Go, Bun, Node and pnpm from `mise.toml`; list tasks with `mise tasks`
- Build: `mise run build` (plain `go build` also works)
- Build app + TUI: `mise run build:all`
- Package (Arch): `mise run package:arch`
- Install (Arch): `yay -U dist/go-automate-<version>-1-x86_64.pkg.tar.zst`
- Run: `mise run run` (or `go run main.go [command]`)
- Test: `mise run test`
- Lint/format: `mise run lint` (`go fmt ./...` + `go vet ./...`)
- Docs: `mise run docs:dev`, `mise run docs:build`
- Update deps: `mise run deps` (or `go mod tidy`)
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

## Home Assistant Bridge Watch Policy (Go Automate)

- For entity watchers, prefer bridge-backed commands by default: `go-automate ha bridge watch entity ...`
- Treat direct websocket watcher usage as exceptional: allow only for explicit troubleshooting (`--direct`) and surface a warning in CLI output/help text
- In help and flag descriptions, strongly recommend bridge watch for lower network usage
- When output is plain text (no `--bar-json`), warn users that machine consumers should prefer `--bar-json` JSON output
