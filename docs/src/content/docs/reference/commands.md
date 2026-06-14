---
title: Commands
description: Every Go Automate command, alias and flag in one place.
---

This page lists the full Go Automate command tree. Run any command with `--help` to see its
usage at the terminal.

## `go-automate`

The root command. With no subcommand in an interactive terminal, it launches the
[TUI](/using/tui/) when the `go-automate-tui` binary is available; otherwise it shows CLI
help.

| Flag | Effect |
| --- | --- |
| `--version` | Print the version and exit. |
| `--help` | Show help for the command. |

## `home-assistant` (alias `ha`)

Interact with Home Assistant.

### `ha watch entity <entity_id>` (alias `ha w e`)

Watch an entity for state changes. Uses the [bridge](/running/) by default and falls back
to a direct WebSocket connection if the bridge is unavailable. Prefer
`ha bridge watch entity` (below) for anything long-running.

| Flag | Effect |
| --- | --- |
| `--bar-json` | Emit JSON lines (`text`, `tooltip`, `class`) for status bars. |
| `--icon` | Text/icon to show for the state in bar JSON mode. |
| `--text-on` | Text appended when the state is `on`. |
| `--text-off` | Text appended when the state is not `on`. |
| `--tooltip-on` | Tooltip when the state is `on`. |
| `--tooltip-off` | Tooltip when the state is not `on`. |
| `--class-on` | Status-bar class when the state is `on`. |
| `--class-off` | Status-bar class when the state is not `on`. |
| `--hide-off` | Hide the module when the state is not `on`. |
| `--direct` | Bypass the bridge and connect directly (higher network usage). |
| `--bridge-socket` | Path to the bridge socket to try before falling back. |

See [Bar JSON](/reference/bar-json/) for the output shape and how the `--bar-json` flags
combine.

### `ha bridge serve`

Serve the shared Home Assistant bridge. See [Running the Bridge](/running/).

| Flag | Effect |
| --- | --- |
| `--socket` | Path to the bridge socket. Defaults to `$XDG_RUNTIME_DIR/go-automate/home-assistant.sock`. |

### `ha bridge watch entity <entity_id>` (alias `ha b w e`)

Watch an entity through the bridge (recommended). Takes the same `--bar-json` output flags
as `ha watch entity`.

| Flag | Effect |
| --- | --- |
| `--bar-json` | Emit JSON lines for status bars. |
| `--icon` | Text/icon to show for the state in bar JSON mode. |
| `--text-on` / `--text-off` | Text appended for the on / not-on states. |
| `--tooltip-on` / `--tooltip-off` | Tooltip for the on / not-on states. |
| `--class-on` / `--class-off` | Status-bar class for the on / not-on states. |
| `--hide-off` | Hide the module when the state is not `on`. |
| `--socket` | Path to the bridge socket. |

### `ha assist_satellite announce <area_id> <message>` (alias `ha as a`)

Announce a message to an area through an assist satellite. The first argument is the area
ID; the second is the message.

```bash
go-automate ha assist_satellite announce living_room "Dinner is ready"
```

### `ha light` (alias `ha l`)

Control light entities. Each subcommand takes the entity name without its domain.

| Subcommand | Alias | Service |
| --- | --- | --- |
| `turn-on <name>` | `on` | `light.turn_on` |
| `turn-off <name>` | `off` | `light.turn_off` |
| `toggle <name>` | `t` | `light.toggle` |

### `ha switch` (alias `ha s`)

Control switch entities. Same subcommands as `light`.

| Subcommand | Alias | Service |
| --- | --- | --- |
| `turn-on <name>` | `on` | `switch.turn_on` |
| `turn-off <name>` | `off` | `switch.turn_off` |
| `toggle <name>` | `t` | `switch.toggle` |

### `ha input_boolean` (alias `ha ib`)

Control input boolean helpers. Same subcommands as `light`.

| Subcommand | Alias | Service |
| --- | --- | --- |
| `turn-on <name>` | `on` | `input_boolean.turn_on` |
| `turn-off <name>` | `off` | `input_boolean.turn_off` |
| `toggle <name>` | `t` | `input_boolean.toggle` |

## `notify <summary> [body]` (alias `n`)

Send a desktop notification through `notify-send`. The first argument is the summary; the
optional second argument is the body. See [Notifications](/using/notifications/).

```bash
go-automate notify "Build complete" "Your build finished successfully"
```

## `tui`

Launch the interactive [TUI](/using/tui/) menu.
