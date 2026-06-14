# 🎛️ Go Automate

A utility to run common tasks.

I call this app with keyboard shortcuts and patched apps in linux to trigger automations in my home using [Home Assistant](https://home-assistant.io).

## Documentation

Full documentation lives at **<https://go-automate.timmo.dev>**:

- [Install](https://go-automate.timmo.dev/install/) — AUR package or build from source
- [Configuration](https://go-automate.timmo.dev/configuration/) — Home Assistant URL and token
- [Running the Bridge](https://go-automate.timmo.dev/running/) — the systemd user service
- [CLI](https://go-automate.timmo.dev/using/cli/) and [TUI](https://go-automate.timmo.dev/using/tui/)
- [Home Assistant](https://go-automate.timmo.dev/using/home-assistant/), [Watching Entities](https://go-automate.timmo.dev/using/watching/) and [Notifications](https://go-automate.timmo.dev/using/notifications/)
- [Reference](https://go-automate.timmo.dev/reference/) — commands and the bridge protocol

## Quick start

```bash
go build              # build the binary
./go-automate         # run it (interactive setup on first launch)
```

See [Install](https://go-automate.timmo.dev/install/) for the AUR package and
[Configuration](https://go-automate.timmo.dev/configuration/) to set your Home Assistant
URL and token.

## Contributing

The documentation site lives in [`docs/`](docs/). Issues and pull requests are welcome on
[GitHub](https://github.com/timmo001/go-automate).
