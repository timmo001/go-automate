---
title: Home Assistant
description: Control lights, switches, covers, climate entities, input helpers and assist satellites from the command line.
---

The `home-assistant` command (aliased `ha`) calls Home Assistant services over the
WebSocket API. It connects using the URL and token from your [configuration](/configuration/),
calls the service, and exits.

## Entity naming

Control commands take the entity name **without** its domain prefix. The domain comes
from the command you run, so:

```bash
go-automate ha light turn-on bedroom_lamp
```

acts on the entity `light.bedroom_lamp`.

## Lights

Turn a light on, off, or toggle it:

```bash
go-automate ha light turn-on bedroom_lamp
go-automate ha light turn-off bedroom_lamp
go-automate ha light toggle bedroom_lamp
```

Aliases keep it short for keyboard shortcuts — `ha l on`, `ha l off`, and `ha l t`:

```bash
go-automate ha l t bedroom_lamp
```

## Switches

Switches work the same way under the `switch` command (alias `s`):

```bash
go-automate ha switch turn-on desk_fan
go-automate ha switch toggle desk_fan
```

## Input booleans

Flip helper booleans that drive your automations with the `input_boolean` command
(alias `ib`):

```bash
go-automate ha input_boolean turn-on guest_mode
go-automate ha input_boolean toggle guest_mode
```

## Input numbers

Increment or decrement an input number by its configured step, or set an exact value, with
the `input_number` command (alias `in`). Increment and decrement read the helper's current
state and call `input_number.set_value`, avoiding helper implementations whose built-in
increment and decrement actions fail:

```bash
go-automate ha input_number increment target_temperature
go-automate ha input_number decrement target_temperature
go-automate ha input_number set-value target_temperature 23.8
```

## Covers

Set a cover such as a curtain or blind to a position or tilt percentage, or close
it, with the `cover` command (alias `c`):

```bash
go-automate ha cover position curtain 30
go-automate ha cover tilt-position living_room_left_blind 40
go-automate ha cover close curtain
```

Watch a cover's state and current tilt position through the shared bridge:

```bash
go-automate ha cover watch living_room_left_blind
```

## Climate

Set a climate entity's fan mode with the `climate` command (alias `cl`):

```bash
go-automate ha climate fan-mode air_conditioner 1
```

Watch its HVAC state, labelled fan mode, and target temperature through the shared bridge. Cooling is displayed as `Cool`:

```bash
go-automate ha climate watch air_conditioner
```

## Assist satellites

Announce a message to an area through an assist satellite with
`assist_satellite announce` (alias `as a`). Pass the area ID first, then the message:

```bash
go-automate ha assist_satellite announce living_room "Dinner is ready"
```

:::note
The first argument is the Home Assistant **area ID**, not an entity. Wrap the message in
quotes so it is passed as a single argument.
:::

## Service actions at a glance

| Command | Alias | Service | Target |
| --- | --- | --- | --- |
| `ha light turn-on` / `turn-off` / `toggle` | `ha l on` / `off` / `t` | `light.turn_on` / `turn_off` / `toggle` | `light.<name>` |
| `ha switch turn-on` / `turn-off` / `toggle` | `ha s on` / `off` / `t` | `switch.turn_on` / `turn_off` / `toggle` | `switch.<name>` |
| `ha cover position` / `tilt-position` / `close` | `ha c position` / `tilt-position` / `close` | `cover.set_cover_position` / `set_cover_tilt_position` / `close_cover` | `cover.<name>` |
| `ha climate fan-mode` | `ha cl fan-mode` | `climate.set_fan_mode` | `climate.<name>` |
| `ha input_boolean turn-on` / `turn-off` / `toggle` | `ha ib on` / `off` / `t` | `input_boolean.turn_on` / `turn_off` / `toggle` | `input_boolean.<name>` |
| `ha input_number increment` / `decrement` / `set-value` | `ha in increment` / `decrement` / `set-value` | `input_number.increment` / `decrement` / `set_value` | `input_number.<name>` |
| `ha assist_satellite announce` | `ha as a` | `assist_satellite.announce` | `area_id` |

## Next steps

- [Watch entities](/using/watching/) to react to state changes.
- See every flag in the [Commands reference](/reference/commands/).

## Omarchy plugin

The generated
[`omarchy-go-automate`](https://github.com/timmo001/omarchy-go-automate)
plugin provides a configurable Omarchy Quattro bar widget and panel. Its
user-owned JSON modules file defines entity watchers, actions, labels, and
icons, so personal Home Assistant entity IDs stay outside this repository.

The plugin source and configuration format live in
[`omarchy-plugin/`](https://github.com/timmo001/go-automate/tree/main/omarchy-plugin).
