# Claude Notification Macoy

Claude Notification Macoy is a macOS notification system for Claude Code. It sends native desktop alerts with audio when agent tasks finish, errors occur, or user input is required.

This plugin runs on macOS only. It uses native macOS notification infrastructure. Windows and Linux are not supported.

## Installation (macOS)

### Step 1: Download the Repository

Open Terminal and run:

```bash
git clone https://github.com/sigmacoy/claude-notification-macoy.git
cd claude-notification-macoy
```

Or download the ZIP from GitHub (Code > Download ZIP) and extract it.

### Step 2: Verify Binaries

Run the setup script:

```bash
chmod +x setup.sh bin/*
./setup.sh
```

Pre-compiled universal binaries for Apple Silicon and Intel Macs are included in the `bin/` directory. The `ClaudeNotifier.app` bundle and `play-notification` binary are also included.

### Step 3: Register the Plugin in Claude Code

Open Claude Code and run:

```text
/plugin marketplace add /absolute/path/to/claude-notification-macoy
/plugin install claude-notifications-go@claude-notification-macoy
/claude-notifications-go:init
```

Replace `/absolute/path/to/claude-notification-macoy` with the real path (use `pwd` to get it).

### Step 4: Restart Claude Code

Restart Claude Code. The lifecycle hooks will activate. Custom notifications will trigger on events.

## Features

- **No emojis** — Clean, professional appearance
- **Concise titles** — Format: `Claude Code - [Project Name]` (or `Claude Code` if no project name)
- **Minimal status bodies** — Three states only:
  - `Done` — Task or review completed
  - `Waiting Input` — Question, plan, or permission request
  - `Error` — Agent stop, limit, or API error
- **Session label cleanup** — Removes internal session hashes and boilerplate greetings
- **Claude logo icon** — Bundled `claude_icon_v2.png` appears in notifications
- **Click-to-focus** — Click a notification to bring your terminal or IDE to front
- **Subagent noise suppression** — Background subagent completions do not trigger alerts
- **Universal macOS binary** — Supports Apple Silicon (M1/M2/M3/M4) and Intel Macs
- **Duplicate notification prevention** — `play-notification` binary uses file-lock deduplication. Sounds play once per session per hook event within a 3-second window.

## Notification Output Format

| Trigger Event | Title | Body | Audio |
|---------------|-------|------|-------|
| Task Complete | `Claude Code - [Project Name]` | `Done` | task-complete.mp3 |
| Review Complete | `Claude Code - [Project Name]` | `Done` | review-complete.mp3 |
| User Question | `Claude Code - [Project Name]` | `Waiting Input` | question.mp3 |
| Plan Ready | `Claude Code - [Project Name]` | `Waiting Input` | plan-ready.mp3 |
| Permission Request | `Claude Code - [Project Name]` | `Waiting Input` | question.mp3 |
| Error / Interruption | `Claude Code - [Project Name]` | `Error` | error.mp3 |

## Configuration

### Interactive Settings

Run inside Claude Code:

```text
/claude-notifications-go:settings
```

The menu lets you:
- Test notification sounds
- Adjust volume (0.0 to 1.0)
- Enable or disable alerts for specific events
- Toggle Click-to-Focus
- Set deduplication window (default: 3 seconds)

### Configuration File

View the config path:

```bash
./bin/claude-notifications config path
```

Config file location on macOS:
`~/.claude/claude-notifications-go/config.json`

Example configuration:

```json
{
  "notifications": {
    "desktop": {
      "enabled": true,
      "sound": true,
      "volume": 1.0,
      "appIcon": "${CLAUDE_PLUGIN_ROOT}/claude_icon_v2.png",
      "clickToFocus": true,
      "terminalBundleId": ""
    },
    "suppressForSubagents": true,
    "dedupWindowSeconds": 3
  }
}
```

## How Duplicate Prevention Works

### The Problem

Claude Code fires two hooks for the same interaction:
1. **Stop** hook — Agent finishes responding
2. **Notification** hook — Permission prompt or question appears

Both hooks played sounds. This caused duplicate audio.

### The Solution

The `play-notification` binary handles sound playback with atomic file-lock deduplication:

1. Hook event occurs (Stop, Notification, SubagentStop, TeammateIdle)
2. Notifier calls `play-notification` with session ID and hook event
3. Binary builds lock key: `sessionID-hookEvent` (example: `abc-Stop`)
4. Binary checks `~/.tmp/claude-play-notification-<key>.lock`
5. If lock exists and is less than 3 seconds old — exit (suppressed)
6. Otherwise — create lock atomically, play sound
7. Background process cleans up lock after window plus 1 second

