package shell

import (
	"github.com/Noudea/glyph/internal/registry"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case splashTickMsg:
		return m.updateSplashTick()
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
	case ModeSplash:
		return m.updateSplash(msg)
	case ModeMain:
		return m.updateMain(msg)
	case ModeLauncher:
		return m.updateLauncher(msg)
	}

	return m, nil
}

func (m *Model) updateSplash(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Any key skips the splash and opens the launcher.
	_ = msg
	m.finishSplash()
	return m, nil
}

func (m *Model) updateSplashTick() (tea.Model, tea.Cmd) {
	if m.mode != ModeSplash {
		return m, nil
	}
	m.splashFrame++
	if m.splashFrame >= splashFrameCount {
		m.finishSplash()
		return m, nil
	}
	return m, splashTickCmd()
}

func (m *Model) updateMain(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	if commandID, ok := m.resolveMainCommandIDForKey(key); ok {
		return m, m.executeCommand(commandID)
	}

	if module, ok := m.activeModule(); ok {
		updated, cmd := module.Update(m.context(), msg)
		m.modules[updated.ID()] = updated
		return m, cmd
	}
	return m, nil
}

func (m *Model) updateLauncher(msg tea.Msg) (tea.Model, tea.Cmd) {
	var inputCmd tea.Cmd
	m.launcherInput, inputCmd = m.launcherInput.Update(msg)
	m.clampLauncherCursor()
	var actionCmd tea.Cmd

	key := ""
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		key = keyMsg.String()
	}

	switch key {
	case "esc":
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
			selected := cmds[m.launcherCursor]
			actionCmd = m.executeCommand(selected.ID)
		}
		m.launcherInput.Blur()
		if m.mode == ModeLauncher {
			m.mode = ModeMain
		}
	default:
		if commandID, ok := m.resolveMainCommandIDForKey(key); ok {
			if commandID == commandLauncherOpen {
				m.launcherInput.Blur()
				m.mode = ModeMain
				return m, nil
			}
			actionCmd = m.executeCommand(commandID)
			m.launcherInput.Blur()
			if m.mode == ModeLauncher {
				m.mode = ModeMain
			}
		}
	}

	return m, tea.Batch(inputCmd, actionCmd)
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
