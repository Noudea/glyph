package shell

import (
	"sort"
	"strings"

	"github.com/Noudea/glyph/internal/core"
)

const (
	commandSourceGlobal  = "global"
	commandSourceProject = "project"
	commandSourceManaged = "managed"
)

func (m Model) launcherCommands() []core.Command {
	if m.state == nil || len(m.state.Commands) == 0 {
		return nil
	}

	out := make([]core.Command, 0, len(m.state.Commands))
	for _, command := range m.state.Commands {
		item := command
		item.Shortcut = m.primaryShortcut(item.ID, item.Shortcut)
		out = append(out, item)
	}

	sort.SliceStable(out, func(i, j int) bool {
		left := strings.ToLower(strings.TrimSpace(out[i].Label))
		right := strings.ToLower(strings.TrimSpace(out[j].Label))
		if left == right {
			return out[i].ID < out[j].ID
		}
		return left < right
	})

	return out
}