Result: Only the first hook event plays the sound. Subsequent hooks within 3 seconds for the same session and event are suppressed.

### Protected Hook Events

| Hook Event | Lock Key Suffix | Use Case |
|------------|-----------------|----------|
| `Stop` | `-Stop` | Task completion |
| `Notification` | `-Notification` | Permission prompts, questions |
| `SubagentStop` | `-SubagentStop` | Subagent completion |
| `TeammateIdle` | `-TeammateIdle` | Team collaboration |

## Testing Notifications

Run from Terminal:

```bash
# Test completion hook (Stop event)
echo '{"session_id":"test","transcript_path":""}' | ./bin/claude-notifications handle-hook Stop

# Test manual notification
./bin/claude-notifications notify

# Test play-notification directly (with deduplication)
./bin/play-notification -sound "sounds/task-complete.mp3" -session "test-123" -hook "Stop"

# Test deduplication — run twice within 3 seconds, second is suppressed
./bin/play-notification -sound "sounds/task-complete.mp3" -session "test-123" -hook "Stop"
./bin/play-notification -sound "sounds/task-complete.mp3" -session "test-123" -hook "Stop"  # suppressed
```

## macOS Permissions

On macOS 10.15 and later, Terminal applications need notification permissions:

1. Open **System Settings > Notifications**
2. Find **Agent Notifications** or your terminal app (Terminal, iTerm2, Ghostty, Warp, Alacritty, VS Code)
3. Enable **Allow Notifications**
4. Set **Alert style** to **Banners** or **Alerts**

## Building From Source (Developers)

### Prerequisites

- Go 1.21 or later (`go version`)

### Build

```bash
make build
```

### Test Build

```bash
./bin/claude-notifications version
./bin/play-notification -version
```

### Cross-Platform Build (macOS only)

```bash
make build-all
```

This creates optimized release binaries for both architectures.

## Project Structure

```
claude-notification-macoy/
├── .claude-plugin/
│   ├── marketplace.json        # Marketplace definition
│   └── plugin.json             # Plugin definition and commands
├── assets/                     # Notification icon assets
│   ├── icon.png
│   └── icon_v2.png             # Custom Claude logo icon
├── bin/
│   ├── ClaudeNotifier.app      # Native macOS notification bundle
│   ├── claude-notifications    # Universal binary (Apple Silicon + Intel)
│   ├── play-notification       # Sound player with deduplication
│   ├── hook-wrapper.sh         # POSIX hook wrapper
│   └── setup.sh                # Setup verification script
├── cmd/
│   ├── claude-notifications/   # Main CLI and hook entry point
│   ├── play-notification/      # Sound player with deduplication
│   ├── sound-preview/          # Sound preview utility
│   └── list-sounds/            # List available sounds
├── commands/                   # Slash commands (/init, /settings)
├── config/
│   └── config.json             # Default configuration template
├── hooks/
│   └── hooks.json              # Hook event bindings
├── internal/
│   ├── notifier/               # Desktop notification formatters
│   ├── summary/                # Notification body summarizer
│   ├── sessionname/            # Session label filters
│   ├── hooks/                  # Hook handling with deduplication
│   ├── audio/                  # Audio playback
│   └── platform/               # Platform utilities
├── sounds/                     # Audio cues (.mp3)
│   ├── task-complete.mp3
│   ├── review-complete.mp3
│   ├── question.mp3
│   ├── plan-ready.mp3
│   └── error.mp3
├── claude_icon_v2.png          # High-resolution Claude icon
├── Makefile                    # Build recipes
├── README.md                   # This file
└── setup.sh                    # Top-level setup script
```

## Sharing With Others

### Option 1: Share the Repository (Recommended)

1. Push to GitHub
2. Others clone and install:
   ```bash
   git clone https://github.com/YOUR_USERNAME/claude-notification-macoy.git
   cd claude-notification-macoy
   chmod +x setup.sh bin/*
   ./setup.sh
   ```
3. In Claude Code:
   ```text
   /plugin marketplace add /absolute/path/to/claude-notification-macoy
   /plugin install claude-notifications-go@claude-notification-macoy
   /claude-notifications-go:init
   ```

### Option 2: Distribute Pre-Built Binaries

1. Run `make build-all` to create binaries in `dist/`
2. Share the `dist/` folder:
   - `claude-notifications-darwin-amd64` (Intel Macs)
   - `claude-notifications-darwin-arm64` (Apple Silicon Macs)
   - `play-notification-darwin-amd64`
   - `play-notification-darwin-arm64`
3. Recipients place binaries in their `bin/` folder and run setup

### Option 3: Package as a Plugin Release

See `.github/workflows/` for automated signing, notarization, and GitHub Releases distribution.
