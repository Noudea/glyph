package shell

import (
	hintbarview "github.com/Noudea/glyph/internal/modules/hintbar/view"
	launcherview "github.com/Noudea/glyph/internal/modules/launcher/view"
	topbarview "github.com/Noudea/glyph/internal/modules/topbar/view"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	hints := hintbarview.Render(hintbarview.ViewState{
		Width: m.width,
		Text:  m.hintText(),
	})

	contentHeight := m.height
	if m.height > 0 {
		contentHeight = m.height - hintbarview.Height(m.width)
		if contentHeight < 1 {
			contentHeight = 1
		}
	}

	content := m.renderContent(contentHeight)

	if hints == "" {
		return content
	}
	return lipgloss.JoinVertical(lipgloss.Left, content, hints)
}

func (m *Model) renderContent(contentHeight int) string {
	switch m.mode {
	case ModeMain:
		return renderMain(m, contentHeight)
	case ModeLauncher:
		return launcherview.Render(launcherview.ViewState{
			InputView: m.launcherInput.View(),
			Commands:  m.filteredCommands(),
			Cursor:    m.launcherCursor,
			Width:     m.width,
			Height:    contentHeight,
		})
	default:
		return renderMain(m, contentHeight)
	}
}

func renderMain(m *Model, contentHeight int) string {
	top := topbarview.Render(topbarview.ViewState{
		Width:     m.width,
		Title:     "glyph",
		Tabs:      m.openApps(),
		ActiveApp: m.state.ActiveApp,
	})
	if top == "" {
		return ""
	}

	body := ""
	if contentHeight > 0 {
		bodyHeight := contentHeight - topbarview.Height(m.width)
		if bodyHeight < 1 {
			bodyHeight = 1
		}
		body = m.renderActiveModule(m.width, bodyHeight)
	}

	return lipgloss.JoinVertical(lipgloss.Left, top, body)
}

func (m *Model) renderActiveModule(width, height int) string {
	if module, ok := m.activeModule(); ok {
		return module.View(width, height)
	}
	return lipgloss.NewStyle().Width(width).Height(height).Render("")
}
