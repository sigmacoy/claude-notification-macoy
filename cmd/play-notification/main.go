package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/777genius/agent-notifications/internal/audio"
	"github.com/777genius/agent-notifications/internal/logging"
	"github.com/777genius/agent-notifications/internal/platform"
)

const (
	appName    = "claude-play-notification"
	appVersion = "1.0.0"
)

func main() {
	// Parse flags
	var (
		soundFile   = flag.String("sound", "", "Path to sound file (required)")
		volume      = flag.Float64("volume", 1.0, "Volume 0.0-1.0")
		device      = flag.String("device", "", "Audio device name (empty for default)")
		sessionID   = flag.String("session", "", "Session ID for deduplication")
		hookEvent   = flag.String("hook", "", "Hook event name for deduplication")
		dedupWindow = flag.Int("dedup-window", 3, "Deduplication window in seconds")
		showVersion = flag.Bool("version", false, "Show version")
		help        = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("%s v%s\n", appName, appVersion)
		os.Exit(0)
	}

	if *help || *soundFile == "" {
		printUsage()
		os.Exit(1)
	}

	// Initialize logging to stderr (for debug output)
	logger := logging.NewLoggerFromFile(os.Stderr)
	logger.SetPrefix(fmt.Sprintf("%s pid=%d", appName, os.Getpid()))

	// Resolve sound file path
	soundPath := *soundFile
	if !filepath.IsAbs(soundPath) {
		// Try relative to plugin root if CLAUDE_PLUGIN_ROOT is set
		if pluginRoot := os.Getenv("CLAUDE_PLUGIN_ROOT"); pluginRoot != "" {
			soundPath = filepath.Join(pluginRoot, "sounds", soundPath)
		}
	}

	// Check if sound file exists
	if !platform.FileExists(soundPath) {
		logger.Error("Sound file not found: %s", soundPath)
		os.Exit(1)
	}

	// Deduplication check
	if *sessionID != "" {
		// Build dedup key: session + hook event (if provided)
		dedupKey := *sessionID
		if *hookEvent != "" {
			dedupKey = *sessionID + "-" + *hookEvent
		}

		lockPath := filepath.Join(platform.TempDir(), fmt.Sprintf("claude-play-notification-%s.lock", dedupKey))

		// Check for recent duplicate
		if platform.FileExists(lockPath) {
			age := platform.FileAge(lockPath)
			if age >= 0 && age < int64(*dedupWindow) {
				logger.Debug("Duplicate notification suppressed: session=%s hook=%s age=%ds", *sessionID, *hookEvent, age)
				os.Exit(0)
			}
			// Lock file exists but is old - remove it so we can create a fresh one
			_ = os.Remove(lockPath)
		}

		// Create lock file
		created, err := platform.AtomicCreateFile(lockPath)
		if err != nil || !created {
			logger.Warn("Failed to create dedup lock: created=%v err=%v", created, err)
		}

		// Schedule cleanup of old locks
		go func() {
			time.Sleep(time.Duration(*dedupWindow+1) * time.Second)
			_ = os.Remove(lockPath)
		}()
	}

	// Initialize audio player
	player, err := audio.NewPlayer(*device, *volume)
	if err != nil {
		logger.Error("Failed to init audio player: %v", err)
		os.Exit(1)
	}
	defer func() { _ = player.Close() }()

	// Play sound
	logger.Debug("Playing sound: %s (volume=%.2f)", soundPath, *volume)
	if err := player.Play(soundPath); err != nil {
		logger.Error("Failed to play sound: %v", err)
		os.Exit(1)
	}

	logger.Debug("Notification played successfully")
}

func printUsage() {
	fmt.Printf(`%s v%s - Play notification sounds with deduplication

Usage:
  %s -sound <file> [options]

Required:
  -sound <file>       Path to sound file (relative to plugin sounds/ dir or absolute)

Options:
  -volume <float>     Volume 0.0-1.0 (default: 1.0)
  -device <name>      Audio device name (default: system default)
  -session <id>       Session ID for deduplication
  -hook <name>        Hook event name for deduplication (e.g., Stop, Notification)
  -dedup-window <sec> Deduplication window in seconds (default: 3)
  -version            Show version
  -help               Show this help

Examples:
  # Play a sound with session-based deduplication
  %s -sound "complete.mp3" -session "abc123" -hook "Stop"

  # Play without deduplication
  %s -sound "/path/to/sound.mp3" -volume 0.5

Environment:
  CLAUDE_PLUGIN_ROOT  Plugin root directory (for resolving relative sound paths)
`, appName, appVersion, appName, appName, appName)
}