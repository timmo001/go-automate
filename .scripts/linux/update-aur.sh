#!/usr/bin/env bash
# Update the AUR package for go-automate-git
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
AUR_PACKAGE_NAME="go-automate-git"
AUR_REPO_URL="ssh://aur@aur.archlinux.org/${AUR_PACKAGE_NAME}.git"

# Detect if running in CI
IS_CI="${CI:-false}"

echo "Updating $AUR_PACKAGE_NAME AUR package..."
echo "Running in CI: $IS_CI"

# Setup SSH authentication if running in CI
if [ "$IS_CI" = "true" ]; then
  echo "Setting up SSH authentication for AUR..."

  if [ -z "$AUR_SSH_PRIVATE_KEY" ]; then
    echo "Error: AUR_SSH_PRIVATE_KEY environment variable not set"
    exit 1
  fi

  # Setup SSH
  mkdir -p ~/.ssh
  chmod 700 ~/.ssh

  # Write SSH key (use printf to preserve newlines)
  printf '%s\n' "$AUR_SSH_PRIVATE_KEY" > ~/.ssh/aur_rsa
  chmod 600 ~/.ssh/aur_rsa

  # Verify key was written correctly
  if ! grep -q "BEGIN.*PRIVATE KEY" ~/.ssh/aur_rsa; then
    echo "Error: SSH key format appears invalid"
    echo "Key should start with '-----BEGIN OPENSSH PRIVATE KEY-----' or similar"
    exit 1
  fi

  # Add AUR to known hosts
  ssh-keyscan -H aur.archlinux.org >> ~/.ssh/known_hosts 2>/dev/null

  # Configure SSH for AUR
  cat << 'EOF' > ~/.ssh/config
Host aur.archlinux.org
  IdentityFile ~/.ssh/aur_rsa
  User aur
  StrictHostKeyChecking accept-new
EOF
  chmod 600 ~/.ssh/config

  # Test SSH connection
  echo "Testing SSH connection to AUR..."
  if ! ssh -T aur@aur.archlinux.org 2>&1 | grep -q "successfully authenticated"; then
    echo "Warning: SSH authentication test did not return expected response"
    echo "Attempting to clone anyway..."
  fi

  # Create temporary directory for AUR repo
  AUR_REPO_PATH="$(mktemp -d)/aur-repo"

  # Clone AUR repository
  echo "Cloning AUR repository..."
  git clone "$AUR_REPO_URL" "$AUR_REPO_PATH"
else
  # Local mode: use ~/repos/aur location
  AUR_REPO_PATH="$HOME/repos/aur/$AUR_PACKAGE_NAME"

  # Check if AUR repository exists locally
  if [ ! -d "$AUR_REPO_PATH" ]; then
    echo "Error: AUR repository not found at $AUR_REPO_PATH"
    echo "Clone it first with: git clone $AUR_REPO_URL"
    exit 1
  fi
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

if [ "$IS_CI" = "true" ]; then
  # In CI with Arch container, script should already be running as builduser
  # Just set up the build directory
  export BUILDDIR="/tmp/makepkg-build"
  mkdir -p "$BUILDDIR"

  # Ensure we can access the working directory
  chmod 755 "$(dirname "$AUR_REPO_PATH")"

  # Run makepkg
  cd "$AUR_REPO_PATH"
  makepkg --printsrcinfo > .SRCINFO
else
  # Local mode: run directly
  cd "$AUR_REPO_PATH"
  makepkg --printsrcinfo > .SRCINFO
fi

# Show changes
echo ""
echo "Changes:"
git diff

# Check if there are any changes
if git diff --quiet; then
  echo ""
  echo "No changes detected. AUR package is already up to date."

  # Cleanup SSH key if in CI
  if [ "$IS_CI" = "true" ]; then
    rm -f ~/.ssh/aur_rsa
  fi

  exit 0
fi

# Commit and push changes
if [ "$IS_CI" = "true" ]; then
  echo ""
  echo "Committing and pushing to AUR..."

  git add -f PKGBUILD .SRCINFO
  git commit -m "Update to version $PKGVER"
  git push origin master

  echo "Successfully updated AUR package to version $PKGVER"

  # Cleanup SSH key
  rm -f ~/.ssh/aur_rsa
  rm -f ~/.ssh/config
else
  echo ""
  echo "Ready to commit and push to AUR"
  echo ""
  echo "To commit and push, run:"
  echo "  cd $AUR_REPO_PATH"
  echo "  git add PKGBUILD .SRCINFO"
  echo "  git commit -m 'Update to version $PKGVER'"
  echo "  git push"
fi
