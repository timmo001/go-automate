---
title: CLI
description: Run tasks and control Home Assistant from the Go Automate command line.
---

The Go Automate CLI groups its commands by area. Most live under the `home-assistant`
command (aliased `ha`), with `notify` available at the top level. Every command
group has short aliases so you can bind them to keyboard shortcuts comfortably.

## Command groups

| Command | Alias | Purpose |
| --- | --- | --- |
| `home-assistant` | `ha` | Control Home Assistant and watch entities. |
| `notify` | `n` | Send a desktop notification. |

For the full command tree with every flag, see the [Commands reference](/reference/commands/).

## Home Assistant

The `ha` command controls entities and watches their state:

```bash
# Toggle a light entity (light.bedroom_lamp)
go-automate ha light toggle bedroom_lamp

# Turn a switch on
go-automate ha switch turn-on desk_fan

# Announce to an area through an assist satellite
go-automate ha assist_satellite announce living_room "Dinner is ready"
```

See [Home Assistant](/using/home-assistant/) for the full set of control commands, and
[Watching Entities](/using/watching/) for live state.

## Notify

Send a desktop notification through `notify-send`:

```bash
go-automate notify "Build complete" "Your build finished successfully"
```

See [Notifications](/using/notifications/) for the details.

## Help

Running Go Automate with no command shows CLI help after any required
[first-run setup](/configuration/#first-run-setup). Use `--help` to see a command's usage:

```bash
go-automate --help
go-automate ha --help
```

## Tips

- Use the short aliases for shortcuts, for example `go-automate ha l toggle bedroom_lamp`.
- Press <kbd>Tab</kbd> to complete commands, aliases, and flags. See [Shell completions](/using/completions/).
- Entity commands take the entity name **without** its domain prefix. `go-automate ha light turn-on bedroom_lamp` acts on `light.bedroom_lamp`.
- Wrap multi-word arguments in quotes, for example announce messages.

## Next steps

- Control [Home Assistant](/using/home-assistant/).
- [Watch entities](/using/watching/) and feed [status bars](/using/watching/#status-bars).
- Browse everything in the [Commands reference](/reference/commands/).
