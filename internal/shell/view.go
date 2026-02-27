package shell

import (
	"strings"

	hintbarview "github.com/Noudea/glyph/internal/view/hintbar"
	launcherview "github.com/Noudea/glyph/internal/view/launcher"
	splashview "github.com/Noudea/glyph/internal/view/splash"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	showHints := m.mode != ModeSplash

	hints := hintbarview.Render(hintbarview.ViewState{
		Width: m.width,
		Text:  m.hintText(),
	})

	contentHeight := m.height
	if showHints && m.height > 0 {
		contentHeight = m.height - hintbarview.Height(m.width)
		if contentHeight < 1 {
			contentHeight = 1
		}
	}

	content := m.renderContent(contentHeight)

	if !showHints || hints == "" {
		return content
	}
	return lipgloss.JoinVertical(lipgloss.Left, content, hints)
}

func (m *Model) renderContent(contentHeight int) string {
	switch m.mode {
	case ModeSplash:
		return splashview.Render(splashview.ViewState{
			Frame:  m.splashFrame,
			Width:  m.width,
			Height: contentHeight,
		})
	case ModeLauncher:
		return launcherview.Render(launcherview.ViewState{
			InputView: m.launcherInput.View(),
			Commands:  m.filteredCommands(),
			Cursor:    m.launcherCursor,
			Width:     m.width,
			Height:    contentHeight,
		})
	case ModeMain:
		fallthrough
	default:
		return m.renderMain(contentHeight)
	}
}

func (m *Model) renderMain(height int) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF9F68"))
	mutedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8A90A6"))

	lines := []string{
		titleStyle.Render(m.workspaceTitle()),
		"",
		mutedStyle.Render("Press " + m.shortcutsHint(commandLauncherOpen, "ctrl+p/ctrl+k/alt+p") + " to open the command palette"),
		mutedStyle.Render("Run commands in current terminal · current folder"),
	}

	content := strings.Join(lines, "\n")
	if m.width > 0 && height > 0 {
		return lipgloss.Place(m.width, height, lipgloss.Center, lipgloss.Center, content)
	}
	return content
}
