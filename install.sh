#!/bin/bash
# Install Filex desktop entry and icon for the current user.
# Re-run after rebuilding to update the binary path.

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BINARY="$SCRIPT_DIR/filex"

if [ ! -f "$BINARY" ]; then
    echo "Binary not found at $BINARY — building..."
    (cd "$SCRIPT_DIR" && go build -o filex .)
fi

mkdir -p ~/.local/share/applications
mkdir -p ~/.local/share/icons/hicolor/scalable/apps

# Install icon
cp "$SCRIPT_DIR/assets/filex.svg" ~/.local/share/icons/hicolor/scalable/apps/filex.svg

# Install desktop entry with correct binary path
sed "s|Exec=.*|Exec=$BINARY|" "$SCRIPT_DIR/assets/filex.desktop" \
    > ~/.local/share/applications/filex.desktop

# Refresh caches
update-desktop-database ~/.local/share/applications/ 2>/dev/null || true
gtk-update-icon-cache ~/.local/share/icons/hicolor/ 2>/dev/null || true

echo "Installed. You can now search for 'Filex' in your app launcher and pin it to the taskbar."
