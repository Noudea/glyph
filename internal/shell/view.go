package shell

import (
	hintbarview "github.com/Noudea/glyph/internal/view/hintbar"
	launcherview "github.com/Noudea/glyph/internal/view/launcher"
	splashview "github.com/Noudea/glyph/internal/view/splash"
	topbarview "github.com/Noudea/glyph/internal/view/topbar"
	"github.com/charmbracelet/lipgloss"
)

const (
	mainContentPadX = 1
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
		Title:     m.workspaceTitle(),
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
		body = m.renderActiveModulePadded(m.width, bodyHeight)
	}

	return lipgloss.JoinVertical(lipgloss.Left, top, body)
}

func (m *Model) renderActiveModulePadded(width, height int) string {
	if width <= 0 || height <= 0 {
		return m.renderActiveModule(width, height)
	}

	if width <= (mainContentPadX*2)+1 {
		return m.renderActiveModule(width, height)
	}

	innerWidth := width - (mainContentPadX * 2)
	content := m.renderActiveModule(innerWidth, height)
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Padding(0, mainContentPadX).
		Render(content)
}

func (m *Model) renderActiveModule(width, height int) string {
	if module, ok := m.activeModule(); ok {
		return module.View(width, height)
	}
	return lipgloss.NewStyle().Width(width).Height(height).Render("")
}
