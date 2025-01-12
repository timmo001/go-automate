# Go Automate

A utility to run common tasks.

I call this app with keyboard shortcuts and patched apps in linux to trigger automations in my home using [Home Assistant](https://home-assistant.io).

## Setup

1. Install go
1. Set up your go workspace and make sure that your `GOPATH` is set correctly.

```zsh
export GOPATH=$HOME/go
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
```

## Build and Install

1. Clone this repo
1. Run `go build`
1. Run `go install`
