# Go Automate

A utility to run common tasks.

I call this app with keyboard shortcuts and patched apps in linux to trigger automations in my home using [Home Assistant](https://home-assistant.io).

## Setup

1. Install Go
1. Install Arch packaging tools if you plan to build a local package (`base-devel`, `yay`)
1. Create `~/.config/go-automate/config.yml` with your Home Assistant URL and long-lived access token, or run `go-automate` once interactively to create it.

Example config:

```yaml
homeassistant:
  url: http://homeassistant.local:8123
  token: your-long-lived-access-token
```

## Build

1. Clone this repo
1. Run `go build`

## Arch Package

1. Run `make create_arch`
1. Install the generated package with `yay -U dist/go-automate-<version>-1-x86_64.pkg.tar.zst`

The package installs:

- `/usr/bin/go-automate`
- `/usr/lib/systemd/user/go-automate-home-assistant-bridge.service`

The package post-install script will:

- globally enable `go-automate-home-assistant-bridge.service` for future user logins,
- try to start it immediately for the installing user when that user has an active systemd user session,
- leave Home Assistant config in each user's `~/.config/go-automate/config.yml`.

If the service is not running yet in the current session, run:

```bash
systemctl --user daemon-reload
systemctl --user enable --now go-automate-home-assistant-bridge.service
```

Check service status with:

```bash
systemctl --user status go-automate-home-assistant-bridge.service
```
