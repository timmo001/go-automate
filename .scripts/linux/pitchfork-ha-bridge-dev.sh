#!/usr/bin/env bash
set -euo pipefail

service="go-automate-home-assistant-bridge.service"
socket="${XDG_RUNTIME_DIR:-/tmp}/go-automate/home-assistant.sock"

restore_service() {
	if /usr/bin/systemctl --user cat "$service" >/dev/null 2>&1; then
		/usr/bin/systemctl --user start "$service" >/dev/null 2>&1 || true
	fi
}

cleanup() {
	trap - EXIT INT TERM HUP
	restore_service
}

trap cleanup EXIT
trap 'exit 143' INT TERM HUP

/usr/bin/systemctl --user stop "$service" >/dev/null 2>&1 || true
fuser -TERM -k "$socket" >/dev/null 2>&1 || true
sleep 1
fuser -KILL -k "$socket" >/dev/null 2>&1 || true

mise run build
./go-automate ha bridge serve
