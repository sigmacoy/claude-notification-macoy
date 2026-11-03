package notifier

import (
	"strings"

	"github.com/777genius/agent-notifications/internal/analyzer"
	"github.com/777genius/agent-notifications/internal/notification"
	"github.com/777genius/agent-notifications/internal/sessionname"
)

// legacyDesktopPresentation is the presentation boundary that parses hook labels
// and translates analyzer status.
type legacyDesktopPresentation struct {
	notification.Content
	TimeSensitive bool
}

func legacyPresentation(status analyzer.Status, message, statusTitle string, sessionLabel bool) legacyDesktopPresentation {
	sessionName, gitBranch, _ := extractSessionInfo(message)

	folderName := ""
	if gitBranch != "" {
		parts := strings.SplitN(gitBranch, " ", 2)
		if len(parts) == 2 {
			folderName = parts[1]
		} else {
			folderName = gitBranch
		}
	} else if sessionName != "" {
		cleaned := sessionname.StripSessionLabel(sessionName)
		if cleaned != "" {
			folderName = cleaned
		} else if !sessionname.IsSessionLabel(sessionName) && sessionName != "unknown" {
			folderName = sessionName
		}
	}

	title := "Claude Code"
	if folderName != "" && folderName != "." {
		title = "Claude Code - " + folderName
	}

	body := statusToText(status)

	// No subtitle: removes redundant second line repeating folder name
	var subtitle string

	timeSensitive := isTimeSensitiveStatus(status)

	return legacyDesktopPresentation{Content: notification.Content{Title: title, Body: body, Subtitle: subtitle}, TimeSensitive: timeSensitive}
}

func hookPresentation(status analyzer.Status, content HookPresentation, statusTitle string, sessionLabel bool) legacyDesktopPresentation {
	name := ""
	if content.SessionName != "" && !sessionname.IsSessionLabel(content.SessionName) && content.SessionName != "unknown" {
		name = content.SessionName
	} else if content.Folder != "" && content.Folder != "." {
		name = content.Folder
	}

	title := "Claude Code"
	if name != "" {
		title = "Claude Code - " + name
	}

	body := statusToText(status)

	// No subtitle: removes redundant second line repeating folder name
	subtitle := ""

	return legacyDesktopPresentation{
		Content:       notification.Content{Title: title, Subtitle: subtitle, Body: body},
		TimeSensitive: isTimeSensitiveStatus(status),
	}
}

// HookDesktopContent shares literal hook presentation with trusted integrations.
// Delivery policy, consent and routing remain the caller's responsibility.
func HookDesktopContent(status analyzer.Status, content HookPresentation, statusTitle string, sessionLabel bool) notification.Content {
	return hookPresentation(status, content, statusTitle, sessionLabel).Content
}

// Platforms without a native subtitle keep literal context in the body.
func desktopBodyWithSubtitle(content notification.Content) string {
	if content.Subtitle == "" {
		return content.Body
	}
	if content.Body == "" {
		return content.Subtitle
	}
	return content.Subtitle + "\n" + content.Body
}

// statusToText maps analyzer status to simplified user-facing text.
func statusToText(status analyzer.Status) string {
	switch status {
	case analyzer.StatusTaskComplete, analyzer.StatusReviewComplete:
		return "Done"
	case analyzer.StatusQuestion, analyzer.StatusPlanReady, analyzer.StatusPermissionRequest:
		return "Waiting Input"
	case analyzer.StatusAgentError, analyzer.StatusSessionLimitReached, analyzer.StatusAPIError, analyzer.StatusAPIErrorOverloaded, analyzer.StatusAgentStopping:
		return "Error"
	default:
		return "Done"
	}
}

// appendSessionLabel keeps custom status titles intact while avoiding a redundant
// default completion word when the session name already identifies the task.
func appendSessionLabel(status analyzer.Status, statusTitle, sessionName string) string {
	if status == analyzer.StatusTaskComplete && statusTitle == "✅ Completed" && strings.TrimSpace(sessionName) != "" {
		statusTitle = "✅"
	}
	return statusTitle + " [" + sessionName + "]"
}

func joinContext(first, second string) string {
	if second == "." {
		second = ""
	}
	if first == "" {
		return second
	}
	if second == "" {
		return first
	}
	return first + " · " + second
}

func shortenTitle(text string, limit int) string {
	runes := []rune(strings.Join(strings.Fields(text), " "))
	if len(runes) > limit {
		return string(runes[:limit-1]) + "…"
	}
	return string(runes)
}
