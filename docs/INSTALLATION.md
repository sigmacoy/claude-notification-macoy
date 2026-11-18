[Back to README](../README.md)

# Installation (macOS Only)

> **Platform**: This plugin is macOS only. Windows and Linux are not supported.

## Prerequisites

- Install the selected agent CLIs first. The installer does not impose CLI version minimums, maximums or exact-version allowlists.
- `curl`, Bash, and a working `python3` or `node` on `PATH` for installer metadata and checksum validation.
- Git for Claude marketplace setup, `tar` for Codex bundles, and `unzip` for the macOS desktop helper.
- Python is additionally needed for optional iTerm2 exact tab/pane targeting; Node.js can satisfy the installer requirement but cannot replace the iTerm2 Python API.

### Quick Install (Recommended)

The short setup loader handles release lookup and validation internally, then downloads the installer from its qualified immutable source snapshot. Run it and choose your agents:

```bash
(set -o pipefail; curl -fsSL https://agent-notifications.com/install.sh | bash)
```

For automation or terminals without a controlling TTY, choose explicitly and preserve download failures in the exit status:

```bash
(set -o pipefail; curl -fsSL https://agent-notifications.com/install.sh | bash -s -- --product claude)
```

Use `claude`, `codex`, or `both` for Claude/Codex. The other product IDs are `opencode`, `gemini`, `qoder`, `cursor` and `pi`. The selected host CLI must already be on `PATH`; this installs notifications only.

Select agents together in the guided setup, or pass one comma-separated list to the common loader:

```bash
(set -o pipefail; curl -fsSL https://agent-notifications.com/install.sh | bash -s -- --products claude,codex --desktop)
```

`--desktop` and `--webhook` apply to each selected integration, with consent saved independently. `--desktop` selects desktop only, `--webhook` selects webhooks only, and both select both channels. Configure webhook URLs separately. Fresh installs default to desktop on and webhooks off. Updates without channel flags preserve each selected agent's complete saved decision, including both channels off. Use `--desktop=false --webhook=false` to disable those channels explicitly while retaining the integration. `--skip-agent-notify` skips notification-tool setup for selected Claude/Codex clients. Existing registrations and consent remain unchanged.

Add `,pi` or use `--product pi --desktop` for Pi. Install its official standalone CLI first. Pi uses the same common installer, with separate desktop/webhook consent and local-TUI sound/click behavior. [Pi installation, configuration and lifecycle](pi-notifications.md).

The installer shows recent safe checkpoints for each selected client above its live status, including its current operation and elapsed waiting time. Saved registration remains distinct from completed setup; a timeout or failed activation keeps that distinction visible. The bottom status shows the current stage and observed client counts. Resizing reduces history first; very small terminals switch to plain output.

Use `--activity compact` (the default), `--activity expanded` to also print the bounded session history at completion, or `--activity off` to hide history. `--activity-limit 1` through `16` controls visible checkpoints per client when space permits. These options require an installer with presentation support; an older pinned release reports a compatibility error before installation. `--ui plain`, `NO_COLOR`, `TERM=dumb`, redirected output and `BOOTSTRAP_VERBOSE=1` use append-only text. Activity contains fixed checkpoint labels; detailed stage output stays in the private diagnostic log. Observers stop before a prompt reads input, and setup success does not verify notification delivery.

After installation:

- **Claude:** restart Claude. Optionally run `/claude-notifications-go:settings` to configure sounds.
- **Codex:** start Codex, run `/hooks`, then review and trust the installed hooks. The installer registers them automatically; no JSON editing or manual registration command is needed. Trust approval remains yours.
- **Both:** complete both steps above.
- **Qoder:** install the Qoder CLI and make `qoder` or `qodercli` available on `PATH`, then restart Qoder after installation. [Native setup and limits](QODER.md).
- **Cursor CLI:** restart the selected agent CLI to load its Stop observer. Cursor IDE is outside this integration. [Qualified hosts and limits](cursor-notifications.md).
- **Gemini CLI:** restart Gemini to reload its hooks. [Setup and limits](gemini-notifications.md).
- **Pi:** restart Pi or run `/reload` to load the native extension. [Modes, sound and terminal focus](pi-notifications.md).
- **OpenCode:** restart OpenCode to load its global plugin. On Mac, explicitly grant notification permission. OpenCode bootstrap requires release v1.46.0 or newer; see [OpenCode setup, channels and limits](opencode-notifications.md).

