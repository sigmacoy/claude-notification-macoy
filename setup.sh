#!/bin/bash
# setup.sh - Verify and set up Claude Notification Macoy plugin

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "=========================================="
echo " Claude Notification Macoy - Setup"
echo "=========================================="
echo ""

# Ensure bin directory scripts are executable
chmod +x "$SCRIPT_DIR"/bin/* 2>/dev/null || true

# Check if binary exists
if [ -f "$SCRIPT_DIR/bin/claude-notifications" ]; then
    echo "✓ Pre-compiled universal binary detected in bin/claude-notifications"
else
    echo "⚙️ Compiling binary from source..."
    cd "$SCRIPT_DIR" && go build -ldflags="-s -w" -trimpath -o bin/claude-notifications ./cmd/claude-notifications
    echo "✓ Built bin/claude-notifications"
fi

echo "✓ Permissions and scripts verified"
echo ""
echo "=========================================="
echo " Setup Complete!"
echo "=========================================="
echo ""
echo "To use this plugin inside Claude Code:"
echo ""
echo "Option A (Local installation):"
echo "  /plugin marketplace add $SCRIPT_DIR"
echo "  /plugin install claude-notifications-go@claude-notification-macoy"
echo ""
echo "Option B (From GitHub once pushed):"
echo "  /plugin marketplace add <your-github-username>/claude-notification-macoy"
echo "  /plugin install claude-notifications-go@claude-notification-macoy"
echo ""
echo "Then restart Claude Code or run:"
echo "  /claude-notifications-go:init"
echo ""
