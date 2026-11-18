# Do Not Disturb

On macOS, the notification service plays the configured per-status sound. The
helper converts it to a short PCM resource with the configured volume already
applied, then submits an ordinary notification. macOS applies Focus and per-app
notification settings to the banner and sound together. No Focus-file access or
Full Disk Access is required.

## Configuration

In the shared file selected by `config path`
(`"$NOTIFICATIONS_BIN" config path --json`; see the
[Configuration guide](CONFIGURATION.md)):

```json
{
  "notifications": {
    "respectDoNotDisturb": "silent"
  }
}
```

| Value | Behaviour while DND is active |
|---|---|
| `"off"` (default) | Current behaviour. DND state is never queried, so the option costs nothing when it is not in use. |
| `"silent"` | The banner is still delivered — it lands in the notification centre and shows when DND lifts — but the plugin's own sound is not played. |
| `"suppress"` | Nothing is delivered: no banner, no sound. |

Webhooks are unaffected in all three modes. DND is a statement about *this*
desktop's attention, not about your team's Slack channel.

DND is read at delivery time, after any `notifyDelaySeconds` wait, so switching
DND on during the grace period is honoured.

### Why the default is `"off"`

Playing audio through DND is the bug, so "respect DND" arguably belongs on by
default. Two things argue against flipping it:

- A silent behaviour change on upgrade is a surprise. `notifyOnlyWhenUnfocused`
  shipped as `false` for the same reason.
- macOS Focus modes apply automatically to native notifications. The plugin does
  not query Focus state, so `"suppress"` cannot remove a notification from
  Notification Center on macOS.

### Terminal bell

On macOS the plugin does not send a terminal bell, because an audible BEL can
bypass Focus.

## Detection fails open

Detection returns "in DND" **only** on positive confirmation. Any error,
unsupported desktop, missing daemon, timeout, or unparseable value is treated as
"unknown", which means the notification and its sound are delivered exactly as
before.

This bias is deliberate: a probe error or unknown state keeps the existing
delivery behavior. A confirmed positive result in `"suppress"` mode intentionally
skips desktop delivery; a false positive could therefore hide a notification.

## What is detected on macOS

Focus applies automatically to native notifications, including custom sounds,
regardless of `respectDoNotDisturb`. The plugin does not query Focus state, so
`"suppress"` cannot remove a notification from Notification Center on macOS.
Normal notifications do not request time-sensitive interruption.

Custom sounds are decoded to signed 16-bit PCM CAF files in `~/Library/Sounds`,
with the configured volume applied to the samples. Resources are cached by
source content and volume. The source must be a regular file of at most 32 MiB
and shorter than 30 seconds. Invalid or unsupported sounds do not fall back to a
default beep or separate Go playback. macOS controls the alert output device and
system alert volume in addition to the plugin's configured gain.

`respectDisplaySleep` continues to mute the custom sound when the display sleeps.
Older native helpers do not understand custom audio arguments and remain silent;
update both the sender and native helper together.