When the selected release supports the portable wizard, the installer also registers the `agent-notifications` MCP server and `agent-notifications` skill for the selected Claude/Codex clients. Restart each selected client and open a new session before checking its MCP tools. The agent can call `notify` during a task, not only after a Stop hook. Check the read-only `notification_status` first: when `navigation.capability` is `eligible`, prefer `navigation: "required"` to keep the configured return to the chat, including for informational notifications. Eligibility uses this call's client context and configuration; it does not verify that the app, profile or chat can actually open. Use `navigation: "none"` only for an alert that does not need a return to the chat, including setups without a desktop route. Never silently downgrade a request that requires a return to the chat. `notify` returns `submitted` when the OS accepted the request, which does not prove that a banner was visible or the chat opened. Desktop notification permission and a visible test send must be checked separately.

The installer distinguishes saved registration, failed activation and activation that was not started. If setup fails, keep the selected profile and the private log if `Diagnostics` reports a path. Follow its single next action after resolving the cause; a saved activation command uses the acquired package. `--skip-agent-notify` is an explicit hooks-only choice and does not remove an existing MCP installation. Automatic setup keeps absent MCP clients off when an existing managed portable installation or MCP binding is detected. Skipping setup alone does not record a permanent per-client opt-out. Use an explicit `--agent-notify` request to add that client later.

Codex requires a published stable plugin release v1.42.0 or newer. The installer downloads matching source and binaries, respects `CODEX_HOME`, and keeps a permanent runtime copy there. It reports an error if no supported release is published yet.

