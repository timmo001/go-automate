---
title: Shell completions
description: Enable tab completion for Go Automate commands, subcommands, and flags in zsh, bash, and fish.
---

Go Automate ships shell completion for its whole command tree. Press <kbd>Tab</kbd> to
complete `ha`, `bridge`, `watch` and the rest, their aliases, and every flag such as
`--bar-json`.

## Arch package

The [Arch package](/install/) installs the completion scripts for you, so completion works
in zsh, bash, and fish as soon as you open a new shell. There is nothing else to do.

## Build from source

When you build from source, generate the script for your shell and load it.

### zsh

```bash
go-automate completion zsh > "${fpath[1]}/_go-automate"
```

Start a new shell, or run `source ~/.zshrc`, to pick it up. To try it in the current shell
without installing:

```bash
source <(go-automate completion zsh)
```

### bash

```bash
go-automate completion bash > ~/.local/share/bash-completion/completions/go-automate
```

Or load it for the current shell only:

```bash
source <(go-automate completion bash)
```

### fish

```bash
go-automate completion fish > ~/.config/fish/completions/go-automate.fish
```

## Aliases

If you alias `go-automate`, most shells reuse the completion through the alias. In zsh, an
alias like `alias ga=go-automate` completes automatically. If your setup does not expand it,
tell zsh to reuse the completion function:

```bash
compdef _go-automate ga
```

## Notes

- Completion never contacts Home Assistant. It only lists commands, aliases, and flags, so
  it stays fast and works before you have configured a URL or token.
- Completion is static: it does not yet suggest live entity IDs. That is planned as a
  follow-up.
