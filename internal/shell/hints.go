package shell

import "strings"

func (m Model) hintText() string {
	switch m.mode {
	case ModeLauncher:
		return "type to filter · ↑/↓ move · enter run · esc/ctrl+p/ctrl+k/alt+p close"
	case ModeMain:
		hints := []string{"ctrl+p, ctrl+k, or alt+p command palette"}
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
