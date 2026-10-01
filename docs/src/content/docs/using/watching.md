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

## How watching works

`ha bridge watch entity` connects to the shared [bridge](/running/), so many watchers reuse
one connection to Home Assistant. Start the bridge with
[`go-automate ha bridge serve`](/running/), or run it as a service.

## Watch through the bridge

```bash
go-automate ha bridge watch entity input_boolean.guest_mode
```

The watcher prints the current state immediately, then prints again on every change. Pass
`--socket` to use a non-default bridge socket path:

```bash
go-automate ha bridge watch entity input_boolean.guest_mode --socket /tmp/go-automate-ha.sock
```

## Status bars

Add `--bar-json` to emit machine-readable JSON lines instead of plain text. Each line is an
object with `text`, `tooltip` and `class`, plus an optional `name` (the entity's display
name), that any status bar, shell or script can consume, including
[Waybar](https://github.com/Alexays/Waybar) and [Quickshell](https://quickshell.org/).
See [Bar JSON](/reference/bar-json/) for the full output contract and every flag.

```bash
go-automate ha bridge watch entity input_boolean.guest_mode \
  --bar-json \
  --text-on "Guest" \
  --tooltip-on "Guest mode is on" \
  --tooltip-off "Guest mode is off" \
  --class-on "active" \
  --class-off "inactive"
```

:::tip
When the output is consumed by a status bar or another program, always use `--bar-json`.
Without it, Go Automate prints plain text and warns that machine consumers should switch to
JSON.
:::

For the full output contract, every `--bar-json` flag and a complete Waybar module, see
[Bar JSON](/reference/bar-json/).

## Next steps

- See [Bar JSON](/reference/bar-json/) to wire a watcher into a status bar.
- Make sure the [bridge](/running/) is running for the lowest network usage.
- See the [Bridge Protocol](/reference/bridge/) for how watchers talk to the bridge.
