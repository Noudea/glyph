package shell

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/Noudea/glyph/internal/core"
)

func (m Model) openApps() []core.Command {
	if m.state == nil {
		return nil
	}
	if len(m.state.OpenApps) == 0 {
		return nil
	}
	out := make([]core.Command, 0, len(m.state.OpenApps))
	for _, id := range m.state.OpenApps {
		if cmdItem, ok := m.findCommandByID(id); ok && cmdItem.Kind == core.CommandApp && m.hasModule(id) {
			out = append(out, cmdItem)
		}
	}
	return out
}

func (m *Model) cycleOpenApp(delta int) {
	open := m.openApps()
	if len(open) == 0 || m.state == nil {
		return
	}
	index := 0
	for i, a := range open {
		if a.ID == m.state.ActiveApp {
			index = i
			break
		}
	}
	index = (index + delta + len(open)) % len(open)
	m.state.ActiveApp = open[index].ID
}

func (m Model) filteredCommands() []core.Command {
	commands := m.launcherCommands()
	query := strings.TrimSpace(strings.ToLower(m.launcherInput.Value()))
	if query == "" {
		return commands
	}
	out := make([]core.Command, 0, len(commands))
	for _, cmd := range commands {
		if strings.Contains(strings.ToLower(cmd.Label), query) || strings.Contains(strings.ToLower(cmd.ID), query) {
			out = append(out, cmd)
		}
	}
	return out
}

func (m Model) launcherCommands() []core.Command {
	size := len(m.workspaceCommands())
	if m.state != nil {
		size += len(m.state.Commands)
	}
	commands := make([]core.Command, 0, size)
	commands = append(commands, m.workspaceCommands()...)
	if m.state != nil {
		commands = append(commands, m.state.Commands...)
	}
	return commands
}

func (m Model) workspaceCommands() []core.Command {
	projectID := actionWorkspaceProject
	projectLabel := "workspace: use project"
	globalLabel := "workspace: use global"
	if m.workspace.Kind == core.WorkspaceProject {
		projectLabel = "workspace: use project [active]"
	} else {
		globalLabel = "workspace: use global [active]"
		project, err := m.resolver.ResolveProject(false)
		switch {
		case err == nil:
			projectName := filepath.Base(project.ProjectPath)
			if projectName == "." || projectName == string(filepath.Separator) || projectName == "" {
				projectName = "project"
			}
			projectLabel = "workspace: use project (" + projectName + ")"
		case errors.Is(err, core.ErrNoProjectWorkspace):
			projectID = actionWorkspaceCreate
			projectLabel = "workspace: create project workspace"
		}
	}
	return []core.Command{
		{
			ID:    actionWorkspaceToggle,
			Label: "workspace: toggle global/project",
			Kind:  core.CommandAction,
			Group: "workspace",
		},
		{
			ID:    projectID,
			Label: projectLabel,
			Kind:  core.CommandAction,
			Group: "workspace",
		},
		{
			ID:    actionWorkspaceGlobal,
			Label: globalLabel,
			Kind:  core.CommandAction,
			Group: "workspace",
		},
	}
}

func (m *Model) openAppByID(id string) {
	if m.state == nil || id == "" {
		return
	}
	if cmd, ok := m.findCommandByID(id); !ok || cmd.Kind != core.CommandApp {
		return
	}
	if !m.hasModule(id) {
		return
	}
	for _, open := range m.state.OpenApps {
		if open == id {
			m.state.ActiveApp = id
			return
		}
	}
	m.state.OpenApps = append(m.state.OpenApps, id)
	m.state.ActiveApp = id
}

func (m Model) findCommandByID(id string) (core.Command, bool) {
	if m.state == nil {
		return core.Command{}, false
	}
	for _, cmd := range m.state.Commands {
		if cmd.ID == id {
			return cmd, true
		}
	}
	return core.Command{}, false
}

func (m Model) activeModule() (core.Module, bool) {
	if m.state == nil {
		return nil, false
	}
	return m.moduleByID(m.state.ActiveApp)
}

func (m Model) moduleByID(id string) (core.Module, bool) {
	if id == "" {
		return nil, false
	}
	module, ok := m.modules[id]
	return module, ok
}

func (m Model) hasModule(id string) bool {
	_, ok := m.modules[id]
	return ok
}

func (m Model) context() core.CoreContext {
	return core.CoreContext{RootPath: m.root}
}

func (m Model) workspaceTitle() string {
	if m.workspace.Kind == core.WorkspaceProject {
		name := filepath.Base(m.workspace.ProjectPath)
		if name == "." || name == string(filepath.Separator) || name == "" {
			name = "project"
		}
		return "glyph · project:" + name
	}
	return "glyph · global"
}

func (m Model) workspaceHint() string {
	if m.workspace.Kind == core.WorkspaceProject {
		if m.workspace.ProjectPath == "" {
			return "workspace: project"
		}
		name := filepath.Base(m.workspace.ProjectPath)
		if name == "." || name == string(filepath.Separator) || name == "" {
			name = m.workspace.ProjectPath
		}
		return "workspace: project (" + name + ")"
	}
	return "workspace: global"
}
