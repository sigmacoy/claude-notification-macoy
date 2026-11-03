package sessionname

import (
	"fmt"
	"strings"
)

// Lists for friendly name generation (same as bash version)
var adjectives = []string{
	"bold", "brave", "bright", "calm", "clever",
	"cool", "cosmic", "crisp", "daring", "eager",
	"fair", "fancy", "fast", "gentle", "glad",
	"grand", "happy", "kind", "lively", "lucky",
	"merry", "noble", "proud", "quick", "quiet",
	"rapid", "smart", "solid", "swift", "warm",
	"wise", "witty", "zesty", "agile", "alert",
}

var nouns = []string{
	"bear", "bird", "cat", "deer", "eagle",
	"fish", "fox", "hawk", "lion", "owl",
	"star", "moon", "sun", "wind", "wave",
	"tree", "river", "mountain", "ocean", "cloud",
	"tiger", "wolf", "dragon", "phoenix", "falcon",
	"comet", "galaxy", "planet", "nova", "meteor",
	"forest", "canyon", "valley", "peak", "storm",
}

// GenerateSessionName generates a friendly single-word name from a session ID (UUID).
// Returns a deterministic name like "cat" or "eagle".
//
// Args:
//   - sessionID: UUID string (e.g., "73b5e210-ec1a-4294-96e4-c2aecb2e1063")
//
// Returns:
//   - Friendly name string (e.g., "peak")
func GenerateSessionName(sessionID string) string {
	// Return "unknown" if no session ID
	if sessionID == "" || sessionID == "unknown" {
		return "unknown"
	}

	// Remove dashes and convert to lowercase
	cleanID := strings.ToLower(strings.ReplaceAll(sessionID, "-", ""))

	// Use first 8 hex chars as seed for word selection
	if len(cleanID) < 8 {
		// Fallback for short IDs
		return "unknown"
	}

	seed := cleanID[0:8]

	// Combine adjectives and nouns into a single pool for more variety
	allWords := append(adjectives, nouns...)

	// Convert hex to decimal for array indexing
	index := hexToInt(seed) % len(allWords)

	return allWords[index]
}

// GenerateSessionLabel generates a friendly name with session ID prefix.
// Returns a string like "bold 06ddb8f7" for better session identification.
func GenerateSessionLabel(sessionID string) string {
	name := GenerateSessionName(sessionID)
	if name == "unknown" {
		return "unknown"
	}

	// Extract first 8 chars of UUID (before first dash)
	prefix := sessionID
	if idx := strings.Index(sessionID, "-"); idx != -1 {
		prefix = sessionID[:idx]
	}
	if len(prefix) > 8 {
		prefix = prefix[:8]
	}

	return name + " " + prefix
}

// hexToInt converts hex string to int (takes first 6 characters for safety)
func hexToInt(hex string) int {
	if len(hex) > 6 {
		hex = hex[0:6]
	}

	var result int
	if _, err := fmt.Sscanf(hex, "%x", &result); err != nil {
		return 0 // Return 0 on parse error
	}
	return result
}

// IsSessionLabel checks if s matches a generated session label like "bold 06ddb8f7" or "cosmic ba9deea3"
func IsSessionLabel(s string) bool {
	fields := strings.Fields(strings.TrimSpace(s))
	if len(fields) != 2 {
		return false
	}
	word, prefix := fields[0], fields[1]
	if len(prefix) < 4 || len(prefix) > 8 {
		return false
	}
	for _, c := range prefix {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	for _, a := range adjectives {
		if a == word {
			return true
		}
	}
	for _, n := range nouns {
		if n == word {
			return true
		}
	}
	return false
}

// StripSessionLabel removes a leading session label (e.g. "cosmic ba9deea3") from text
func StripSessionLabel(text string) string {
	fields := strings.Fields(strings.TrimSpace(text))
	if len(fields) >= 2 {
		candidate := fields[0] + " " + fields[1]
		if IsSessionLabel(candidate) {
			if len(fields) == 2 {
				return ""
			}
			return strings.Join(fields[2:], " ")
		}
	}
	return strings.TrimSpace(text)
}
