# AGENTS.md

## Deleting or renaming a file in the package

A file crossed out from the plugin package can "resurface" again from three accumulation points (lesson from issue #490: the deleted `skills/agent-notify` resurfaced right next to `skills/agent-notifications`). If you delete or rename a shipped file, clean up all three places and cover them with tests:

1. **Copy-back of the same version** (`bin/bootstrap.sh`) — the list of restored fingerprints must not return the deleted file into the updated plugin root.

2. **Leftovers in the plugin root** (cleanup after refresh in `bin/bootstrap.sh`) — the updated root must match the package: the Claude Code overlay update does not delete crossed-out files on its own.

3. **Ledger of managed-runtime** (`internal/installruntime`, `cmd/claude-notifications/install_runtime.go`) — the owned-record is retired by the same transaction; user-edited or unowned legacy is not silently deleted — the installation must refuse (see `TestEmbeddedSkillRename*`).

Verification: `bin/release_channel_native_test.sh` (CI Ubuntu) + tests in `cmd/claude-notifications` (`TestEmbeddedSkillRename*`).