#!/usr/bin/env sh

set -e

EXT_DIR="$HOME/.vscode/extensions/gorth-language-support"
SRC_DIR="vscode"

echo "Removing existing Gorth VS Code extension..."
rm -rf "$EXT_DIR"

echo "Copying new Gorth VS Code extension..."
cp -R "$SRC_DIR" "$EXT_DIR"

echo "Done."
