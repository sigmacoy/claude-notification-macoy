# Troubleshooting

Common installation and runtime issues for macOS.

## OpenCode notifications missing or duplicated

- Confirm the host version with `opencode --version`. Published support: 1.18.33. The dual-API candidate is tested with 1.18.33,
  2.0.0 and 2.0.21; SDK 0.2.0 publication and registry qualification are pending. Update the notifications plugin using the guided OpenCode setup.
- Restart OpenCode after install/update/remove. Check that the same
  `OPENCODE_CONFIG_DIR` / `XDG_CONFIG_HOME` environment is used for setup and the
  host; the global plugin should be `plugins/agent-notifications.js` there.
- Choose desktop/webhook consent explicitly. Shared status/channel settings may
  disable delivery; webhook consent also needs an enabled endpoint.
- On macOS, use the installed executable's `setup-opencode permission-status`
  and `setup-opencode request-permission`. Check ClaudeNotifier in System
  Settings > Notifications and Focus. Missing/damaged attestation requires a
  verified native helper update, not an unverified fallback.
- Duplicate banners can come from OpenCode's native Agent/Permissions/Errors
  switches or another notification plugin. Disable the overlapping source.
- Silent alerts and no click-to-focus are expected for OpenCode. One-shot host
  exit can end asynchronous delivery before completion; use a persistent session
  when validating delivery.
- Ownership conflicts preserve foreign/edited files. Inspect the reported path;
  do not delete unrelated plugins. Interrupted transactions use
  `setup-opencode recover` before retrying setup.

See [OpenCode setup, supported targets and visual evidence](opencode-notifications.md).


## `Failed to add marketplace: ... its network source differs from the one declared for it in settings`

### Symptom

Running the installer to update fails during "Setting up marketplace..." with an error
mentioning the marketplace's network source differing from what's declared in settings.

### Why it happens

You installed before the project's GitHub repository was renamed from `claude-notifications-go`
to `agent-notifications`. Claude remembers which repo you originally added the marketplace
from and refuses to silently switch it to a different one under the same name — that's a
deliberate safety check, not a bug in Claude.

### Fix

The installer (`bin/bootstrap.sh`) detects this specific case and re-registers the marketplace
automatically; your saved notification settings are not affected. If you're still seeing this
error (e.g. from a cached copy of an older installer script), fix it manually and re-run the
installer:

```bash
/plugin marketplace remove claude-notifications-go
```

then run the [installer](INSTALLATION.md#quick-install-recommended) again.

## macOS: VS Code click-to-focus focuses the wrong window

### Symptom

Clicking a notification activates VS Code but raises the wrong window (or the last-active window) instead of the project-specific one.

### Why it happens

VS Code window focus requires **Screen Recording** permission (macOS 10.15+) to read window titles across all Spaces. Without it, the binary falls back to plain app activation.

### Fix

On first use the binary requests Screen Recording access automatically — a macOS dialog will appear. If you dismissed it:

1. Open **System Settings → Privacy & Security → Screen Recording**
2. Enable access for the `claude-notifications` binary (or the terminal running Claude)
3. Click the notification again

Once granted, the correct VS Code window will be raised even if it is on a different Space.

## The notification sound plays even though Do Not Disturb is on

### Symptom

Your desktop is in Do Not Disturb. The banner is correctly held back by the
desktop, but the plugin's sound plays anyway.

### Why it happens

The plugin plays its own audio cue in its own process, rather than asking the
notification server to play one - that is what makes per-status sounds and volume
control work. The notification server is explicitly asked not to play anything
(`suppress-sound`), and it has no say over a separate process's audio.

### Fix

Opt in to DND handling in the shared file selected by `config path`
(`"$NOTIFICATIONS_BIN" config path --json`; see the
[Configuration guide](CONFIGURATION.md)):

```json
{
  "notifications": {
    "respectDoNotDisturb": "silent"
  }
}
```

`"silent"` keeps the banner (it still lands in the notification centre) and drops
the sound; `"suppress"` drops both. The default is `"off"`, which is the old
behaviour.

### If it still plays

Detection deliberately fails open: anything it cannot read counts as "not in
DND", so the notification is delivered. Check, in order:

1. **Platform.** Detection for macOS Focus modes is supported. See
   [DO_NOT_DISTURB.md](DO_NOT_DISTURB.md).
2. **The value is spelled correctly.** An unrecognised value falls back to
   `"off"` and logs `Unknown respectDoNotDisturb value ...` rather than failing
   the config, so a typo looks exactly like the feature being switched off.
3. **Your desktop actually exposes the state.** Run the probe for macOS:

   ```bash
   defaults read com.apple.ncprefs dnd_prefs
   ```

   If this does not report DND while your desktop says it is on, the plugin has
   no source to read and will keep delivering.
4. **Which source fired.** With debug logging on, look for a `DND: active (...)`
   line in the plugin log. If the outcome looks right but the source is not the
   one you expect, the detection is wrong even though the result happens to match.

### Related: the terminal bell still rings

`"silent"` mutes the plugin's sound, not the terminal bell - BEL is a tab
indicator on Ghostty, tmux and other terminals, and whether it makes noise is
your terminal's own setting. Set `"terminalBell": false` to turn it off.