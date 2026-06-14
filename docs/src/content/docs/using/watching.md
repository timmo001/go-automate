---
title: Watching Entities
description: Stream Home Assistant entity state changes into scripts and status bars.
---

Go Automate can watch a Home Assistant entity and print its state every time it changes.
This powers shell scripts and status-bar modules that react to your home in real time.

:::note
Watch commands take the **full** entity ID, including its domain, for example
`input_boolean.guest_mode` or `sensor.living_room_temperature`. This is different from the
[control commands](/using/home-assistant/), which take the name without the domain.
:::

## Two ways to watch

- **Through the bridge (recommended)** — `ha bridge watch entity` connects to the shared
  [bridge](/running/), so many watchers reuse one connection to Home Assistant. This is the
  best choice for status bars and long-running watchers.
- **Direct (troubleshooting only)** — `ha watch entity` opens its own WebSocket connection.
  It uses the bridge by default and only falls back to a direct connection if the bridge is
  unavailable. Force a direct connection with `--direct` when debugging.

:::caution
Direct connections open a new WebSocket per watcher and increase network usage. Prefer the
bridge for anything that runs continuously. Start it with
[`go-automate ha bridge serve`](/running/).
:::

## Watch through the bridge

```bash
go-automate ha bridge watch entity input_boolean.guest_mode
```

The watcher prints the current state immediately, then prints again on every change. Pass
`--socket` to use a non-default bridge socket path:

```bash
go-automate ha bridge watch entity input_boolean.guest_mode --socket /tmp/go-automate-ha.sock
```

## Watch directly

`ha watch entity` uses the bridge when it can and falls back to a direct connection
otherwise:

```bash
go-automate ha watch entity sensor.living_room_temperature
```

Force a direct connection for troubleshooting with `--direct`:

```bash
go-automate ha watch entity sensor.living_room_temperature --direct
```

## Status bars

Add `--bar-json` to emit machine-readable JSON lines instead of plain text. Each line is an
object with `text`, `tooltip` and `class` that any status bar, shell or script can consume,
including [Waybar](https://github.com/Alexays/Waybar) and [Quickshell](https://quickshell.org/).
See [Bar JSON](/reference/bar-json/) for the full output contract and every flag.

```bash
go-automate ha bridge watch entity input_boolean.guest_mode \
  --bar-json \
  --text-on "Guest" \
  --tooltip-on "Guest mode is on" \
  --tooltip-off "Guest mode is off" \
  --class-on "active" \
  --hide-off
```

:::tip
When the output is consumed by a status bar or another program, always use `--bar-json`.
Without it, Go Automate prints plain text and warns that machine consumers should switch to
JSON.
:::

For the full output contract, every `--bar-json` flag and a complete Waybar module, see
[Bar JSON](/reference/bar-json/).

Connection flags differ by command:

| Flag | Command | Effect |
| --- | --- | --- |
| `--socket` | `ha bridge watch entity` | Path to the bridge socket. |
| `--direct` | `ha watch entity` | Bypass the bridge and connect directly. |
| `--bridge-socket` | `ha watch entity` | Path to the bridge socket to try before falling back. |

## Next steps

- See [Bar JSON](/reference/bar-json/) to wire a watcher into a status bar.
- Make sure the [bridge](/running/) is running for the lowest network usage.
- See the [Bridge Protocol](/reference/bridge/) for how watchers talk to the bridge.
