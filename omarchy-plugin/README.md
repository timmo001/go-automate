# Go Automate for Omarchy

An Omarchy bar widget and panel driven by Go Automate commands. It runs
commands that emit bar JSON, shows active values in the bar, and presents all
configured modules in a keyboard-filterable panel. The initial configuration
format supports Home Assistant statuses and actions without embedding personal
entities in the plugin.

## Requirements

- Omarchy Quattro
- `go-automate` available on `PATH`
- A running Go Automate Home Assistant bridge for bridge-backed watchers
- A JSON modules file created and owned by the user

## Install

Review the repository, then add the plugin:

```bash
omarchy plugin add https://github.com/timmo001/omarchy-go-automate.git
```

Accept the prompt to enable the plugin during installation. Set `modulesFile`
to the absolute path of your JSON configuration. `modules.example.json` shows
the supported base shape without containing personal entity IDs.

## Modules

The modules file contains a `modules` array. Every module requires `id` and
`command`. Commands must emit one JSON object per update with `text`, `tooltip`,
and `class` fields. Use `stream: true` for a long-running watcher or `interval`
in milliseconds for polling.

Optional fields include `group`, `label`, `icon`, `action`, `panelOnly`,
`background`, `hideUnavailable`, `hideClasses`, `barHideClasses`, `activeText`,
`inactiveText`, `barIconOnly`, `severityClasses`, and `colors`. The plugin runs
`command` and `action` through `bash -lc`, so review every configured command.

The recommended watcher is:

```bash
go-automate ha bridge watch entity light.example --bar-json --hide-off
```

## Use

Select the widget to open its panel. Type to filter modules, use Up and Down to
move through the list, and press Escape to clear the filter or close the panel.
Right-click the widget or press Ctrl+R in the panel to refresh polled modules.

The plugin preserves the `timmo.home-assistant` shell IPC target:

```bash
omarchy-shell timmo.home-assistant toggle
```

## Settings

- `modulesFile`: absolute path to the user-owned JSON modules file
- `primaryOnly`: show the widget only on the selected output
- `primaryOutput`: optional output name used when `primaryOnly` is enabled

## Update

```bash
omarchy plugin update timmo.home-assistant
```

## Remove

```bash
omarchy plugin remove timmo.home-assistant
```

Removing the plugin does not remove the user-owned modules file.

## Validate from source

```bash
omarchy plugin validate .
```

## Security

This plugin runs unsandboxed inside `omarchy-shell` when enabled. Review its
source and your modules file before using it.

It reads the configured modules file and starts every configured command. Go
Automate watchers read Home Assistant credentials from the Go Automate config
and connect through the local bridge. Actions may access Home Assistant over
the network. The plugin does not install software, run privileged commands, or
write Omarchy configuration.
