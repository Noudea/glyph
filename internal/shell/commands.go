package shell

import (
	"errors"
	"path/filepath"

	"github.com/Noudea/glyph/internal/core"
)

const (
	actionWorkspaceProject = "workspace.project"
	actionWorkspaceGlobal  = "workspace.global"
	actionWorkspaceToggle  = "workspace.toggle"
	actionWorkspaceCreate  = "workspace.create"
	actionTabsNext         = "tabs.next"
)

func (m Model) launcherCommands() []core.Command {
	size := len(m.workspaceActionCommands())
	if m.state != nil {
		size += len(m.state.Commands)
	}
	commands := make([]core.Command, 0, size)
	commands = append(commands, m.workspaceActionCommands()...)
	if m.state != nil {
		for _, cmd := range m.state.Commands {
			item := cmd
			item.Shortcut = m.primaryShortcut(item.ID, item.Shortcut)
			commands = append(commands, item)
		}
	}
	return commands
}

func (m Model) workspaceActionCommands() []core.Command {
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
			ID:       actionTabsNext,
			Label:    "tabs: next",
			Kind:     core.CommandAction,
			Group:    "tabs",
			Shortcut: m.primaryShortcut(actionTabsNext, ""),
		},
		{
			ID:       actionWorkspaceToggle,
			Label:    "workspace: toggle global/project",
			Kind:     core.CommandAction,
			Group:    "workspace",
			Shortcut: m.primaryShortcut(actionWorkspaceToggle, ""),
		},
		{
			ID:       projectID,
			Label:    projectLabel,
			Kind:     core.CommandAction,
			Group:    "workspace",
			Shortcut: m.primaryShortcut(projectID, ""),
		},
		{
			ID:       actionWorkspaceGlobal,
			Label:    globalLabel,
			Kind:     core.CommandAction,
			Group:    "workspace",
			Shortcut: m.primaryShortcut(actionWorkspaceGlobal, ""),
		},
	}
}
