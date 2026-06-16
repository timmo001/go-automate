---
title: Bridge Protocol
description: The Unix socket JSON protocol the Go Automate Home Assistant bridge speaks.
---

The [bridge](/running/) keeps one WebSocket connection to Home Assistant and serves cached
state and live updates to local clients over a Unix domain socket. Most of the time you use
it through [`ha bridge watch entity`](/using/watching/), but the protocol is simple enough
to speak directly from your own scripts.

## Socket location

By default the bridge listens at:

```text
$XDG_RUNTIME_DIR/go-automate/home-assistant.sock
```

If `XDG_RUNTIME_DIR` is not set, it falls back to the system temp directory. The socket
directory is created with `0700` permissions and the socket with `0600`, so only your user
can connect. Override the path with `--socket` on
[`ha bridge serve`](/reference/commands/#ha-bridge-serve).

## Transport

The protocol is newline-delimited JSON over the socket. The client sends a single request
object, and the bridge replies with one or more response objects on the same connection.

## Request format

```json title="Request"
{
  "action": "watch_entity",
  "entity_id": "light.bedroom_lamp"
}
```

- `action`: the action to perform, either `get_entity` or `watch_entity`.
- `entity_id`: the full entity ID, including its domain. Required for both actions.

## Response format

```json title="Response"
{
  "type": "snapshot",
  "entity_id": "light.bedroom_lamp",
  "state": {
    "entity_id": "light.bedroom_lamp",
    "state": "on"
  },
  "name": "Bedroom Lamp"
}
```

- `type`: the response type, one of `snapshot`, `state_changed`, or `error`.
- `entity_id`: the entity the response relates to.
- `state`: an object with `entity_id`, the `state` string, and any `attributes`, or `null` if
  the bridge has no cached state for that entity.
- `name`: the entity's structured display name (device and entity name combined), resolved
  from the registries. Omitted when no name can be resolved. See
  [Bar JSON](/reference/bar-json/) for the naming model.
- `error`: a human-readable message, present only when `type` is `error`.

## Actions

### `get_entity`

Returns a single `snapshot` response with the entity's current cached state, then closes the
connection. `state` is `null` if the entity is unknown to the bridge.

```bash
SOCKET="$XDG_RUNTIME_DIR/go-automate/home-assistant.sock"
echo '{"action":"get_entity","entity_id":"light.bedroom_lamp"}' | socat - "UNIX-CONNECT:$SOCKET"
```

### `watch_entity`

Sends a `snapshot` immediately if the entity is known, then a `state_changed` response every
time the entity changes. The connection stays open until the client disconnects or the
bridge shuts down.

```bash
SOCKET="$XDG_RUNTIME_DIR/go-automate/home-assistant.sock"
echo '{"action":"watch_entity","entity_id":"light.bedroom_lamp"}' | socat - "UNIX-CONNECT:$SOCKET"
```

A change then arrives as:

```json title="state_changed"
{
  "type": "state_changed",
  "entity_id": "light.bedroom_lamp",
  "state": { "entity_id": "light.bedroom_lamp", "state": "off" },
  "name": "Bedroom Lamp"
}
```

## Errors

A bad request returns an `error` response:

```json title="Error response"
{
  "type": "error",
  "error": "entity_id is required"
}
```

Errors are returned when:

- `entity_id` is missing.
- `action` is not `get_entity` or `watch_entity`.
- The request cannot be decoded as JSON.

## Upstream behaviour

The bridge dials Home Assistant, authenticates with your token, reads all states once to
build its cache, then subscribes to `state_changed` events. On each connection it also fetches
the entity and device registries once and caches them to resolve entity names, so adding the
`name` field costs no extra request per watcher. If the connection drops it reconnects every
five seconds, refreshes the registries, and re-broadcasts a fresh snapshot to active watchers,
so clients recover automatically.
