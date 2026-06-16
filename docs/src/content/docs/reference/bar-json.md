---
title: Bar JSON
description: The machine-readable JSON output that entity watchers emit for status bars and shell modules.
---

`--bar-json` switches an entity watcher from plain text to machine-readable JSON. Each state
change prints one JSON object on its own line. It is a generic JSON-lines contract, so it
works with any status bar, shell or script that reads line-delimited JSON, including
[Waybar](https://github.com/Alexays/Waybar) and [Quickshell](https://quickshell.org/).

It is supported by both watch commands:

- [`ha bridge watch entity`](/reference/commands/) (recommended)
- [`ha watch entity`](/reference/commands/)

See [Watching Entities](/using/watching/) for when to use each one.

:::tip
Whenever a watcher's output is consumed by a status bar or another program, use `--bar-json`.
Without it, Go Automate prints plain text and warns that machine consumers should switch to
JSON.
:::

## Output shape

Every line is a JSON object. Three string fields are always present, and an optional
`name` field is added when Go Automate can resolve the entity's name:

```json
{ "text": "437", "tooltip": "437", "class": "437", "name": "Living Room Thermostat Temperature" }
```

| Field | Purpose |
| --- | --- |
| `text` | The label to render. |
| `tooltip` | The hover tooltip. |
| `class` | A class name for styling. |
| `name` | Optional. The entity's display name, with the device and entity name combined. Omitted when no name can be resolved. |

The `text`, `tooltip` and `class` field names match Waybar's custom-module schema
(`"return-type": "json"`), but the object is plain JSON. `name` is an extra field beyond
Waybar's schema, so bars that do not use it ignore it. Any consumer can read the line and use
whichever fields it supports, so the same command drives a Waybar module, a Quickshell widget
or a custom script.

## How the fields are derived

A watcher treats the entity state as either `on` or not. The check is a literal comparison
against the string `on`, so any other value (`off`, `unavailable`, a sensor reading, and so
on) is handled by the "not on" branch.

All three fields default to the raw entity state. The output flags then refine them:

**When the state is `on`:**

- `text` becomes `--icon` if set, then `--text-on` is appended (space-joined).
- `tooltip` becomes `--tooltip-on` if set.
- `class` becomes `--class-on` if set.

**When the state is not `on`:**

- `text` becomes empty if `--hide-off` is set, otherwise `--icon` if set, otherwise the raw
  state. `--text-off` is then appended.
- `tooltip` becomes `--tooltip-off` if set.
- `class` becomes `--class-off` if set.
- With `--hide-off`, `hidden` is appended to `class` (so the bar can hide the module while
  keeping the line valid JSON).

Any flag you leave unset keeps the field at its default, so a minimal `--bar-json` run still
produces valid output built from the raw state.

**The `name` field:**

`name` is the entity's structured display name, following Home Assistant's entity naming
model: the device name and the entity-specific name combined (for example
`Living Room Thermostat Temperature`). The bridge watcher resolves it from the entity and
device registries, which it fetches and caches when it connects to Home Assistant, so no
extra request is made per watcher. The direct watcher (`ha watch entity --direct`) falls back
to the entity's `friendly_name`. When no name can be resolved, the field is omitted rather
than emitted empty, and the flags above never change it.

## Flags

These flags apply in `--bar-json` mode on both watch commands:

| Flag | Effect |
| --- | --- |
| `--bar-json` | Emit JSON lines (`text`, `tooltip`, `class`) for status bars. |
| `--icon` | Text or icon to show for the state. Replaces the raw state text. |
| `--text-on` | Text appended when the state is `on`. |
| `--text-off` | Text appended when the state is not `on`. |
| `--tooltip-on` | Tooltip when the state is `on`. |
| `--tooltip-off` | Tooltip when the state is not `on`. |
| `--class-on` | Status-bar class when the state is `on`. |
| `--class-off` | Status-bar class when the state is not `on`. |
| `--hide-off` | Hide the module (empty `text`, `hidden` appended to `class`) when the state is not `on`. |

Connection flags differ by command. See [Commands](/reference/commands/) for `--socket`,
`--direct` and `--bridge-socket`.

## Example

```bash
go-automate ha bridge watch entity input_boolean.guest_mode \
  --bar-json \
  --text-on "Guest" \
  --tooltip-on "Guest mode is on" \
  --tooltip-off "Guest mode is off" \
  --class-on "active" \
  --hide-off
```

With the entity on, this prints:

```json
{ "text": "Guest", "tooltip": "Guest mode is on", "class": "active" }
```

With the entity off, `--hide-off` blanks the text and marks the class hidden:

```json
{ "text": "", "tooltip": "Guest mode is off", "class": "off hidden" }
```

## Consuming the output

Any program that reads line-delimited JSON from a process can consume a watcher. The watcher
streams a JSON line on every state change, so the consumer updates live.

### Waybar

```json title="~/.config/waybar/config.jsonc"
"custom/guest_mode": {
  "exec": "go-automate ha bridge watch entity input_boolean.guest_mode --bar-json --text-on 'Guest' --tooltip-on 'Guest mode on' --tooltip-off 'Guest mode off' --class-on 'active' --hide-off",
  "return-type": "json",
  "restart-interval": 5
}
```

The `restart-interval` makes Waybar restart the watcher if it ever exits (for example while
the bridge restarts).

### Other shells and scripts

Quickshell and similar shells run the same command and parse each line from the process's
standard output. The contract is identical: read a line, decode the JSON, use the `text`,
`tooltip` and `class` fields. From a shell, you can read it line by line:

```bash
go-automate ha bridge watch entity input_boolean.guest_mode --bar-json |
  while IFS= read -r line; do
    printf '%s\n' "$line" | jq -r '.text'
  done
```

Restart the watcher if it exits, the same way Waybar's `restart-interval` does, so the
consumer reconnects after a bridge restart.

## Next steps

- See [Watching Entities](/using/watching/) for choosing between the bridge and direct
  watchers.
- Run the [bridge](/running/) so many bar modules share one connection to Home Assistant.
