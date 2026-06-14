---
title: Bar JSON
description: The machine-readable JSON output that entity watchers emit for status bars and shell modules.
---

`--bar-json` switches an entity watcher from plain text to machine-readable JSON. Each state
change prints one JSON object on its own line, ready to feed straight into a status bar such
as [Waybar](https://github.com/Alexays/Waybar) or any shell module that reads JSON.

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

Every line is a JSON object with exactly three string fields, always present:

```json
{ "text": "Guest", "tooltip": "Guest mode is on", "class": "active" }
```

| Field | Purpose |
| --- | --- |
| `text` | The label the bar renders. |
| `tooltip` | The hover tooltip. |
| `class` | The CSS class the bar applies for styling. |

These map directly onto a Waybar custom module with `"return-type": "json"`.

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

### Waybar module

```json title="~/.config/waybar/config.jsonc"
"custom/guest_mode": {
  "exec": "go-automate ha bridge watch entity input_boolean.guest_mode --bar-json --text-on 'Guest' --tooltip-on 'Guest mode on' --tooltip-off 'Guest mode off' --class-on 'active' --hide-off",
  "return-type": "json",
  "restart-interval": 5
}
```

The watcher streams a JSON line on every state change, so the module updates live. The
`restart-interval` makes Waybar restart the watcher if it ever exits (for example while the
bridge restarts).

## Next steps

- See [Watching Entities](/using/watching/) for choosing between the bridge and direct
  watchers.
- Run the [bridge](/running/) so many bar modules share one connection to Home Assistant.
