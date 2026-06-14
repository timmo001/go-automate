---
title: LLMs
description: Feed the Go Automate documentation to an agent using llms.txt.
---

The Go Automate documentation is published in formats an agent can read directly, following
the [llms.txt](https://llmstxt.org) convention. Point an agent at one of these files to give
it the docs as context.

## Available files

- [`/llms.txt`](https://go-automate.timmo.dev/llms.txt) - an index of the documentation with links to each section. Use this when you want the agent to pick what it needs.
- [`/llms-full.txt`](https://go-automate.timmo.dev/llms-full.txt) - the entire documentation in a single file. Use this for the most complete context.
- [`/llms-small.txt`](https://go-automate.timmo.dev/llms-small.txt) - a trimmed version for smaller context windows.

These files are generated from the docs at build time, so they stay in sync with the rest of
the site.

## How to use them

- Paste one of the URLs above into an agent and ask it to read the page.
- Use the actions in the page header on any docs page to copy the page as Markdown, then
  paste it straight into an agent. This gives the agent the content directly, so it does not
  need to fetch the page over the web.
