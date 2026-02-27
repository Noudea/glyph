package shell

import (
	"path/filepath"
	"strings"

	"github.com/Noudea/glyph/internal/core"
)

func (m Model) filteredCommands() []core.Command {
	commands := m.launcherCommands()
	query := strings.TrimSpace(strings.ToLower(m.launcherInput.Value()))
	if query == "" {
		return commands
	}

	out := make([]core.Command, 0, len(commands))
	for _, command := range commands {
		if strings.Contains(strings.ToLower(command.Label), query) ||
			strings.Contains(strings.ToLower(command.ID), query) ||
			strings.Contains(strings.ToLower(command.Shortcut), query) {
			out = append(out, command)
		}
	}
	return out
}

func (m Model) findCommandByID(id string) (core.Command, bool) {
	if m.state == nil || id == "" {
		return core.Command{}, false
	}
	for _, command := range m.state.Commands {
		if command.ID == id {
			return command, true
		}
	}
	return core.Command{}, false
}

func (m Model) workspaceTitle() string {
	abs := filepath.Clean(m.startDir)
	if abs == "." || abs == string(filepath.Separator) {
		return "glyph · " + abs
	}
	base := filepath.Base(abs)
	if base == "." || base == string(filepath.Separator) || base == "" {
		return "glyph · " + abs
	}
	return "glyph · " + base
}

func (m Model) workspaceHint() string {
	return "cwd: " + filepath.Clean(m.startDir)
}
