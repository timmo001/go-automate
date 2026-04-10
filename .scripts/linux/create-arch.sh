#!/bin/bash

set -euo pipefail

# If running as root in CI container, re-run as non-root build user
if [ "$(id -u)" -eq 0 ] && id -u builduser >/dev/null 2>&1; then
  echo "Switching to builduser for packaging..."
  chown -R builduser:builduser "$(pwd)"
  exec sudo --preserve-env=VERSION -u builduser -H bash "$0" "$@"
fi

# Ensure required tools for Arch packaging are installed (makepkg)
if ! command -v makepkg >/dev/null 2>&1; then
  echo "makepkg not found, attempting installation..."
  if command -v pacman >/dev/null 2>&1; then
    sudo pacman -Syu --noconfirm --needed base-devel
  elif command -v apt-get >/dev/null 2>&1; then
    echo "makepkg is Arch-specific. Please run this script on an Arch-based system, or use the CI containerized build." >&2
    exit 1
  else
    echo "Unsupported or unknown package manager. Please run on Arch or install makepkg (pacman base-devel)." >&2
    exit 1
  fi
fi

# Check if binary exists
if [ ! -f "go-automate" ]; then
  echo "go-automate not found, please build the application first"
  exit 1
fi

# Create build directory
mkdir -p build/arch
cd build/arch

# Copy necessary files
cp ../../go-automate go-automate
cp ../../LICENSE LICENSE
cp ../../.scripts/linux/PKGBUILD.binary PKGBUILD
cp ../../.scripts/linux/go-automate-home-assistant-bridge.service go-automate-home-assistant-bridge.service
cp ../../.scripts/linux/arch-package.install arch-package.install

# Sanitize VERSION for Arch pkgver
ARCH_PKGVER=$(echo "$VERSION" | sed 's/[-+]/./g')
export ARCH_PKGVER

echo "ARCH_PKGVER: $ARCH_PKGVER"

# Generate new sha256sums and update PKGBUILD
makepkg -g >new_sums.txt
sed -i '/^sha256sums=(/,/^)/d' PKGBUILD
awk '
  /source=/ { print; while ((getline line < "new_sums.txt") > 0) print line; next }
  { print }
' PKGBUILD >PKGBUILD.new && mv PKGBUILD.new PKGBUILD
rm new_sums.txt

# Build package
makepkg -f --noconfirm

# Move package to dist directory
mkdir -p ../../dist
mv *.pkg.tar.zst ../../dist/

cd ../..
rm -rf build/arch

echo "Package created successfully!"
echo "Install with: yay -U dist/go-automate-${ARCH_PKGVER}-1-x86_64.pkg.tar.zst"
