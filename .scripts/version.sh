#!/usr/bin/env bash
# Compute the build version from git, mirroring the AUR pkgver style.
# Falls back to a revision count when no tags are reachable.
set -euo pipefail

version=$(git describe --long --tags --abbrev=7 2>/dev/null \
  || printf "r%s.%s" "$(git rev-list --count HEAD)" "$(git rev-parse --short=7 HEAD)")

printf "%s" "$version" | sed "s/^v//;s/\([^-]*-g\)/r\1/;s/-/./g"
