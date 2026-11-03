# Claude Notification Macoy

A customized desktop notification system and lifecycle hook handler for Claude Code. It delivers native audio and visual alerts when agent tasks finish, errors occur, or user input is required.

![Claude Code](https://img.shields.io/badge/Claude%20Code-Plugin-blue?style=flat-square&logo=anthropic)
![macOS](https://img.shields.io/badge/macOS-Supported-black?style=flat-square&logo=apple)
![Windows](https://img.shields.io/badge/Windows-Supported-0078D6?style=flat-square&logo=windows)
![Linux](https://img.shields.io/badge/Linux-Supported-FCC624?style=flat-square&logo=linux)
![CLI](https://img.shields.io/badge/CLI-Tool-lightgrey?style=flat-square)
![Notifications](https://img.shields.io/badge/Notifications-Desktop-success?style=flat-square)

---

## Overview

- **Concise Formatting**: Formats notification titles as `Claude Code - [Project Name]`.
- **Minimal Bodies**: Simplifies notification messages to `Done` or `Waiting Input`.
- **Session Cleanup**: Removes random session hashes and redundant folder subtitles.
- **Click-to-Focus**: Clicking an alert automatically focuses the active terminal or workspace.
- **Audio Prompts**: Plays non-intrusive sound cues for completions, questions, and approval requests.

---

## Installation

Step-by-step installation instructions separated by operating system.

### macOS

#### 1. Prerequisites
- Claude Code CLI installed and available on `PATH`.
- macOS 10.15 (Catalina) or later (Apple Silicon and Intel supported).

#### 2. Install Plugin
Run these slash commands inside your Claude Code session:

```text
/plugin marketplace add 777genius/agent-notifications
/plugin install claude-notifications-go@claude-notifications-go
/claude-notifications-go:init
```

*Alternatively, install via Terminal using the bootstrap script:*
```bash
(set -o pipefail; curl -fsSL https://agent-notifications.com/install.sh | bash -s -- --product claude --desktop)
```

#### 3. Enable Permissions
When prompted, allow notifications for **Agent Notifications**. You can verify permissions under:
`System Settings > Notifications > Agent Notifications`

#### 4. Restart
Restart your Claude Code session to load hooks.

---

### Windows

#### 1. Prerequisites
- Claude Code CLI installed.
- **Git Bash** (included with [Git for Windows](https://git-scm.com/download/win)). Do not install via WSL.

#### 2. Install Plugin
Open Git Bash and run:

```bash
(set -o pipefail; curl -fsSL https://agent-notifications.com/install.sh | bash -s -- --product claude --desktop)
```

*Or install directly within a Claude Code session:*
```text
/plugin marketplace add 777genius/agent-notifications
/plugin install claude-notifications-go@claude-notifications-go
/claude-notifications-go:init
```

#### 3. Enable Notifications
Ensure notifications and banners are enabled for your terminal application under:
`Settings > System > Notifications`

#### 4. Restart
Restart your terminal and relaunch Claude Code.

---

### Linux

#### 1. Prerequisites
- Claude Code CLI installed.
- D-Bus notification daemon and `libnotify-bin`:

```bash
# Debian / Ubuntu
sudo apt-get install -y libnotify-bin alsa-utils

# Fedora / RHEL
sudo dnf install -y libnotify alsa-utils

# Arch Linux
sudo pacman -S libnotify alsa-utils
```

#### 2. Install Plugin
Run the automated installer:

```bash
(set -o pipefail; curl -fsSL https://agent-notifications.com/install.sh | bash -s -- --product claude --desktop)
```

*Or install within Claude Code:*
```text
/plugin marketplace add 777genius/agent-notifications
/plugin install claude-notifications-go@claude-notifications-go
/claude-notifications-go:init
```

#### 3. Restart
Restart your Claude Code session to initialize hooks.

---

## Configuration

To adjust sounds, volume, and notification preferences, run inside Claude Code:

```text
/claude-notifications-go:settings
```

To display the path to your active configuration file:

```bash
claude-notifications config path
```

---

## Notification Formats

| Trigger Event | Notification Title | Notification Body |
| :--- | :--- | :--- |
| **Task Complete** | `Claude Code - [Project Name]` | `Done` |
| **User Question** | `Claude Code - [Project Name]` | `Waiting Input` |
| **Plan Ready** | `Claude Code - [Project Name]` | `Waiting Input` |
| **Permission Request** | `Claude Code - [Project Name]` | `Waiting Input` |
| **Error / Interruption** | `Claude Code - [Project Name]` | `Error` |
