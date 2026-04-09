# Go Automate

A utility to run common tasks.

I call this app with keyboard shortcuts and patched apps in linux to trigger automations in my home using [Home Assistant](https://home-assistant.io).

## Setup

1. Install Go
1. Install Arch packaging tools if you plan to build a local package (`base-devel`, `yay`)

## Build

1. Clone this repo
1. Run `go build`

## Arch Package

1. Run `make create_arch`
1. Install the generated package with `yay -U dist/go-automate-<version>-1-x86_64.pkg.tar.zst`
