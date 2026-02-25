package shell

import (
	"github.com/Noudea/glyph/internal/core"
	"github.com/Noudea/glyph/internal/registry"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case registry.ChangedMsg:
		m.refreshRegistry()
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "ctrl+c":
		return m, tea.Quit
	}

	switch m.mode {
	case ModeMain:
		return m.updateMain(msg)
	case ModeLauncher:
		return m.updateLauncher(msg)
	}

	return m, nil
}

func (m *Model) updateMain(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	if key == "ctrl+p" || key == "ctrl+k" || key == "alt+p" {
		m.mode = ModeLauncher
		m.launcherInput.SetValue("")
		m.launcherInput.Focus()
		return m, nil
	}

	if module, ok := m.activeModule(); ok {
		updated, cmd := module.Update(m.context(), msg)
		m.modules[updated.ID()] = updated
		return m, cmd
	}
	return m, nil
}

func (m *Model) updateLauncher(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.launcherInput, cmd = m.launcherInput.Update(msg)
	m.clampLauncherCursor()

	key := ""
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		key = keyMsg.String()
	}

	switch key {
	case "esc", "ctrl+p", "ctrl+k", "alt+p":
		m.launcherInput.Blur()
		m.mode = ModeMain
		return m, nil
	case "up", "k":
		if m.launcherCursor > 0 {
			m.launcherCursor--
		}
	case "down", "j":
		if m.launcherCursor < len(m.filteredCommands())-1 {
			m.launcherCursor++
		}
	case "enter":
		cmds := m.filteredCommands()
		if len(cmds) > 0 && m.launcherCursor >= 0 && m.launcherCursor < len(cmds) {
			cmd := cmds[m.launcherCursor]
			if cmd.Kind == core.CommandApp {
				m.openAppByID(cmd.ID)
			}
		}
		m.launcherInput.Blur()
		m.mode = ModeMain
	}

	return m, cmd
}

func (m *Model) clampLauncherCursor() {
	cmds := m.filteredCommands()
	if len(cmds) == 0 {
		m.launcherCursor = 0
		return
	}
	if m.launcherCursor < 0 {
		m.launcherCursor = 0
		return
	}
	if m.launcherCursor >= len(cmds) {
		m.launcherCursor = len(cmds) - 1
	}
}
