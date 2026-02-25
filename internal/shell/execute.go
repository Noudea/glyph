package shell

import (
	"github.com/Noudea/glyph/internal/core"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) executeCommand(commandID string) tea.Cmd {
	switch commandID {
	case commandLauncherOpen:
		m.openLauncher()
		return nil
	case actionWorkspaceToggle, actionWorkspaceProject, actionWorkspaceCreate, actionWorkspaceGlobal:
		return m.runWorkspaceAction(commandID)
	default:
		m.openAppByID(commandID)
		return nil
	}
}

func (m *Model) openLauncher() {
	m.mode = ModeLauncher
	m.launcherInput.SetValue("")
	m.launcherInput.Focus()
	m.clampLauncherCursor()
}

func (m *Model) runWorkspaceAction(actionID string) tea.Cmd {
	resolve := func(kind core.WorkspaceKind) (core.Workspace, error) {
		if kind == core.WorkspaceProject {
			return m.resolver.ResolveProject(true)
		}
		return m.resolver.ResolveGlobal()
	}

	switch actionID {
	case actionWorkspaceToggle:
		if m.workspace.Kind == core.WorkspaceProject {
			workspace, err := resolve(core.WorkspaceGlobal)
			if err != nil {
				m.err = err.Error()
				return nil
			}
			return m.setWorkspace(workspace)
		}
		workspace, err := resolve(core.WorkspaceProject)
		if err != nil {
			m.err = err.Error()
			return nil
		}
		return m.setWorkspace(workspace)
	case actionWorkspaceProject, actionWorkspaceCreate:
		workspace, err := resolve(core.WorkspaceProject)
		if err != nil {
			m.err = err.Error()
			return nil
		}
		return m.setWorkspace(workspace)
	case actionWorkspaceGlobal:
		workspace, err := resolve(core.WorkspaceGlobal)
		if err != nil {
			m.err = err.Error()
			return nil
		}
		return m.setWorkspace(workspace)
	default:
		return nil
	}
}
