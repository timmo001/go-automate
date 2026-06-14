---
title: Home Assistant
description: Control lights, switches, input booleans and assist satellites from the command line.
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
| `ha input_boolean turn-on` / `turn-off` / `toggle` | `ha ib on` / `off` / `t` | `input_boolean.turn_on` / `turn_off` / `toggle` | `input_boolean.<name>` |
| `ha assist_satellite announce` | `ha as a` | `assist_satellite.announce` | `area_id` |

## Next steps

- [Watch entities](/using/watching/) to react to state changes.
- See every flag in the [Commands reference](/reference/commands/).
