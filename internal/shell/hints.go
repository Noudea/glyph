package shell

import "strings"

func (m Model) hintText() string {
	switch m.mode {
	case ModeLauncher:
		return "type to filter · ↑/↓ move · enter run · esc/" + m.shortcutsHint(commandLauncherOpen, "ctrl+p/ctrl+k/alt+p") + " close"
	case ModeMain:
		hints := []string{
			m.workspaceHint(),
			m.primaryShortcut(actionWorkspaceToggle, "ctrl+w") + " toggle workspace",
			m.shortcutsHint(commandLauncherOpen, "ctrl+p/ctrl+k/alt+p") + " command palette",
		}
		if m.err != "" {
			hints = append([]string{"error: " + m.err}, hints...)
		}
		if appHint := m.appHint(); appHint != "" {
			hints = append([]string{appHint}, hints...)
		}
		return strings.Join(hints, " · ")
	default:
		return ""
	}
}

func (m Model) appHint() string {
	if module, ok := m.activeModule(); ok {
		return module.Hint()
	}
	return ""
}
