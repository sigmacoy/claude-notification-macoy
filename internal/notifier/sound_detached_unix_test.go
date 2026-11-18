//go:build !windows

package notifier

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/777genius/agent-notifications/internal/config"
	"github.com/777genius/agent-notifications/internal/logging"
)

// TestPlaySoundDetached_ChildStderrReachesDebugLog guards issue #131: the
// detached play-sound child must write its diagnostics into the same
// notification-debug.log as the parent instead of /dev/null.
func TestPlaySoundDetached_ChildStderrReachesDebugLog(t *testing.T) {
	// Set CLAUDE_PLUGIN_ROOT so findPlayNotificationBinary can locate the binary
	_, thisFile, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(thisFile)))
	t.Setenv("CLAUDE_PLUGIN_ROOT", projectRoot)

	logRoot := t.TempDir()
	if _, err := logging.InitLogger(logRoot); err != nil {
		t.Fatalf("InitLogger: %v", err)
	}
	t.Cleanup(func() { _ = logging.Close() })

	sound := filepath.Join(t.TempDir(), "sound.mp3")
	if err := os.WriteFile(sound, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	originalExecCommand := execCommand
	execCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("sh", "-c", "echo 'play-sound: child diagnostic' >&2")
	}
	t.Cleanup(func() { execCommand = originalExecCommand })

	cfg := config.DefaultConfig()
	cfg.Notifications.Desktop.Sound = true
	n := New(cfg)
	t.Cleanup(func() { _ = n.Close() })

	n.playSoundDetached(sound, sendOptions{})

	logPath := filepath.Join(logRoot, "notification-debug.log")
	deadline := time.Now().Add(5 * time.Second)
	for {
		data, _ := os.ReadFile(logPath)
		if strings.Contains(string(data), "play-sound: child diagnostic") {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("child stderr never reached the debug log; log contents:\n%s", data)
		}
		time.Sleep(20 * time.Millisecond)
	}
}
