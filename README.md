# Claude Notification Macoy

A customized, clean desktop macOS notification system and lifecycle hook handler for **Claude Code**. It delivers native audio and visual alerts when agent tasks finish, errors occur, or user input is required — featuring **zero emojis**, clean titles, simplified status bodies, and click-to-focus window management.

![Claude Code](https://img.shields.io/badge/Claude%20Code-Plugin-blue?style=flat-square&logo=anthropic)
![macOS](https://img.shields.io/badge/macOS-Supported-black?style=flat-square&logo=apple)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)
![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg?style=flat-square)

> **Platform**: macOS only. This plugin is designed exclusively for macOS and uses native macOS notification infrastructure (ClaudeNotifier.app, terminal-notifier). Windows and Linux are not supported.

---

## Features

- **No Emojis**: Removed emoji badges for a minimal, professional aesthetic.
- **Concise Title Formatting**: Notification titles display as `Claude Code - {Project Name}` (or `Claude Code` if no project name is detected).
- **Minimal Status Bodies**: Messages are streamlined to essential states:
  - `Done` (for completed tasks and reviews)
  - `Waiting Input` (for questions, plans, and permission approvals)
  - `Error` (for agent stops, limits, or API errors)
- **Session Label Cleanup**: Automatically filters out random internal session hashes (e.g. `bold 06ddb8f7`, `cosmic ba9deea3`) and boilerplate greetings (`"How can I help..."`).
- **Claude Logo Icon**: Ships with the dedicated Claude icon logo (`claude_icon_v2.png`) bundled directly into the notification alert.
- **Click-to-Focus**: Clicking any alert immediately brings your active terminal window or IDE to the front.
- **Subagent Noise Suppression**: Prevents notification spam by suppressing background subagent completions while still alerting on main agent milestones.
- **Universal macOS Binary Included**: Comes pre-built with a universal Mach-O binary supporting both **Apple Silicon (M1/M2/M3/M4)** and **Intel Macs**.
- **✨ Duplicate Notification Prevention**: New `play-notification` binary uses atomic file-lock deduplication to ensure sounds play only once per session per hook event (Stop, Notification, etc.) within a 3-second window — fixes the "double notification" issue when both Stop and Notification hooks fire for the same interaction.

---

## Notification Output Format

| Trigger Event | Notification Title | Notification Body | Audio Cue |
| :--- | :--- | :--- | :--- |
| **Task Complete** | `Claude Code - [Project Name]` | `Done` | `task-complete.mp3` |
| **Review Complete** | `Claude Code - [Project Name]` | `Done` | `review-complete.mp3` |
| **User Question** | `Claude Code - [Project Name]` | `Waiting Input` | `question.mp3` |
| **Plan Ready** | `Claude Code - [Project Name]` | `Waiting Input` | `plan-ready.mp3` |
| **Permission Request** | `Claude Code - [Project Name]` | `Waiting Input` | `question.mp3` |
| **Error / Interruption** | `Claude Code - [Project Name]` | `Error` | `error.mp3` |

---

## Installation

Users should download or clone this GitHub repository directly to their machine.

### Step 1: Download or Clone the Repository

Open your terminal and clone the repository:

```bash
git clone https://github.com/sigmacoy/claude-notification-macoy.git
cd claude-notification-macoy
```

*(Alternatively, download the ZIP from GitHub by clicking **Code > Download ZIP**, and extract it into a folder).*

---

### Step 2: Verify Permissions and Binaries

Run the setup script to verify executable permissions and pre-compiled binaries:

```bash
chmod +x setup.sh bin/*
./setup.sh
```

> **Note**: Pre-compiled universal binaries for macOS (Apple Silicon + Intel) and the native `ClaudeNotifier.app` bundle are already included in the `bin/` directory. The `play-notification` binary is also included for duplicate sound prevention.

---

### Step 3: Register the Plugin in Claude Code

Open Claude Code and add your downloaded repository folder as a local marketplace:

```text
/plugin marketplace add /absolute/path/to/claude-notification-macoy
/plugin install claude-notifications-go@claude-notification-macoy
/claude-notifications-go:init
```

*(Replace `/absolute/path/to/claude-notification-macoy` with the real path where you cloned or extracted the repository, e.g. `pwd`).*

---

### Step 4: Restart Claude Code

Restart Claude Code once for the lifecycle hooks to activate. Your custom notifications will now trigger on events!

---

## Configuration & Preferences

### Interactive Settings

To configure notification sounds, volumes, and alert toggles interactively inside Claude Code:

```text
/claude-notifications-go:settings
```

This interactive menu allows you to:
- Test and preview notification sounds (`sound-preview`)
- Adjust volume (0.0 to 1.0)
- Enable or disable alerts for specific trigger events
- Toggle Click-to-Focus
- Configure deduplication window (default: 3 seconds)

### Configuration File

To view the path to your active configuration file:

```bash
./bin/claude-notifications config path
```

Configurations are stored at:
- **macOS**: `~/.claude/claude-notifications-go/config.json`

Example `config.json`:
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

---

## The Duplicate Notification Fix (Technical Details)

### The Problem

Claude Code fires **two hooks** for the same user interaction:
1. **Stop** hook — fires when the agent finishes responding
2. **Notification** hook — fires for permission prompts, questions, etc.

Both hooks would independently play notification sounds, resulting in **duplicate audio** — the first one often cutting off or sounding "bad."

### The Solution: `play-notification` Binary

A dedicated Go binary (`bin/play-notification`) handles sound playback with **atomic file-lock deduplication**:

```
┌─────────────────────────────────────────────────────────────┐
│  Hook Event (Stop, Notification, SubagentStop, etc.)        │
└─────────────────────┬───────────────────────────────────────┘
                      ▼
┌─────────────────────────────────────────────────────────────┐
│  Notifier.SendDesktop()                                     │
│  • Passes sessionID + hookEvent to sendOptions              │
└─────────────────────┬───────────────────────────────────────┘
                      ▼
┌─────────────────────────────────────────────────────────────┐
│  playSoundDetached() spawns:                                │
│  ./play-notification -sound <file> -session <id> -hook <evt>│
└─────────────────────┬───────────────────────────────────────┘
                      ▼
┌─────────────────────────────────────────────────────────────┐
│  play-notification (deduplication logic):                   │
│  1. Build lock key: "sessionID-hookEvent" (e.g., "abc-Stop")│
│  2. Check ~/.tmp/claude-play-notification-<key>.lock        │
│  3. If lock exists AND age < 3s → EXIT (suppressed)         │
│  4. Else: create lock atomically (O_EXCL), play sound       │
│  5. Background goroutine cleans up lock after window + 1s   │
└─────────────────────────────────────────────────────────────┘
```

**Result**: Only the first hook event plays the sound. Subsequent hooks within 3 seconds for the same session+event are silently suppressed.

### Hook Events Protected

| Hook Event | Lock Key Suffix | Use Case |
|------------|-----------------|----------|
| `Stop` | `-Stop` | Task completion |
| `Notification` | `-Notification` | Permission prompts, questions |
| `SubagentStop` | `-SubagentStop` | Subagent completion |
| `TeammateIdle` | `-TeammateIdle` | Team collaboration |

---

## Testing Notifications

You can verify that notifications work properly directly from your terminal:

```bash
# Test a completion hook (Stop event)
echo '{"session_id":"test","transcript_path":""}' | ./bin/claude-notifications handle-hook Stop

# Test manual notification
./bin/claude-notifications notify

# Test play-notification binary directly (with deduplication)
./bin/play-notification -sound "sounds/task-complete.mp3" -session "test-123" -hook "Stop"

# Test deduplication - run twice within 3 seconds, second is suppressed
./bin/play-notification -sound "sounds/task-complete.mp3" -session "test-123" -hook "Stop"
./bin/play-notification -sound "sounds/task-complete.mp3" -session "test-123" -hook "Stop"  # suppressed
```

---

## macOS Permissions

On macOS 10.15+, macOS requires notification permissions for terminal applications:

1. Open **System Settings > Notifications**.
2. Locate **Agent Notifications** or your terminal application (e.g., Terminal, iTerm2, Ghostty, Warp, Alacritty, VS Code).
3. Ensure **Allow Notifications** is enabled and **Alert style** is set to **Banners** or **Alerts**.

---

## 🔨 Building From Source (Developers)

If you modify the Go source code in `cmd/` or `internal/` and wish to rebuild the binaries:

1. **Prerequisites**: Go 1.21+ installed (`go version`).
2. Run `make build`:
   ```bash
   make build
   ```
3. Test your build:
   ```bash
   ./bin/claude-notifications version
   ./bin/play-notification -version
   ```

### Cross-Platform Build (macOS only)

```bash
# Build optimized release binaries for macOS (both architectures)
make build-all
```

---

## Project Structure

```text
claude-notification-macoy/
├── .claude-plugin/
│   ├── marketplace.json        # Marketplace definition for Claude Code
│   └── plugin.json             # Plugin definition and command registrations
├── assets/                     # Notification icon assets
│   ├── icon.png
│   └── icon_v2.png             # Custom Claude logo icon
├── bin/
│   ├── ClaudeNotifier.app      # Standalone native macOS notification bundle
│   ├── claude-notifications    # Universal binary (Apple Silicon + Intel)
│   ├── play-notification       # Sound playback with deduplication
│   ├── hook-wrapper.sh         # POSIX wrapper for hook execution
│   └── setup.sh                # Helper setup script
├── cmd/
│   ├── claude-notifications/   # Main CLI and hook entry point
│   ├── play-notification/      # Dedicated sound player with deduplication
│   ├── sound-preview/          # Sound preview utility
│   └── list-sounds/            # List available sounds
├── commands/                   # Claude Code slash commands (/init, /settings, etc.)
├── config/
│   └── config.json             # Default configuration template
├── hooks/
│   └── hooks.json              # Hook event bindings (PreToolUse, Stop, Notification, SubagentStop, TeammateIdle)
├── internal/
│   ├── notifier/               # Desktop notification presentations and formatters
│   ├── summary/                # Notification body and status text summarizer
│   ├── sessionname/            # Session label and prefix filters
│   ├── hooks/                  # Hook handling logic with deduplication
│   ├── audio/                  # Audio playback (malgo/miniaudio)
│   ├── platform/               # Platform-specific utilities
│   └── ...
├── sounds/                     # Bundled audio prompt cues (.mp3)
│   ├── task-complete.mp3
│   ├── review-complete.mp3
│   ├── question.mp3
│   ├── plan-ready.mp3
│   └── error.mp3
├── claude_icon_v2.png          # High-resolution Claude icon
├── Makefile                    # Build and development recipes
├── README.md                   # Documentation and user guide
└── setup.sh                    # Top-level setup verification script
```

---

## Sharing With Others

### Option 1: Share the Repository (Recommended)

1. Push this repository to GitHub (or your preferred Git host)
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

1. Run `make build-all` to create optimized binaries in `dist/`
2. Share the `dist/` folder containing:
   - `claude-notifications-darwin-amd64` (Intel Macs)
   - `claude-notifications-darwin-arm64` (Apple Silicon Macs)
   - `play-notification-darwin-amd64`
   - `play-notification-darwin-arm64`
3. Recipients place binaries in their `bin/` folder and run setup

### Option 3: Package as a Plugin Release

See the release workflows in `.github/workflows/` for automated signing, notarization, and distribution via GitHub Releases.

---