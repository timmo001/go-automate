#!/usr/bin/env bash
# Update the AUR package for go-automate-git
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
AUR_REPO_PATH="$HOME/repos/aur/go-automate-git"

echo "Updating go-automate-git AUR package..."

# Check if AUR repository exists
if [ ! -d "$AUR_REPO_PATH" ]; then
  echo "Error: AUR repository not found at $AUR_REPO_PATH"
  echo "Clone it first with: git clone ssh://aur@aur.archlinux.org/go-automate-git.git"
  exit 1
fi

# Generate version from git
cd "$REPO_ROOT"

# Get version information
LAST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "0.1.0")
REV_COUNT=$(git rev-list --count HEAD)
SHORT_HASH=$(git rev-parse --short=7 HEAD)

# Remove 'v' prefix if present
PREFIX_VERSION="${LAST_TAG#v}"

# Generate pkgver
PKGVER="${PREFIX_VERSION}.r${REV_COUNT}.g${SHORT_HASH}"

echo "Generated version: $PKGVER"

# Copy PKGBUILD to AUR repository
cp "$SCRIPT_DIR/PKGBUILD" "$AUR_REPO_PATH/PKGBUILD"

# Update pkgver in PKGBUILD
cd "$AUR_REPO_PATH"
sed -i "s/^pkgver=.*/pkgver=${PKGVER}/" PKGBUILD

# Generate .SRCINFO
echo "Generating .SRCINFO..."
makepkg --printsrcinfo > .SRCINFO

# Show changes
echo ""
echo "Changes:"
git diff

echo ""
echo "Ready to commit and push to AUR"
echo ""
echo "To commit and push, run:"
echo "  cd $AUR_REPO_PATH"
echo "  git add PKGBUILD .SRCINFO"
echo "  git commit -m 'Update to version $PKGVER'"
echo "  git push"
