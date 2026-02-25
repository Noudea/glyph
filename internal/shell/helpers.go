package shell

import (
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
	if m.state == nil {
		return nil
	}
	query := strings.TrimSpace(strings.ToLower(m.launcherInput.Value()))
	if query == "" {
		return m.state.Commands
	}
	out := make([]core.Command, 0, len(m.state.Commands))
	for _, cmd := range m.state.Commands {
		if strings.Contains(strings.ToLower(cmd.Label), query) || strings.Contains(strings.ToLower(cmd.ID), query) {
			out = append(out, cmd)
		}
	}
	return out
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

func (m Model) inputFocused() bool {
	switch m.mode {
	case ModeLauncher:
		return m.launcherInput.Focused()
	default:
		if module, ok := m.activeModule(); ok {
			return module.InputFocused()
		}
		return false
	}
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
