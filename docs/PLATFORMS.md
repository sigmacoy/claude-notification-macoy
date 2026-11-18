[Back to README](../README.md)

# Platform Support

**Supported platform:**
- macOS (Intel & Apple Silicon)

**Self-contained runtime after setup:**
- ✅ Binaries auto-download from GitHub Releases
- ✅ Pure Go - no C compiler needed
- ✅ All libraries bundled
- ✅ Works offline after first setup

> **Note**: This plugin is macOS only. Windows and Linux are not supported.

### Click-to-Focus (macOS)

Clicking a notification activates your terminal window. Auto-detects terminal.

**macOS** — via AX API with bundle ID detection:

| Terminal | Focus method |
|----------|-------------|
| Ghostty | Exact tab focus via Ghostty AppleScript, with AXDocument fallback |
| VS Code / Insiders / Cursor | AXTitle (focus-window subcommand) |
| iTerm2 | Exact tab/pane targeting via iTerm2 Python API when available, otherwise app-level iTerm activation |
| Warp | Exact window/tab/pane via `WARP_FOCUS_URL` (`warp://session/<uuid>`), AXTitle fallback on older Warp |
| kitty, WezTerm, Alacritty, Hyper, Apple Terminal | AXTitle (focus-window subcommand) |
| Any other (custom `terminalBundleId`) | AXTitle (focus-window subcommand) |

**Multiplexers:** tmux and kitty (macOS only) — click switches to the correct pane/tab.

**iTerm2 note:** to open the exact iTerm2 tab or split pane, enable `iTerm2 > Settings > General > Magic > Enable Python API`. If you just toggled it, restart iTerm2 once. Without the Python API, the plugin falls back to app-level iTerm activation instead of exact tab targeting.

See **[Click-to-Focus Guide](CLICK_TO_FOCUS.md)** for configuration details.