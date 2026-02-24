package shell

import "strings"

func (m Model) hintText() string {
	switch m.mode {
	case ModeLauncher:
		return "type to filter · ↑/↓ move · enter run · esc close"
	case ModeMain:
		hints := []string{"ctrl+p or c command palette"}
		if appHint := m.appHint(); appHint != "" {
			hints = append([]string{appHint}, hints...)
		}
		if len(m.openApps()) > 1 {
			hints = append(hints, "tab next app")
		}
		hints = append(hints, "q quit")
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
