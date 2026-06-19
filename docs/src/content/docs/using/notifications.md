---
title: Notifications
description: Send desktop notifications from Go Automate with notify-send.
---

The `notify` command (aliased `n`) sends a desktop notification. It shells out to
`notify-send`, so it works with any notification daemon that implements the freedesktop
notification spec (for example [mako](https://github.com/emersion/mako) or
[dunst](https://github.com/dunst-project/dunst)).

## Requirements

- `notify-send`, provided by `libnotify`.
- A running notification daemon.

## Usage

Pass a summary and, optionally, a body:

```bash
go-automate notify "Build complete" "Your build finished successfully"
```

The first argument is the summary (title) and the second is the body. Use the `n` alias for
shortcuts:

```bash
go-automate n "Backup finished"
```

Wrap multi-word arguments in quotes so they are passed as a single value.

## Using it in automations

Because it is a single command, `notify` slots into scripts, build steps and keyboard
shortcuts:

```bash
mise run build && go-automate notify "Build complete" || go-automate notify "Build failed"
```

:::note
The notification is delivered by `notify-send`. If nothing appears, confirm that
`notify-send` is installed and that a notification daemon is running on your session.
:::

## Next steps

- Control [Home Assistant](/using/home-assistant/) from the same shortcuts.
- Browse every command in the [Commands reference](/reference/commands/).