> If installation fails, use [manual Claude installation](#manual-install) or [manual Codex registration](CODEX.md#manual-codex-registration), depending on the product.

### Manual Install

<details>
<summary>Step-by-step installation inside Claude (if bootstrap doesn't work)</summary>

Run these slash commands in the Claude chat, not in your system terminal:

```text
# 1) Add marketplace
/plugin marketplace add 777genius/agent-notifications
# 2) Install plugin
/plugin install claude-notifications-go@claude-notifications-go
# 3) Restart Claude
# 4) Download binary
/claude-notifications-go:init
# 5) (Optional) Configure sounds and settings
/claude-notifications-go:settings
```

> **Compatibility:** `claude-notifications-go` is the frozen Claude marketplace, plugin, and command namespace. The public product is **Agent Notifications**, but changing these technical identifiers breaks existing installations and updates. See [Claude plugin identity compatibility](CLAUDE_PLUGIN_IDENTITY.md).

</details>

> Having issues with installation? See [Troubleshooting](troubleshooting.md).

### Updating

Run the [secure install command](#quick-install-recommended) again and choose the product(s) you want to update.

For OpenCode, rerun with explicitly chosen desktop/webhook flags; the idempotent install action updates the registered runtime. Restart OpenCode. For Claude, restart Claude. For Codex, restart Codex and inspect `/hooks`; changed hook definitions may need trust approval again. The installer refreshes the Codex runtime and registration automatically. Existing foreign hooks and shared settings in the file selected by `config path` are preserved.

<details>
<summary>Manual Claude update (if bootstrap didn't work)</summary>

Claude also periodically checks for plugin updates automatically. Binaries are updated on the next hook invocation when a version mismatch is detected.

To update manually via Claude UI:

1. Run `/plugin`, select **Marketplaces**, choose `claude-notifications-go`, then select **Update marketplace**
2. Select **Installed**, choose `claude-notifications-go`, then select **Update now**

If the binary auto-update didn't work (e.g. no internet at the time), run `/claude-notifications-go:init` to download it manually. If hook definitions changed in the new version, restart Claude to apply them.

</details>

### Uninstalling

**OpenCode:** use the installer's printed **Remove** command, then restart OpenCode. On macOS it invokes the installed executable directly. Removal revokes its consent before deleting owned files and preserves other consumers. [Detailed removal/recovery commands](opencode-notifications.md#change-channels-remove-or-recover).

**Claude:**

```text
/plugin uninstall claude-notifications-go@claude-notifications-go
```

Optionally also remove the marketplace registration: `/plugin marketplace remove claude-notifications-go`.

**Codex:** remove the hooks and runtime registered by this installer manually. Use the same Codex home selected during setup: the explicit `--codex-home` path, otherwise `CODEX_HOME`, otherwise `~/.codex`. In that directory, delete only this installer's entries from `hooks.json`, then remove the `claude-notifications-go` directory. Preserve hooks registered for other tools.

**Configuration:** uninstalling does not delete your saved settings. Run `agent-notifications config path` to find the active file, and remove it yourself if you no longer want it.

### Reading the installation result

The installer ends with a summary for the selected agents. Installed means that registration succeeded; it does not confirm notification delivery. Restart the agents, review and trust Codex entries with `/hooks`, and send a test notification. If a notification does not arrive, check OS permissions and notification settings. Skipped agent-notify setup is reported separately from installed automatic hooks. Failed stages show a compact cause and retain earlier stage output in a private durable log outside `TMPDIR` on POSIX systems. The default `agent-notifications-installer-diagnostics` directory sits beside the managed control root, allowing fresh installation to provision its own transaction state. Its parent directory has mode 0700 and the log has mode 0600. Each retained stream holds at most 64 KiB, with an explicit truncation marker; the session log is bounded to 2 MiB. Read the existing log for details without running installation again. `BOOTSTRAP_VERBOSE=1` displays printable diagnostic details at the end of a requested install. If logging fails, the installer reports that limitation and keeps raw command output out of the default display.

### Restore a pruned Claude cache

The bootstrap installer checks for managed files lost when Claude prunes an old versioned plugin cache. If exact recovery is possible, interactive setup offers to restore those files before preparing the normal installation plan. Declining leaves the installation unchanged. Explicit noninteractive installation performs the same repair automatically. Existing files with changed content or permissions remain conflicts; repair never overwrites them or resets configuration. If the read-only repair preview cannot complete, bootstrap leaves the decision to the normal installer's existing validation, reinstall and journal recovery. Failures after repair starts stop setup instead of being ignored.

Recovery uses matching owned copies or verified payloads retained outside the Claude cache by newer installers. The payloads are optional support, not another installation or a policy source. Native application links require the recorded native generation to remain intact. An older installation without an exact copy or retained payload still requires recovery from a trusted matching package; the installer does not bypass its fingerprint check.

To run only this repair, using an installed helper that supports it:

```sh
claude-notifications setup-products repair-cache
# Explicit unattended repair, with the same identity checks:
claude-notifications setup-products repair-cache --yes
```

This restores the existing integrations. Removing a fully abandoned consumer is a separate operation described below.

### Recover an entirely absent installer consumer

Ordinary install and removal remain strict about every managed asset. For an abandoned secondary runtime, use the exact consumer ID, physical runtime root, installation ID and positive generation from that installation's ownership record. A missing registration alone does not establish an orphan. Keep retained assets intact and preview first:

```sh
claude-notifications internal-install-runtime --recover-orphan-consumer \
  --consumer 'codex:/absolute/TEST/codex/hooks.json' \
  --runtime-root /absolute/TEST/codex/claude-notifications-go \
  --control-root /absolute/TEST/config/agent-notifications \
  --expected-installation-id INSTALLATION_ID --expected-generation 18 \
  --dry-run --json
```