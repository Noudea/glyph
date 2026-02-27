package shell

import "strings"

func (m Model) hintText() string {
	switch m.mode {
	case ModeLauncher:
		return "type to filter · ↑/↓ move · enter run · esc/" + m.shortcutsHint(commandLauncherOpen, "ctrl+p/ctrl+k/alt+p") + " close"
	case ModeMain:
		hints := []string{
			m.workspaceHint(),
			m.shortcutsHint(commandLauncherOpen, "ctrl+p/ctrl+k/alt+p") + " command palette",
		}
		if m.err != "" {
			hints = append([]string{"error: " + m.err}, hints...)
		}
		return strings.Join(hints, " · ")
	default:
		return ""
	}
}
