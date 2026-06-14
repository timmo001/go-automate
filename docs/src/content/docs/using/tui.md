---
title: TUI
description: Browse and run Go Automate commands from an interactive terminal menu.
---

The TUI is an interactive terminal menu for Go Automate. It is a separate binary,
`go-automate-tui`, that shells out to the main `go-automate` CLI, so anything you can do
in the menu maps to a command you could run yourself.

## Launching

There are three ways to open the TUI:

- **Automatically** — run `go-automate` with no command in an interactive terminal. When
  the `go-automate-tui` binary is found, the menu opens.
- **Explicitly** — run the launcher command:

  ```bash
  go-automate tui
  ```

- **Directly** — run the TUI binary itself:

  ```bash
  go-automate-tui
  ```

Go Automate looks for the TUI binary next to the running `go-automate` executable first,
then on your `PATH`. If it cannot be found, the CLI continues normally.

## The menu

The main menu covers the most common tasks:

- **Bridge Serve** — start the [Home Assistant bridge](/running/) in the foreground.
- **Home Assistant** — open a submenu for lights, switches, input booleans and assist
  satellites.
- **Watch Entity** — watch an entity's state through the bridge.
- **Notify** — send a desktop notification.
- **Quit** — leave the menu.

Selecting a control action opens a small popup to collect its arguments, such as the
entity name to toggle or the message to announce, then runs the matching command.

## Search

Start typing to fuzzy-search the menu. Items also match helpful keywords, so typing
`lamp`, `monitor` or `announce` jumps to the right action. Vim-style shortcuts work too:
`:w` for watch, `:n` for notify, and `:q` to quit.

## Building the TUI

The TUI is built with [Bun](https://bun.sh) and [OpenTUI](https://github.com/sst/opentui).
Build it alongside the main binary:

```bash
make build_all
```

Or build just the TUI from the `tui/` directory:

```bash
cd tui
bun install
bun run build
```

During development you can run it with live reload:

```bash
cd tui
bun run dev
```

## Next steps

- Prefer the command line? See the [CLI](/using/cli/).
- Set up the [bridge](/running/) so Watch Entity stays cheap on the network.
